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
	"encoding/json"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	v1 "github.com/trainyao/playground/k8s/crd-controller/handmade/pkg/apis/handmade/v1"
	clientset "github.com/trainyao/playground/k8s/crd-controller/handmade/pkg/client/clientset/versioned"
	handmadeScheme "github.com/trainyao/playground/k8s/crd-controller/handmade/pkg/client/clientset/versioned/scheme"
	informers "github.com/trainyao/playground/k8s/crd-controller/handmade/pkg/client/informers/externalversions/handmade/v1"
	listers "github.com/trainyao/playground/k8s/crd-controller/handmade/pkg/client/listers/handmade/v1"
)

const (
	controllerAgentName = "goddessmoment-controller"
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
	buyFoodTimeCostMap = map[string]time.Duration{
		"珍珠奶茶": time.Second,
		"麻辣烫":  5 * time.Second,
		"螺蛳粉":  10 * time.Second,
	}
)

// Controller is the controller implementation for Foo resources
type Controller struct {
	name string
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// handmadeClientset is a clientset for our own API group
	handmadeClientset clientset.Interface

	gmLister listers.GoddessMomentLister
	gmSynced cache.InformerSynced

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
	handmadeClientset clientset.Interface,
	goddessMomentInformer informers.GoddessMomentInformer) *Controller {

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
		name:              name,
		kubeclientset:     kubeclientset,
		handmadeClientset: handmadeClientset,
		gmLister:          goddessMomentInformer.Lister(),
		gmSynced:          goddessMomentInformer.Informer().HasSynced,
		workqueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "GoddessMoments"),
		recorder:          recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Foo resources change
	goddessMomentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueGoddessMoment,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueGoddessMoment(new)
		},
		DeleteFunc: controller.enqueueGoddessMoment,
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
	if ok := cache.WaitForCacheSync(stopCh, c.gmSynced); !ok {
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

	// Get the Foo resource with this namespace/name
	gm, err := c.gmLister.GoddessMoments(namespace).Get(name)
	if err != nil {
		// The GoddessMoment resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("gm '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	focusFood, ok := foodFocus[c.name]
	if !ok {
		utilruntime.HandleError(fmt.Errorf("%s 朋友圈找不到 %s 的关注食物", gm.Name, focusFood))
		return nil
	}
	klog.Infof("%s 的关注食物是 %s", c.name, focusFood)

	foundFocusFood := false
	for _, food := range gm.Spec.FoodDemand {
		if food.Name == focusFood {
			foundFocusFood = true
			break
		}
	}

	if !foundFocusFood {
		klog.Infof("%s 关注食物 %s 没有在朋友圈 %s 发现", c.name, focusFood, key)
		return nil
	}

	var foodStatus *v1.FoodDemandStatus
	var index = -1
	for i, foodDemandStatus := range gm.Status.FoodDemand {
		if foodDemandStatus.Name == focusFood {
			foodStatus = &gm.Status.FoodDemand[i]
			index = i
		}
	}
	if foodStatus == nil {
		utilruntime.HandleError(fmt.Errorf("关注食物 %s 在女神朋友圈出现, 但是在status里没有找到", focusFood))
		return nil
	}

	klog.Infof("朋友圈 %s %s 关注的食物 %s 状态是 %s", gm.Name, c.name, focusFood, foodStatus.Status)

	if foodStatus.Status == v1.FoodStatusPending {
		// 处理女神朋友圈里关注食物还未被人认领的情况
		err = c.handlePending(gm, index)
	} else {
		// 处理关注食物已经有人认领的情况

		// 认领的人不是我
		if foodStatus.ClaimBy != c.name {
			klog.Infof("朋友圈 %s 食物 %s 已经被人 %s 认领了, 我来的太迟了", gm.Name, focusFood, foodStatus.ClaimBy)
			return nil
		}

		if foodStatus.Status == v1.FoodStatusPendingArrival {
			// 去购买食物
			err = c.handlePendingArrival(gm, index)
		}

		//  状态是已到达, 本条朋友圈任务已完成
		klog.Infof("朋友圈 %s 食物 %s 我已经认领并送到了, 可以不用关注了", gm.Name, focusFood)
	}

	if err != nil {
		err = fmt.Errorf("处理朋友圈 %s 更新失败, err: %s", gm.Name, err.Error())
		utilruntime.HandleError(err)
		return err
	}

	c.recorder.Event(gm, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) handlePending(gm *v1.GoddessMoment, index int) error {
	focusFood := gm.Status.FoodDemand[index].Name
	klog.Infof("开始告诉女神 %s 我来送", focusFood)
	if c.name == "小杰" {
		// 展示patch更新
		data := `[{"op":"replace","path":"/status/foodDemand","value":%s}]`
		foodDemandStatusCopy := gm.Status.DeepCopy()
		foodDemandStatusCopy.FoodDemand[index].Status = v1.FoodStatusPendingArrival
		foodDemandStatusCopy.FoodDemand[index].ClaimTime = metav1.Now()
		foodDemandStatusCopy.FoodDemand[index].ClaimBy = c.name

		valueJSON, err := json.Marshal(&foodDemandStatusCopy.FoodDemand)
		if err != nil {
			return fmt.Errorf("marshal foodDemandStatusCopy.FoodDemand 失败, err: %s", err.Error())
		}

		dataJSON := fmt.Sprintf(data, valueJSON)
		_, err = c.handmadeClientset.HandmadeV1().GoddessMoments(gm.Namespace).Patch(context.TODO(), gm.Name, types.JSONPatchType, []byte(dataJSON), metav1.PatchOptions{})
		if err != nil {
			return fmt.Errorf("patch更新朋友圈告诉女神 %s 我来送 %s 失败, err: %s, dataJSON: %s", gm.Name, focusFood, err.Error(), dataJSON)
		}
	} else {
		// 使用update更新
		gmCopy := gm.DeepCopy()
		gmCopy.Status.FoodDemand[index].Status = v1.FoodStatusPendingArrival
		gmCopy.Status.FoodDemand[index].ClaimTime = metav1.Now()
		gmCopy.Status.FoodDemand[index].ClaimBy = c.name

		_, err := c.handmadeClientset.HandmadeV1().GoddessMoments(gm.Namespace).Update(context.TODO(), gmCopy, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("更新朋友圈 %s 告诉女神 %s 我来送失败, err: %s", gm.Name, focusFood, err.Error())
		}
	}

	klog.Infof("更新朋友圈 %s 告诉女神送来 %s 成功", gm.Name, focusFood)
	return nil
}

func (c *Controller) handlePendingArrival(gm *v1.GoddessMoment, index int) error {
	focusFood := gm.Status.FoodDemand[index].Name
	buyFoodTimeCost := buyFoodTimeCostMap[focusFood]
	time.Sleep(buyFoodTimeCost)
	klog.Infof("为朋友圈 %s 购买 %s 需要 %d 秒, 正在前往购买", gm.Name, focusFood, buyFoodTimeCost/time.Second)

	gmCopy := gm.DeepCopy()
	gmCopy.Status.FoodDemand[index].Status = v1.FoodStatusArrived
	gmCopy.Status.FoodDemand[index].ArrivalTime = metav1.Now()
	_, err := c.handmadeClientset.HandmadeV1().GoddessMoments(gm.Namespace).Update(context.TODO(), gmCopy, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("更新朋友圈 %s 告诉女神送来 %s 失败, err: %s", gm.Name, focusFood, err.Error())
	}
	klog.Infof("更新朋友圈 %s 告诉女神送来 %s 成功", gm.Name, focusFood)
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
