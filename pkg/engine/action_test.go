package engine

import (
	"fmt"
	"testing"

	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestGenerateProvisionedResources(t *testing.T) {
	pvcName := "myclaim"
	pvcResource := []byte(fmt.Sprintf(`
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
    name: %s
`, pvcName))
	screenplayName := "main"
	play := &corev1alpha1.Play{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1alpha1.PlaySpec{
			Screenplays: []corev1alpha1.Screenplay{{
				Name: screenplayName,
				Provision: corev1alpha1.Provision{
					Resources: []runtime.RawExtension{{
						Raw: pvcResource,
					}},
				},
			}},
		},
	}
	provisioned, err := generateProvisionedResources(play, screenplayName)
	if err != nil {
		t.Fatalf("Failed to generate resources: %s", err)
	}

	if want := 1; len(provisioned) != want {
		t.Errorf("Want %d provisioned resources, but got %d", want, len(provisioned))
	}

	if want := fmt.Sprintf("%s-%s", pvcName, play.Name); provisioned[0].GetName() != want {
		t.Errorf("Want '%s' name for provisioned resource, but got %s", want, provisioned[0].GetName())
	}
}

func TestGenerateJob(t *testing.T) {
	pvcName := "myclaim"
	pvcResource := []byte(fmt.Sprintf(`
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
    name: %s
`, pvcName))
	screenplayName := "main"
	play := &corev1alpha1.Play{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1alpha1.PlaySpec{
			Screenplays: []corev1alpha1.Screenplay{{
				Name: screenplayName,
				Provision: corev1alpha1.Provision{
					Resources: []runtime.RawExtension{{
						Raw: pvcResource,
					}},
				},
				Scenes: []corev1alpha1.Scene{{
					Frames: []corev1alpha1.Frame{{
						ID: "a",
						Action: &corev1alpha1.Action{
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Volumes: []corev1.Volume{{
										Name: "a",
										VolumeSource: corev1.VolumeSource{
											PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
												ClaimName: pvcName,
											},
										},
									}},
								},
							},
						},
					}},
				}},
			}},
		},
	}
	provisioned, _ := generateProvisionedResources(play, screenplayName)
	job := newAction(play, "a")
	job = generateJob(play, screenplayName, job)

	if provisioned[0].GetName() != job.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ClaimName {
		t.Errorf("Want '%s' name for provisioned resource, but got %s", provisioned[0].GetName(), job.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ClaimName)
	}
}
