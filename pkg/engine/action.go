package engine

import (
	"fmt"
	"os"

	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	batchv1 "k8s.io/api/batch/v1"

	"sigs.k8s.io/kustomize/api/builtins"
	"sigs.k8s.io/kustomize/api/provider"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/resource"
	ktypes "sigs.k8s.io/kustomize/api/types"
)

const (
	movieKind = "Movie"
)

func provisionedResources(play *corev1alpha1.Play, screenplay string) ([]*resource.Resource, error) {
	var rf = provider.NewDefaultDepProvider().GetResourceFactory()
	var provision []*resource.Resource
	for _, p := range play.Spec.Screenplays[0].Provision.Resources {
		r, err := rf.FromBytes(p.Raw)
		if err != nil {
			return nil, err
		}
		provision = append(provision, r)
	}
	return provision, nil
}

func resourceTransform(play *corev1alpha1.Play, resources ...*resource.Resource) ([]*resource.Resource, error) {
	rm := resmap.NewFactory(provider.NewDefaultDepProvider().GetResourceFactory(), nil).FromResourceSlice(resources)

	nameTransformer := builtins.PrefixSuffixTransformerPlugin{Suffix: fmt.Sprintf("-%s", play.Name), FieldSpecs: []ktypes.FieldSpec{
		{Path: "metadata/name"},
	}}

	err := nameTransformer.Transform(rm)
	if err != nil {
		return nil, err
	}

	return rm.Resources(), nil
}

func generateProvisionedResources(play *corev1alpha1.Play, screenplay string) ([]*resource.Resource, error) {
	provisoned, err := provisionedResources(play, screenplay)
	if err != nil {
		return nil, err
	}

	transformed, err := resourceTransform(play, provisoned...)
	if err != nil {
		return nil, err
	}

	return transformed, nil
}

var (
	falseVal       = false
	trueVal        = true
	zero     int32 = 0
)

const (
	ActionAnnotationFrameID = "core.kuberik.io/frameID"
)

func actionName(play *corev1alpha1.Play, frameID string) string {
	return fmt.Sprintf("%.46s-%.16s", play.Name, frameID)
}

func actionLabels(play *corev1alpha1.Play, frameID string) labels.Set {
	return map[string]string{
		// TODO replace with frame name
		"app.kubernetes.io/name":     frameID,
		"app.kubernetes.io/instance": actionName(play, frameID),
		// TODO replace with actual version of kuberik
		"app.kubernetes.io/version":    os.Getenv("TODOVERSION"),
		"app.kubernetes.io/component":  "action",
		"app.kubernetes.io/part-of":    play.Name,
		"app.kubernetes.io/managed-by": "kuberik",
	}
}

func newAction(play *corev1alpha1.Play, frameID string) batchv1.Job {
	e := play.Frame(frameID).Action

	// TODO replace with owner reference
	annotations := map[string]string{
		ActionAnnotationFrameID: frameID,
	}
	for k, v := range play.GetAnnotations() {
		annotations[k] = v
	}

	e.Template.Labels = labels.Merge(e.Template.Labels, actionLabels(play, frameID))

	if e.BackoffLimit == nil {
		e.BackoffLimit = &zero
	}
	if e.Template.Spec.RestartPolicy == "" {
		e.Template.Spec.RestartPolicy = corev1.RestartPolicyNever
	}

	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			// maximum string for job name is 63 characters.
			Name:        actionName(play, frameID),
			Namespace:   play.Namespace,
			Annotations: annotations,
			Labels:      e.Template.Labels,
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

const (
	// JobLabelPlay is name of a label which stores name of the play that owns frame of this job
	JobLabelPlay = "core.kuberik.io/play"
)

func JobLabelSelector(play *corev1alpha1.Play) labels.Selector {
	ls, _ := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			JobLabelPlay: play.Name,
		},
	})
	return ls
}
