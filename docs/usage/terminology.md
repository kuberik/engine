# Terminology

Kuberik pipelines have an unique terminology inspired by movies.
In short, to run a workflow with Kuberik, you must create a [Screenplay] containing a list of one or more [Scenes][Scene], each consisting of a collection of [Frames][Frame].

## Screenplay
A *Screenplay* is the main abstraction in a *Kuberik* workflow. It contains a list of one or more *Scenes* which are executed serially.
*Screenplays* are designed to be simple and human editable so it's defined in the *yaml* format. Besides containing a list of *Scenes*, *Screenplays* also have their own properties.
[See the full screenplay reference](./screenplay-reference.md).

## Movie
Movie is a [CRD] which describes a screenplay. To run a screenplay, an instance of [Play] needs to be created from the movie's template. An instance can be created manually with Kuberik CLI or automatically by a screener.

## Play
Play is an instance of a [Movie].
Every execution of a [Screenplay] is defined by a separate Play object. This ensures that screenplay doesn't change during the execution and allows users to freely change originating [Movie] at any time.

## Scene
A [Scene] defines execution of multiple frames. [Frames][Frame] are executed in parallel.

## Frame
A [Frame] is the smallest logical piece of the workflow. In itself, it can't decide if something should execute in series or in parallel.
It just defines a single logical piece of work to be executed - either an [Action] or a [Story].

### Action
Defines the `spec` of a [Job][JobSpec] to be executed.

### Story
Story is defined as a nested [Screenplay] inside of a [Frame].

## Examples

```yaml
scenes:
- name: hello
  frames:
  - name: hello
    action:
      template:
        spec:
          containers:
          - name: hello
            command: ["echo", "Hello Dave."]
            image: alpine
```

[Movie]: #movie
[Screenplay]: #screenplay
[Scene]: #scene
[Play]: #play
[Frame]: #frame
[Action]: #action
[Story]: #story
[JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#jobspec-v1-batch
[CRD]: https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#customresourcedefinitions
