package scheduler

import (
	"os/exec"

	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
)

// ShellScheduler runs workloads directly on the local system
type ShellScheduler struct{}

var _ Scheduler = &ShellScheduler{}

// Provision is not implemented for ShellScheduler
func (s *ShellScheduler) Provision(play *corev1alpha1.Play) error {
	panic("not implemented")
}

// Deprovision is not implemented for ShellScheduler
func (s *ShellScheduler) Deprovision(play *corev1alpha1.Play) error {
	panic("not implemented")
}

// Run implements Scheduler interface
func (s *ShellScheduler) Run(play *corev1alpha1.Play, frameID string) error {
	e := play.Frame(frameID).Action

	var args []string
	var command string
	args = append(args, e.Template.Spec.Containers[0].Args...)
	if execCommand := e.Template.Spec.Containers[0].Command; len(execCommand) > 0 {
		args = append(execCommand[1:], args...)
		command = execCommand[0]
	}
	cmd := exec.Command(command, args...)

	cmd.Start()
	go func() {
		cmd.Wait()
	}()

	return nil
}
