package kustomize

import (
	"encoding/json"
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestSingleLayer(t *testing.T) {
	kl := NewKustomizeLayerRoot()

	kl.AddObjectRaw([]byte(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: foo
    `))

	kl.Kustomization.NameSuffix = "-bar"
	rm, err := kl.Run()
	if err != nil {
		t.Fatalf("Kustomize failed: %s", err)
	}

	var configMapFound bool
	for _, r := range rm.Resources() {
		switch r.GetKind() {
		case "ConfigMap":
			configMapFound = true
			if want := "foo-bar"; r.GetName() != want {
				t.Errorf("Want '%s' name for ConfigMap, but got %s", want, r.GetName())
			}
		}
	}

	if !configMapFound {
		t.Error("Wanted ConfigMap not found")
	}
}

func TestNestedLayers(t *testing.T) {
	kl := NewKustomizeLayerRoot()

	kl.AddObjectRaw([]byte(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: foo
    `))

	nl := kl.AddLayer()
	nl.AddObjectRaw([]byte(`
apiVersion: v1
kind: Pod
metadata:
  name: foo
spec:
  containers:
  - name: hello
    envFrom:
    - configMapRef:
        name: foo
    `))

	nl.Kustomization.NameSuffix = "-bar"
	rm, err := nl.Run()
	if err != nil {
		t.Fatalf("Kustomize failed: %s", err)
	}

	var configMapFound bool
	var podFound bool
	for _, r := range rm.Resources() {
		switch r.GetKind() {
		case "ConfigMap":
			configMapFound = true
			if want := "foo-bar"; r.GetName() != want {
				t.Errorf("Want '%s' name for ConfigMap, but got %s", want, r.GetName())
			}
		case "Pod":
			podFound = true
			podMarshaled, _ := r.MarshalJSON()
			fmt.Println(string(podMarshaled))
			pod := corev1.Pod{}
			json.Unmarshal(podMarshaled, &pod)
			got := pod.Spec.Containers[0].EnvFrom[0].ConfigMapRef.Name
			if want := "foo-bar"; got != want {
				t.Errorf("Want '%s' name for ConfigMap, but got %s", want, got)
			}
		}
	}

	if !configMapFound {
		t.Error("Wanted ConfigMap not found")
	}

	if !podFound {
		t.Error("Wanted Pod not found")
	}
}
