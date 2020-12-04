package k8s

import (
	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	// JobLabelPlay is name of a label which stores name of the play that owns frame of this job
	JobLabelPlay = "kuberik.io/play"

	// JobAnnotationFrameID is name of a label which stores ID of the frame that owns the job
	JobAnnotationFrameID = "kuberik.io/frameID"
)

func JobLabelSelector(play *corev1alpha1.Play) labels.Selector {
	ls, _ := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			JobLabelPlay: play.Name,
		},
	})
	return ls
}
