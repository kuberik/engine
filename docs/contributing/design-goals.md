# Design goals

[[toc]]

## Portability

Making Kuberik run containers ensures that it is portable. Although current work is mostly targeted on running the workloads on Kubernetes, for other purposes it would also make sense to run it with other Docker-backed schedulers. For example, if you'd want to test your pipelines, it would make sense to enable executing it on your local machine. In this case, implementing plain Docker scheduler (which comes without advanced feature of Kubernetes) would be a good idea.

## Domain specific

While in its core, Kuberik tries to be domain agnostic, i.e. doesn't differ between CI, CD or data processing pipelines, there's a way to extend Kuberik to make it easier to developers to develop their pipelines. This is tackled by giving the ability to extend the way in which Kuberik pipelines are trigger and providing an opinionated of developing pipelines.

## Testability

Many pipeline engines are plainly said untestable. This comes from the fact that pipelines are generally used as high-level workloads, meaning that they integrate very different software into a one continuos execution. Checking if the whole pipelines runs in a single go is usually unpractical as they last too long or touch production services.

Kuberik aims to solve the problems of testability by having frame-level testing. This would allow for verifying each of the frames functionality and input/output integrity.

Simplest way it could work is to define test suits. Each test suit runs (in order) an initialization frame, frame under test (FUT) and verification frame, which verifies that FUT completed correctly. It would also allow to run additional mock services if FUT needs access to some external services and can't interact with real ones during test.

## Full Circle

E.g. Kuberik pipeline should be able to test Kuberik locally, run the CI system for itself, and deploy Kuberik itself to production.

## Resiliency

Many pipeline engines fail the pipeline on transient error, or end up in a dirty state and are unable to recover without manual intervention. When the pipeline gets scheduled it should get executed eventually without user intervention.

## Extensible pipeline trigger system (Screeners)
Every pipeline depends on some sort of a trigger. That's what defines _when_ the pipeline should run. This can be a simple UI click, or a webhook. Goal of Kuberik is to define an API with which standalone screeners can indirectly trigger the pipelines.

## Pipeline definition

Although Kuberik tries to avoid creating yet another pipeline with a DSL, this is largely unavoidable. However, it is guided by Kubernetes principles, which makes it more of an API than a DSL. Kuberik workflows are defined as a single YAML file which describes the model of the pipeline.
