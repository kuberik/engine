package engine

import (
	"encoding/json"
	"fmt"
	"os"

	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
	"github.com/kuberik/engine/pkg/engine/internal/kustomize"
	"github.com/kuberik/engine/pkg/kubeutils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	batchv1 "k8s.io/api/batch/v1"

	"sigs.k8s.io/kustomize/api/resource"
)

const (
	movieKind = "Movie"
)

func provisionedResourcesLayer(play *corev1alpha1.Play, screenplay string) kustomize.KustomizeLayer {
	kl := kustomize.NewKustomizeLayerRoot()
	for _, p := range play.Spec.Screenplays[0].Provision.Resources {
		kl.AddObjectRaw(p.Raw)
	}
	return kl
}

func generateFinalLayer(play *corev1alpha1.Play, layer kustomize.KustomizeLayer) ([]*resource.Resource, error) {
	layer.Kustomization.NameSuffix = fmt.Sprintf("-%s", play.Name)
	layer.Kustomization.Namespace = play.Namespace
	rm, err := layer.Run()
	if err != nil {
		return nil, err
	}

	return rm.Resources(), nil
}

func generateProvisionedResources(play *corev1alpha1.Play, screenplay string) ([]*resource.Resource, error) {
	return generateFinalLayer(play, provisionedResourcesLayer(play, screenplay))
}

func generateActionJob(play *corev1alpha1.Play, screenplay string, frameID string) batchv1.Job {
	pl := provisionedResourcesLayer(play, screenplay)

	action := newAction(play, frameID)
	jl := pl.AddLayer()
	jl.AddObject(action)

	resources, err := generateFinalLayer(play, jl)
	if err != nil {
		panic("failed creating a job")
	}
	for _, r := range resources {
		if r.GetKind() == "Job" {
			transformedAction := batchv1.Job{}
			transformedActionMarshaled, _ := json.Marshal(r)
			json.Unmarshal(transformedActionMarshaled, &transformedAction)
			return transformedAction
		}
	}

	panic("transformation lost input action")
}

var (
	trueVal       = true
	zero    int32 = 0
)

const (
	ActionAnnotationFrameID = "core.kuberik.io/frameID"
)

func actionLabels(play *corev1alpha1.Play, frameName string) labels.Set {
	return map[string]string{
		// TODO: replace with frame name
		"app.kubernetes.io/name":     frameName,
		"app.kubernetes.io/instance": fmt.Sprintf("%s-%s", frameName, play.Name),
		// TODO: replace with actual version of kuberik
		"app.kubernetes.io/version":    os.Getenv("TODOVERSION"),
		"app.kubernetes.io/component":  "action",
		"app.kubernetes.io/part-of":    play.Name,
		"app.kubernetes.io/managed-by": "kuberik",
	}
}

func newAction(play *corev1alpha1.Play, frameID string) batchv1.Job {
	f := play.Frame(frameID)
	e := f.Action

	// TODO: replace with owner reference
	annotations := map[string]string{
		ActionAnnotationFrameID: frameID,
	}
	for k, v := range play.GetAnnotations() {
		annotations[k] = v
	}

	e.Template.Labels = labels.Merge(e.Template.Labels, actionLabels(play, f.Name))

	if e.BackoffLimit == nil {
		e.BackoffLimit = &zero
	}
	if e.Template.Spec.RestartPolicy == "" {
		e.Template.Spec.RestartPolicy = corev1.RestartPolicyNever
	}

	job := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			// maximum string for job name is 63 characters.
			Name:            f.Name,
			Namespace:       play.Namespace,
			Annotations:     annotations,
			Labels:          e.Template.Labels,
			OwnerReferences: []metav1.OwnerReference{kubeutils.ControllerReference(play)},
		},
		Spec: *e.DeepCopy(),
	}

	return job
}

const (
	// JobLabelPlay is name of a label which stores name of the play that owns frame of this job
	JobLabelPlay = "core.kuberik.io/play"
)

func JobLabelSelector(play *corev1alpha1.Play) labels.Selector {
	ls, _ := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app.kubernetes.io/part-of":    play.Name,
			"app.kubernetes.io/managed-by": "kuberik",
		},
	})
	return ls
}
