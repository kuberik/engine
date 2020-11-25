# Core principles

Core principles of Kuberik design are **minimalism** and **extensibility**.

## Minimalism

As a primary principle, minimalism ensures that Kuberik doesn't develop into a complex software that is too difficult to use for an end user. As an example, Kuberik doesn't support DAGs. They can greatly increase the complexity of any pipeline engine and add unnecessary complexity to developed workflows.

## Extensibility

Extensibility on the other hand brings the somewhat opposite of the first principle - it means that there needs to be a way to extend the functionality of Kuberik without modifying the core. At the same time though, the complexity has to be hidden from the user (minimalism). Kuberik achieves this by being API driven. Prime example of this would be screener functionality. Every pipeline software needs to have a way of triggering the pipeline. However, most of existing tools, build the feature right into the core of the software. Most common use-case is triggering a pipeline via GitHub webhook which is very specific not only to the way it is triggered, but to the provider as well. Kuberik aims to solve this by having a generic API which can be triggered by standalone screeners.
