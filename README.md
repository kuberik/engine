# kuberik

<img src="./docs/.vuepress/public/assets/img/logo.svg" height=100 />

----

Kuberik is an extensible pipeline engine for Kubernetes. It enables
execution of pipelines on top of Kubernetes by leveraging full expressiveness
of Kubernetes API.

----

## Project status

**Project status:** *alpha*

Project is in alpha stage. API is still a subject to change. Please do not use it in production.

## Development

#### Dependencies
  - kubectl 1.14+
  - Kubernetes 1.14+
  - Go 1.13+
  - Operator SDK 1.2.0+

#### Prerequisites
  - Authenticated to a Kubernetes cluster (e.g. [kind](https://kind.sigs.k8s.io/))
  - Applied kuberik CRDs on the cluster
    - `kubectl apply -k deploy/crds`
  - Installed kuberik CLI
    - `go install ./cmd/kuberik`

### Running operator locally

Start up the operator:

```shell
make run
```

You can use one of the pipelines from the `docs/examples` directory to execute some workload on kuberik.
```shell
kubectl apply -f docs/examples/hello-world.yaml
```

Trigger the pipeline with `kuberik` cmd.
```shell
kuberik create play --from=hello-world
```

### Testing operator locally

Following runs e2e and unit tests.

```shell
make -j2 test
```
