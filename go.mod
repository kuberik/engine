module github.com/kuberik/engine

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/prometheus/common v0.4.1
	github.com/sirupsen/logrus v1.4.2
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.3
	sigs.k8s.io/kustomize/api v0.6.6-0.20201204185154-1583cef8d96f
	sigs.k8s.io/kustomize/kyaml v0.10.3-0.20201204185154-1583cef8d96f // indirect
)
