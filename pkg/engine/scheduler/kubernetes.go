package scheduler

import (
	"context"
	"fmt"
	"sync"

	corev1alpha1 "github.com/kuberik/kuberik/pkg/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	batchv1 "k8s.io/api/batch/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// JobLabelPlay is name of a label which stores name of the play that owns frame of this job
	JobLabelPlay = "kuberik.io/play"

	// JobAnnotationFrameID is name of a label which stores ID of the frame that owns the job
	JobAnnotationFrameID = "kuberik.io/frameID"
)

var updateLock sync.Mutex

// KubernetesScheduler defines a Scheduler which executes Plays on Kubernetes
type KubernetesScheduler struct {
	client client.Client
}

var _ Scheduler = &KubernetesScheduler{}

// NewKubernetesScheduler creates a Kubernetes scheduler
func NewKubernetesScheduler(c client.Client) *KubernetesScheduler {
	return &KubernetesScheduler{
		client: c,
	}
}

// Run implements Scheduler interface
func (ks *KubernetesScheduler) Run(play *corev1alpha1.Play, frameID string) error {
	jobDefinition := newRunJob(play, frameID)

	// Try to recover first
	job := &batchv1.Job{}
	err := ks.client.Get(context.TODO(), types.NamespacedName{
		Name:      jobDefinition.Name,
		Namespace: jobDefinition.Namespace,
	}, job)
	if err == nil {
		return nil
	}

	return ks.client.Create(context.TODO(), jobDefinition)
}

var (
	falseVal       = false
	trueVal        = true
	zero     int32 = 0
)

func newRunJob(play *corev1alpha1.Play, frameID string) *batchv1.Job {
	e := play.Frame(frameID).Action

	annotations := map[string]string{
		JobAnnotationFrameID: frameID,
	}
	for k, v := range play.GetAnnotations() {
		annotations[k] = v
	}

	labels := map[string]string{
		JobLabelPlay: play.Name,
	}

	if e.Template.Labels == nil {
		e.Template.Labels = make(map[string]string)
	}
	for lk, lv := range labels {
		e.Template.Labels[lk] = lv
	}

	if e.BackoffLimit == nil {
		e.BackoffLimit = &zero
	}
	if e.Template.Spec.RestartPolicy == "" {
		e.Template.Spec.RestartPolicy = corev1.RestartPolicyNever
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			// maximum string for job name is 63 characters.
			Name:        fmt.Sprintf("%.46s-%.16s", play.Name, frameID),
			Namespace:   play.Namespace,
			Annotations: annotations,
			Labels:      labels,
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: play.APIVersion,
				Kind:       play.Kind,
				Name:       play.Name,
				UID:        play.UID,
				Controller: &trueVal,
			}},
		},
		Spec: *e,
	}

	return job
}
