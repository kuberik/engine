package scheduler

import (
	corev1alpha1 "github.com/kuberik/engine/pkg/apis/core/v1alpha1"
)

// Scheduler implements a way for launching Actions
type Scheduler interface {
	Run(play *corev1alpha1.Play, frameID string) error
}
