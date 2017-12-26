# resource-controller

This repository implements a `resource` controller for watching `ResourceClass` resources as
defined with a `CustomResourceDefinition` (CRD).

This is a part of POC implementation of Resource Class with `CRD`, `Custom Scheduler` and `Device Plugin`.

## Build & Run

```sh
# build
$ make resource-controller

# run
# assumes you have a working kubeconfig, not required if operating in-cluster
$ ./_output/resource-controller --kubeconfig ${configfile}

# create a CustomResourceDefinition
$ kubectl create -f examples/resource-class.yaml

# create a custom resource of type ResourceClass
$ kubectl create -f examples/resource-class-cpu.yaml

```

## Design Proposal

Here is [design proposal](https://docs.google.com/document/d/1EyOUHah_4hx3QCY0FGyM8n-1FHhIIlG47ws6BgDPfUk/edit?usp=sharing), FYI.

## Roadmap

- It is still under development now.

## Issues and Contributions

Contributions are welcome!

* You can report a bug by [filing a new issue](https://github.com/caicloud/resource-controller/issues/new)
* You can contribute by opening a [pull request](https://help.github.com/articles/using-pull-requests/)

