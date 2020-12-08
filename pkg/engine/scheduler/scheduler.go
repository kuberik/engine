package scheduler

import (
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/kustomize/api/resource"
)

// Scheduler implements a way for launching Actions
type Scheduler interface {
	Run(batchv1.Job) error
	provisioner
}

type provisioner interface {
	Provision([]*resource.Resource) error
	Deprovision([]*resource.Resource) error
}
