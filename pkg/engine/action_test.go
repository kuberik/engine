package engine

import (
	"fmt"
	"testing"

	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
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
spec:
    resources:
        requests:
            storage: 8Gi
    accessModes:
    - ReadWriteOnce
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
	fmt.Println(provisioned[0])
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
