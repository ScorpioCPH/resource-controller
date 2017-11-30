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

package controller

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	rcv1alpha1 "github.com/caicloud/resource-controller/pkg/apis/resourceclass/v1alpha1"
	clientset "github.com/caicloud/resource-controller/pkg/client/clientset/versioned"
	rcscheme "github.com/caicloud/resource-controller/pkg/client/clientset/versioned/scheme"
	informers "github.com/caicloud/resource-controller/pkg/client/informers/externalversions"
	listers "github.com/caicloud/resource-controller/pkg/client/listers/resourceclass/v1alpha1"
)

const controllerAgentName = "resourceclass-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a ResourceClass is synced
	SuccessSynced = "Synced"

	// MessageResourceSynced is the message used for an Event fired when a ResourceClass
	// is synced successfully
	MessageResourceSynced = "ResourceClass synced successfully"
)

// Controller is the controller implementation for ResourceClass resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// resourceclassclientset is a clientset for our own API group
	resourceclassclientset clientset.Interface

	resourceClassLister listers.ResourceClassLister
	resourceClassSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new resourceclass controller
func NewController(
	kubeclientset kubernetes.Interface,
	resourceclassclientset clientset.Interface,
	// kubeInformerFactory kubeinformers.SharedInformerFactory,
	resourceclassInformerFactory informers.SharedInformerFactory) *Controller {

	// obtain references to shared index informers for the ResourceClass type

	resourceClassInformer := resourceclassInformerFactory.Resourceclass().V1alpha1().ResourceClasses()

	// Create event broadcaster
	// Add resourceclass-controller types to the default Kubernetes Scheme so Events can be
	// logged for resourceclass-controller types.
	rcscheme.AddToScheme(scheme.Scheme)
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:          kubeclientset,
		resourceclassclientset: resourceclassclientset,
		resourceClassLister:    resourceClassInformer.Lister(),
		resourceClassSynced:    resourceClassInformer.Informer().HasSynced,
		workqueue:              workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ResourceClasss"),
		recorder:               recorder,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when ResourceClass resources change
	resourceClassInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueResourceClass,
		UpdateFunc: func(old, new interface{}) {
			glog.Info("Update ResourceClass")
			newRC := new.(*rcv1alpha1.ResourceClass)
			oldRC := old.(*rcv1alpha1.ResourceClass)
			if newRC.ResourceVersion == oldRC.ResourceVersion {
				glog.Infof("ResourceVersion not changed: %s", newRC.ResourceVersion)
				// Periodic resync will send update events for all known ResourceClasses.
				// Two different versions of the same ResourceClass will always have different RVs.
				return
			}
			controller.enqueueResourceClass(new)
		},
		DeleteFunc: controller.enqueueResourceClass,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting ResourceClass controller")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	// if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.resourceClassSynced); !ok {
	if ok := cache.WaitForCacheSync(stopCh, c.resourceClassSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	// Launch two workers to process ResourceClass resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, 10*time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		glog.Info("processNextWorkItem.1.1, return")
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		// Run the syncHandler, passing it the namespace/name string of the
		// ResourceClass resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the ResourceClass resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the ResourceClass resource with this namespace/name
	resourceClass, err := c.resourceClassLister.ResourceClasses(namespace).Get(name)
	if err != nil {
		// The ResourceClass resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("resourceClass '%s' in work queue no longer exists", key))
			// remove labels and taints
			err := c.updateLabelsAndTaintsOnNodes(name, []rcv1alpha1.NodeSpec{})
			if err != nil {
				return err
			}
			return nil
		}

		return err
	}

	err = c.updateLabelsAndTaintsOnNodes(resourceClass.ObjectMeta.Name, resourceClass.Spec.Nodes)
	if err != nil {
		return err
	}

	c.recorder.Event(resourceClass, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// markLabelsAndTaintsOnNodes takes a ResourceClass name as label and taint,
// check current states of nodes then update it to be desired states
func (c *Controller) updateLabelsAndTaintsOnNodes(rcName string, updateNodes []rcv1alpha1.NodeSpec) error {
	nodeList, _ := c.kubeclientset.CoreV1().Nodes().List(metav1.ListOptions{})

	rcTaint := corev1.Taint{
		Key:    rcName,
		Value:  "true",
		Effect: corev1.TaintEffectNoSchedule,
	}

	// construct current nodes with this label and taint
	currentLabeledNodes := []string{}
	currentTaintedNodes := []string{}
	for _, node := range nodeList.Items {
		if _, ok := node.Labels[rcName]; ok {
			glog.Infof("currentLabeledNodes append: %s", node.Name)
			currentLabeledNodes = append(currentLabeledNodes, node.Name)
		}

		for _, taint := range node.Spec.Taints {
			if reflect.DeepEqual(taint, rcTaint) {
				glog.Infof("currentTaintedNodes append: %s", node.Name)
				currentTaintedNodes = append(currentTaintedNodes, node.Name)
			}
		}

	}

	// construct desired nodes with this label and taint
	desiredLabeledNodes := []string{}
	desiredTaintedNodes := []string{}

	for _, node := range updateNodes {
		desiredLabeledNodes = append(desiredLabeledNodes, node.Name)
		if node.Dedicated {
			desiredTaintedNodes = append(desiredTaintedNodes, node.Name)
		}
	}

	if err := c.updateLabels(currentLabeledNodes, desiredLabeledNodes, rcName); err != nil {
		return fmt.Errorf("updateLabels error: %s", err.Error())
	}

	if err := c.updateTaints(currentTaintedNodes, desiredTaintedNodes, rcTaint); err != nil {
		return fmt.Errorf("updateTaints error: %s", err.Error())
	}

	return nil
}

func diffNodes(current, desired []string) (added, deleted []string) {
	added = []string{}
	deleted = []string{}

	for _, c := range current {
		found := false
		for _, d := range desired {
			if c == d {
				found = true
				break
			}
		}
		// Not found. We add it to deleted slice
		if !found {
			deleted = append(deleted, c)
		}
	}

	for _, d := range desired {
		found := false
		for _, c := range current {
			if c == d {
				found = true
				break
			}
		}
		// Not found. We add it to added slice
		if !found {
			added = append(added, d)
		}
	}

	return added, deleted
}

func (c *Controller) updateLabels(current, desired []string, rcName string) error {
	// diff current and desired nodes for labeling
	needAddLabelNodes, needDeleteLabelNodes := diffNodes(current, desired)

	for _, nodeName := range needAddLabelNodes {
		// get current node info
		nodeInfo, err := c.kubeclientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
		if err != nil {
			glog.Errorf("Get nodeInfo error: %s", err)
			continue
		}

		glog.Infof("add label on node: %s", nodeName)
		copyNode := nodeInfo.DeepCopy()
		copyNode.Labels[rcName] = "true"

		orginal, _ := json.Marshal(nodeInfo)
		modified, _ := json.Marshal(copyNode)
		patch, err := strategicpatch.CreateTwoWayMergePatch(orginal, modified, nodeInfo)
		if err != nil {
			return err
		}

		_, err = c.kubeclientset.CoreV1().Nodes().Patch(nodeName, types.StrategicMergePatchType, patch)
		if err != nil {
			glog.Errorf("update node err: %v", err)
			return err
		}
	}

	for _, nodeName := range needDeleteLabelNodes {
		// get current node info
		nodeInfo, err := c.kubeclientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
		if err != nil {
			glog.Errorf("Get nodeInfo error: %s", err)
			continue
		}

		glog.Infof("delete label on node: %s", nodeName)

		copyNode := nodeInfo.DeepCopy()
		delete(copyNode.Labels, rcName)

		orginal, _ := json.Marshal(nodeInfo)
		modified, _ := json.Marshal(copyNode)
		patch, err := strategicpatch.CreateTwoWayMergePatch(orginal, modified, nodeInfo)
		if err != nil {
			return err
		}

		_, err = c.kubeclientset.CoreV1().Nodes().Patch(nodeName, types.StrategicMergePatchType, patch)
		if err != nil {
			glog.Errorf("update node err: %v", err)
			return err
		}

	}

	return nil
}

func (c *Controller) updateTaints(current, desired []string, taint corev1.Taint) error {
	// diff current and desired nodes for taint
	needAddTaintNodes, needDeleteTaintNodes := diffNodes(current, desired)

	for _, nodeName := range needAddTaintNodes {
		// get current node info
		nodeInfo, err := c.kubeclientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
		if err != nil {
			glog.Errorf("Get nodeInfo error: %s", err)
			continue
		}

		glog.Infof("add taint on node: %s", nodeName)
		copyNode := nodeInfo.DeepCopy()
		copyNode.Spec.Taints = append(copyNode.Spec.Taints, taint)

		orginal, _ := json.Marshal(nodeInfo)
		modified, _ := json.Marshal(copyNode)
		patch, err := strategicpatch.CreateTwoWayMergePatch(orginal, modified, nodeInfo)
		if err != nil {
			return err
		}

		_, err = c.kubeclientset.CoreV1().Nodes().Patch(nodeName, types.StrategicMergePatchType, patch)
		if err != nil {
			glog.Errorf("update node err: %v", err)
			return err
		}
	}

	for _, nodeName := range needDeleteTaintNodes {
		// get current node info
		nodeInfo, err := c.kubeclientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
		if err != nil {
			glog.Errorf("Get nodeInfo error: %s", err)
			continue
		}

		glog.Infof("delete taint on node: %s", nodeName)
		copyNode := nodeInfo.DeepCopy()

		for i, t := range copyNode.Spec.Taints {
			if t == taint {
				copyNode.Spec.Taints = append(copyNode.Spec.Taints[:i], copyNode.Spec.Taints[i+1:]...)
			}
		}

		orginal, _ := json.Marshal(nodeInfo)
		modified, _ := json.Marshal(copyNode)
		patch, err := strategicpatch.CreateTwoWayMergePatch(orginal, modified, nodeInfo)
		if err != nil {
			return err
		}

		_, err = c.kubeclientset.CoreV1().Nodes().Patch(nodeName, types.StrategicMergePatchType, patch)
		if err != nil {
			glog.Errorf("update node err: %v", err)
			return err
		}
	}

	return nil
}

// enqueueResourceClass takes a ResourceClass resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than ResourceClass.
func (c *Controller) enqueueResourceClass(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	glog.Infof("enqueueResourceClass key: %s", key)
	c.workqueue.AddRateLimited(key)
}
