package scheduler

import (
	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
)

// DummyScheduler implements Scheduler interface but doesn't run any workload
type DummyScheduler struct {
	// Result is a value that dummy scheduler sets as a result status of any frame played
	Result corev1alpha1.FrameStatus
}

var _ Scheduler = &DummyScheduler{}

// Run implements Scheduler interface
func (s *DummyScheduler) Run(play *corev1alpha1.Play, frameID string) error {
	play.Status.SetFrameStatus(frameID, s.Result)
	return nil
}
