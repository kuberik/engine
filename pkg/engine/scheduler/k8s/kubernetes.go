package k8s

import (
	"context"

	"github.com/kuberik/engine/pkg/engine/scheduler"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

func resourcesToObjects(resources ...*resource.Resource) (objects []controllerutil.Object) {
	for _, r := range resources {
		objects = append(objects, &unstructured.Unstructured{Object: r.Map()})
	}
	return
}

func (ks *KubernetesScheduler) createObjects(objects ...controllerutil.Object) error {
	for _, r := range objects {
		err := ks.client.Create(context.TODO(), r)
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

func (ks *KubernetesScheduler) Provision(resources []*resource.Resource) error {
	return ks.createObjects(resourcesToObjects(resources...)...)
}

func (ks *KubernetesScheduler) deleteObjects(objects ...controllerutil.Object) error {
	for _, r := range objects {
		err := ks.client.Delete(context.TODO(), r)
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

func (ks *KubernetesScheduler) Deprovision(resources []*resource.Resource) error {
	return ks.deleteObjects(resourcesToObjects(resources...)...)
}

func (ks *KubernetesScheduler) Run(job batchv1.Job) error {
	return ks.createObjects(&job)
}
