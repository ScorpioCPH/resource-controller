/*
Copyright 2017 The Caicloud Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	clientset "github.com/caicloud/resource-controller/pkg/client/clientset/versioned"
	rcinformers "github.com/caicloud/resource-controller/pkg/client/informers/externalversions"
	rccontroller "github.com/caicloud/resource-controller/pkg/controller"
	"github.com/caicloud/resource-controller/pkg/util/signals"
)

var (
	kubeconfig  string
	masterURL   string
	threadiness int = 1
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "~/.kube/config",
		`Path to a kubeconfig, only required if out-of-cluster.`)
	flag.StringVar(&masterURL, "master", "",
		`The url of the Kubernetes API server,
		 will overrides any value in kubeconfig, only required if out-of-cluster.`)
}

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	resourceClassClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	resourceClassInformerFactory := rcinformers.NewSharedInformerFactory(resourceClassClient, time.Second*30)

	controller := rccontroller.NewController(kubeClient, resourceClassClient, resourceClassInformerFactory)

	go resourceClassInformerFactory.Start(stopCh)

	if err = controller.Run(threadiness, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}
