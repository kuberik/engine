package scheduler

import (
	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/kustomize/api/resource"
)

// DummyScheduler implements Scheduler interface but doesn't run any workload
type DummyScheduler struct {
	// Result is a value that dummy scheduler sets as a result status of any frame played
	Result corev1alpha1.FrameStatus
	Play   *corev1alpha1.Play
}

var _ Scheduler = &DummyScheduler{}

// Provision doesn't do anything for DummyScheduler
func (s *DummyScheduler) Provision(resource []*resource.Resource) error {
	return nil
}

// Deprovision doesn't do anything for DummyScheduler
func (s *DummyScheduler) Deprovision(resource []*resource.Resource) error {
	return nil
}

// Run implements Scheduler interface
func (s *DummyScheduler) Run(job batchv1.Job) error {
	// TODO: replace hardcoded value
	s.Play.Status.SetFrameStatus(job.Annotations["core.kuberik.io/frameID"], s.Result)
	return nil
}
