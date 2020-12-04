package k8s

import (
	"context"
	"fmt"

	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
	"github.com/kuberik/engine/pkg/engine/scheduler"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	batchv1 "k8s.io/api/batch/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kustomize/api/provider"
	"sigs.k8s.io/kustomize/api/resource"
)

// KubernetesScheduler defines a Scheduler which executes Plays on Kubernetes
type KubernetesScheduler struct {
	client client.Client
}

var _ scheduler.Scheduler = &KubernetesScheduler{}

// NewKubernetesScheduler creates a Kubernetes scheduler
func NewKubernetesScheduler(c client.Client) *KubernetesScheduler {
	return &KubernetesScheduler{
		client: c,
	}
}

func generateProvisionedResources(play *corev1alpha1.Play) ([]runtime.Object, error) {
	var factory = provider.NewDefaultDepProvider().GetResourceFactory()
	var provision []*resource.Resource
	for _, p := range play.Spec.Screenplays[0].Provision.Resources {
		r, err := factory.FromBytes(p.Raw)
		if err != nil {
			return nil, err
		}
		provision = append(provision, r)

	}

	// var nameSuffix string
	// for _, o := range play.OwnerReferences {
	// 	if o.Kind == "Movie" && o.APIVersion == corev1alpha1.GroupVersion.String() {
	// 		nameSuffix = strings.TrimPrefix(play.Name, o.Name)
	// 	}
	// }
	return nil, nil
}

func (ks *KubernetesScheduler) Provision(play *corev1alpha1.Play) error {
	resources, err := generateProvisionedResources(play)
	if err != nil {
		return err
	}
	for _, r := range resources {
		err := ks.client.Create(context.TODO(), r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ks *KubernetesScheduler) Deprovision(play *corev1alpha1.Play) error {
	resources, _ := generateProvisionedResources(play)
	for _, r := range resources {
		err := ks.client.Delete(context.TODO(), r)
		if err != nil {
			return err
		}
	}
	return nil
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
