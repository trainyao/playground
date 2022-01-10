/*
Copyright 2017 The Kubernetes Authors.

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
	"context"
	"fmt"
	v1 "github.com/trainyao/playground/k8s/crd-controller/handmade/pkg/apis/handmade/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	handmadeClientset "github.com/trainyao/playground/k8s/crd-controller/handmade/pkg/client/clientset/versioned"
	handmadeScheme "github.com/trainyao/playground/k8s/crd-controller/handmade/pkg/client/clientset/versioned/scheme"
	handmadeInformers "github.com/trainyao/playground/k8s/crd-controller/handmade/pkg/client/informers/externalversions/handmade/v1"
	handmadeLister "github.com/trainyao/playground/k8s/crd-controller/handmade/pkg/client/listers/handmade/v1"
)

const controllerAgentName = "sample-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	// MessageResourceSynced is the message used for an Event fired when a Foo
	// is synced successfully
	MessageResourceSynced = "Foo synced successfully"
)

var (
	foodFocus = map[string]string{
		"小明": "珍珠奶茶",
		"小王": "麻辣烫",
		"小杰": "螺蛳粉",
	}
)

// Controller is the controller implementation for Foo resources
type Controller struct {
	name string
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// handmadeClientset is a clientset for our own API group
	handmadeClientset handmadeClientset.Interface

	deploymentsLister    appslisters.DeploymentLister
	deploymentsSynced    cache.InformerSynced
	goddessMomentsLister handmadeLister.GoddessMomentLister
	goddessMomentsSynced cache.InformerSynced

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

// NewController returns a new sample controller
func NewController(
	name string,
	kubeclientset kubernetes.Interface,
	sampleclientset handmadeClientset.Interface,
	goddessMomentsInformer handmadeInformers.GoddessMomentInformer) *Controller {

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(handmadeScheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		name:                 name,
		kubeclientset:        kubeclientset,
		handmadeClientset:    sampleclientset,
		goddessMomentsLister: goddessMomentsInformer.Lister(),
		goddessMomentsSynced: goddessMomentsInformer.Informer().HasSynced,
		workqueue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Foos"),
		recorder:             recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Foo resources change
	goddessMomentsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueGoddessMoment,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueGoddessMoment(new)
		},
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting Foo controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.goddessMomentsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process Foo resources
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

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
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Foo resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Foo resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the GoddessMoment resource with this namespace/name
	goddessMoment, err := c.goddessMomentsLister.GoddessMoments(namespace).Get(name)
	if err != nil {
		// The GoddessMoment resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("goddessMoment '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	focusFood, ok := foodFocus[c.name]
	if !ok {
		utilruntime.HandleError(fmt.Errorf("can't found focus food for %s", c.name))
		return nil
	}

	foundFocusFood := false
	for _, food := range goddessMoment.Spec.FoodDemand {
		if food.Name == focusFood {
			foundFocusFood = true
			break
		}
	}

	if !foundFocusFood {
		klog.Infof("focus food %s for %s is not in goddessMoment %s", focusFood, c.name, key)
		return nil
	}

	var foodStatus *v1.FoodDemandStatus
	var index = -1
	for i, foodDemandStatus := range goddessMoment.Status.FoodDemand {
		if foodDemandStatus.Name == focusFood {
			foodStatus = &goddessMoment.Status.FoodDemand[i]
			index = i
			break
		}
	}

	if foodStatus == nil {
		utilruntime.HandleError(fmt.Errorf("found focus food %s in spec, but not found in status", focusFood))
		return nil
	}

	if foodStatus.Status != v1.FoodStatusPending {
		klog.Infof("food %s status not pending: %s", focusFood, foodStatus.Status)
		return nil
	}
	if c.name == "小杰" {
		data := `{"op":"replace","path":"/status/foodDemand","value":%s}`
		foodDemandStatusCopy := goddessMoment.Status.DeepCopy()
		foodDemandStatusCopy.FoodDemand[index].Status = v1.FoodStatusPendingArrival
		foodDemandStatusCopy.FoodDemand[index].ClaimTime = metav1.Now()
		foodDemandStatusCopy.FoodDemand[index].ClaimBy = c.name
		valueJson, err := json.Marshal(&foodDemandStatusCopy.FoodDemand)
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("error when marshal foodDemand status, err: %s", err.Error()))
			return nil
		}

		dataJson := fmt.Sprintf(data, valueJson)
		_, err = c.handmadeClientset.HandmadeV1().GoddessMoments(namespace).Patch(context.TODO(), name,
			types.JSONPatchType, []byte(dataJson), metav1.PatchOptions{})
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("patch goddessMoment %s failed, data: %s, err: %s", name, dataJson, err.Error()))
			return nil
		}

		return nil
	}

	momentCopy := goddessMoment.DeepCopy()
	momentCopy.Status.FoodDemand[index].Status = v1.FoodStatusPendingArrival
	momentCopy.Status.FoodDemand[index].ClaimTime = metav1.Now()
	momentCopy.Status.FoodDemand[index].ClaimBy = c.name

	_, err = c.handmadeClientset.HandmadeV1().GoddessMoments(namespace).Update(context.TODO(), momentCopy, metav1.UpdateOptions{})
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("update goddessMoment %s failed, err: %s", name, err.Error()))
		return nil
	}

	c.recorder.Event(goddessMoment, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// enqueueGoddessMoment takes a Foo resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Foo.
func (c *Controller) enqueueGoddessMoment(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}