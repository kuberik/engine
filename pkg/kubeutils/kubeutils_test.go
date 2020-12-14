package kubeutils

import (
	"testing"

	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
)

func TestOwnerReference(t *testing.T) {
	play := corev1alpha1.Play{}

	ownerReference := OwnerReference(&play)

	if want := false; *ownerReference.Controller != want {
		t.Errorf("For OwnerReference field, %v false, got %v", want, *ownerReference.Controller)
	}

	if want := "Play"; ownerReference.Kind != want {
		t.Errorf("For Kind field, %s false, got %s", want, ownerReference.Kind)
	}

	if want := "core.kuberik.io/v1alpha1"; ownerReference.APIVersion != want {
		t.Errorf("For APIVersion field, %s false, got %s", want, ownerReference.APIVersion)
	}
}
