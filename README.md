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
