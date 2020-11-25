# API reference

# Screenplay

## Screenplay
| Field                |            Type            |                                              Description |
|----------------------|:--------------------------:|---------------------------------------------------------:|
| scenes               |         \[][Scene]         |                                           List of scenes |
| vars                 |          \[][Var]          |                                        List of variables |
| volumeClaimTemplates | \[][PersistentVolumeClaim] | List of volume claim templates required during execution |

## Scene
| Field        |    Type     |                                           Description |
|--------------|:-----------:|------------------------------------------------------:|
| name         |   string    |                                     Name of the scene |
| frames       | \[][Frame]  |                                        List of frames |
| pass         | [Condition] |                                        Pass condition |
| when         | [Condition] |                                     Trigger condition |
| ignoreErrors |    bool     | If `true` pipelines will continue regardless of error |

## Frame
| Field        |   Type    |                                           Description |
|--------------|:---------:|------------------------------------------------------:|
| name         |  string   |                                     Name of the frame |
| action       | [JobSpec] |                                    Job to be executed |
| ignoreErrors |   bool    | If `true` pipelines will continue regardless of error |
| loop         |    int    |       Number of instances of the task to be scheduled |

## Variable
| Field     |    Type     |           Description |
|-----------|:-----------:|----------------------:|
| Name      |   string    |  Name of the variable |
| Value     |   string    | Value of the variable |
| ValueFrom | [VarSource] | Value of the variable |

## VarSource

| Field           |          Type          |                                          Description |
|-----------------|:----------------------:|-----------------------------------------------------:|
| configMapKeyRef | [ConfigMapKeySelector] | Selects a key of a ConfigMap in the Play's namespace |
| secretKeyRef    |  [SecretKeySelector]   |    Selects a key of a secret in the Play's namespace |

## Condition
Condition is alias for type `[]map[string]string`.

[Scene]: #scene
[Frame]: #frame
[JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#jobspec-v1-batch
[Condition]: #condition
[Var]: #variable
[VarSource]: #varsource
[InputFieldSelector]: #InputFieldSelector
[ConfigMapKeySelector]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#configmapkeyselector-v1-core
[SecretKeySelector]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#secretkeyselector-v1-core
[gjsonpath]: https://github.com/tidwall/gjson#path-syntax
[PersistentVolumeClaim]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#persistentvolumeclaim-v1-core
