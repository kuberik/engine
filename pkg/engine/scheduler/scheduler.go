package scheduler

import (
	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
)

// Scheduler implements a way for launching Actions
type Scheduler interface {
	Provisioner
	Run(play *corev1alpha1.Play, frameID string) error
}

type Provisioner interface {
	Provision(play *corev1alpha1.Play) error
	Deprovision(play *corev1alpha1.Play) error
}
