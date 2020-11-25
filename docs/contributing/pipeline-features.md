# Pipeline features

## Jobs
**Implemented**: yes

**Status**: stable

All frames are implemented as Kubernetes Jobs. This gives users full expressiveness of existing API.

## Copies
**Implemented**: yes

**Status**: beta

Loops can be defined through code and as such are unnecessary to be defined in the model of the pipeline. However, running it for the wide range would mean running all of the tasks in a single container, greatly constraining the performance of the code. To enable such workloads, Kuberik implements frame copies. It spawns an identical copy of already defined task with

## Nested Screenplay
**Implemented**: no

**Status**: planned

As DAGs are not supported in Kuberik, there are some cases where pipelines execution would be suboptimal. To solve this issue, Kuberik could execute a screenplay instead of a frame, giving the user possibility to create more complex workflows. This enabled the same functionality as DAGs, but in a way that's much more easy to reason about.

## Variable registering
**Implemented**: no

**Status**: under consideration

Not all pipelines would need this feature to work. In fact, an easy workaround would be to write necessary information to the disk and read it from another step. This however creates a hard dependency between the frames, i.e. one step has to know where the other one wrote the information. It also makes it inconvenient for the user, as they need to write and read the file in a safe location on the disk.
