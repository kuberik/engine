package kubeutils

import (
	"reflect"

	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NamespaceObject(name string) corev1.Namespace {
	return corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

var (
	falseVal = false
)

func OwnerReference(owner metav1.Object) metav1.OwnerReference {
	ref := *metav1.NewControllerRef(
		owner, corev1alpha1.GroupVersion.WithKind(reflect.ValueOf(owner).Elem().Type().Name()),
	)
	ref.Controller = &falseVal
	return ref
}
