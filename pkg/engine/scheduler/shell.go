package scheduler

import (
	"os/exec"

	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/kustomize/api/resource"
)

// ShellScheduler runs workloads directly on the local system
type ShellScheduler struct{}

var _ Scheduler = &ShellScheduler{}

// Provision is not implemented for ShellScheduler
func (s *ShellScheduler) Provision(resource []*resource.Resource) error {
	panic("not implemented")
}

// Deprovision is not implemented for ShellScheduler
func (s *ShellScheduler) Deprovision(resource []*resource.Resource) error {
	panic("not implemented")
}

// Run implements Scheduler interface
func (s *ShellScheduler) Run(job batchv1.Job) error {
	var args []string
	var command string
	args = append(args, job.Spec.Template.Spec.Containers[0].Args...)
	if execCommand := job.Spec.Template.Spec.Containers[0].Command; len(execCommand) > 0 {
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
