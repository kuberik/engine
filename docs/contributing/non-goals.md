# Non-goals

Following the design principles, there are also some non-goals of this project which are very unlikely to change.

_YAML composition is in grey zone. Although there's no plan to ever implement this functionality in Kuberik core, there probably should be an opinionated way of doing this. Otherwise, this makes Kuberik less usable._

## Defining a complex DSL

DSL of Kuberik should be driven by scheduling capabilities and not by convenience. In modern programming languages, for example, there are various ways of defining a loop. This however results in a higher language complexity. Kuberik aims to remove unnecessary features such as this, and the use of loops in Kuberik would only make sense to schedule a workload with dynamic resource requirements. With that in mind, we for example implement dynamic task copying which creates a number of copies of the task.

## YAML Composition

As Kuberik pipeline is defined a single YAML file, logical demand is to split a huge YAML file into multiple smaller ones. There are already some solutions which are trying to solve this, especially in Kubernetes space. In theory, this can be achieved with any programming language that can output YAML. Kuberik pipeline is **a model** and therefore model generation and composition is highly encouraged.

## Dependency Management

If you want to use YAML composition, it's likely that you'll also need some sort of dependency management. Dependency management is a problem in itself and every software that needs it struggles to implement it well. Therefore, we'd recommend to rely directly on dependency management solutions of whatever technology you use for YAML composition rather than building one from scratch.

## Built-in Integrations

This goes straight against the second principle behind Kuberik - extendability. Kuberik is API driven, and as such, there's no room for integrating with specific services such as GitHub.

## Templating

Ability to template the model greatly increases the complexity and thus goes against first design principle. Only templating that's enabled is the one that Kubernetes itself implements and that is environment variable templating on `command` and `args` fields. To increase the flexibility of the pipeline though, we might implement some sort of dynamic patching mechanism so you can overwrite/add fields on fly.
