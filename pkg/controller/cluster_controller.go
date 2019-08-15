/*
 * Copyright 2019 gosoon.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controller

import (
	"fmt"
	"time"

	"github.com/gosoon/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	clientset "github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned"
	ecsscheme "github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned/scheme"
	informers "github.com/gosoon/kubernetes-operator/pkg/client/informers/externalversions/ecs/v1"
	listers "github.com/gosoon/kubernetes-operator/pkg/client/listers/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"
)

const (
	ComponentName = "kubernetes-operator"

	// SuccessSynced is used as part of the Event 'reason' when a KubernetesCluster is synced
	SuccessSynced = "Synced"

	DeleteJobLabelCreated = "created"
)

// Controller is the controller implementation for KubernetesCluster resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// ecsclientset is a clientset for our own API group
	kubernetesClusterClientset clientset.Interface

	kubernetesClusterLister listers.KubernetesClusterLister
	kubernetesClusterSynced cache.InformerSynced

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

// NewController returns a new kubernetesCluster controller
func NewController(
	kubeclientset kubernetes.Interface,
	kubernetesClusterClientset clientset.Interface,
	kubernetesClusterInformer informers.KubernetesClusterInformer) *Controller {

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(ecsscheme.AddToScheme(scheme.Scheme))
	glog.Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: ComponentName})

	controller := &Controller{
		kubeclientset:              kubeclientset,
		kubernetesClusterClientset: kubernetesClusterClientset,
		kubernetesClusterLister:    kubernetesClusterInformer.Lister(),
		kubernetesClusterSynced:    kubernetesClusterInformer.Informer().HasSynced,
		workqueue:                  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "KubernetesClusters"),
		recorder:                   recorder,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when kubernetesCluster resources change
	kubernetesClusterInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueKubernetesCluster,
		UpdateFunc: func(old, new interface{}) {
			oldKubernetesCluster := old.(*ecsv1.KubernetesCluster)
			newKubernetesCluster := new.(*ecsv1.KubernetesCluster)
			if oldKubernetesCluster.ResourceVersion == newKubernetesCluster.ResourceVersion {
				// Periodic resync will send update events for all known KubernetesClusters.
				// Two different versions of the same KubernetesCluster will always have different RVs.
				return
			}
			controller.enqueueKubernetesCluster(new)
		},
		DeleteFunc: controller.enqueueKubernetesClusterForDelete,
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
	glog.Info("Starting KubernetesCluster control loop")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.kubernetesClusterSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	// Launch two workers to process KubernetesCluster resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
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
		// KubernetesCluster resource to be synced.
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
// converge the two. It then updates the Status block of the KubernetesCluster resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	if len(namespace) == 0 || len(name) == 0 {
		glog.Errorf("invalid key %q: either namespace or name is missing", key)
		return err
	}

	// Get the KubernetesCluster resource with this namespace/name
	kubernetesCluster, err := c.kubernetesClusterLister.KubernetesClusters(namespace).Get(name)
	//fmt.Printf("%+v\n", kubernetesCluster)
	switch {
	case errors.IsNotFound(err):
		// The KubernetesCluster resource may no longer exist, in which case we stop
		// processing.
	case err != nil:
		runtime.HandleError(fmt.Errorf("Unable to retrieve service %v from store: %v", key, err))
	default:
		err = c.processKubernetesClusterCreateOrUpdate(kubernetesCluster)
	}
	return err
}

// processKubernetesClusterCreateOrUpdate is handle all status in kubernetesCluster.
func (c *Controller) processKubernetesClusterCreateOrUpdate(kubernetesCluster *ecsv1.KubernetesCluster) error {
	operation := kubernetesCluster.Annotations[enum.Operation]

	switch kubernetesCluster.Status.Phase {
	// phase is "" express create new kubernetesCluster
	case enum.New:
		// precheck must be performed before New, KubeScalingUp, KubeScalingDown operate
		return c.processOperatePrecheck(kubernetesCluster)

	// TODO: callback controller to ensure create success
	case enum.Creating, enum.Running, enum.Scaling:
		if kubernetesCluster.DeletionTimestamp != nil {
			if operation == enum.KubeTerminating {
				return c.processClusterTerminating(kubernetesCluster)
			}
			return nil
		}
		// annotation
		switch operation {
		// precheck must be performed before New, KubeScalingUp, KubeScalingDown operate
		case enum.KubeScalingUp, enum.KubeScalingDown:
			return c.processOperatePrecheck(kubernetesCluster)

		// other operate is only update status.phase
		case enum.KubeCreating:
			return c.processOperateCreating(kubernetesCluster)
		case enum.KubeCreateFailed:
			return c.processOperateFailed(kubernetesCluster)
		case enum.KubeCreateFinished:
			return c.processOperateFinished(kubernetesCluster)
		case enum.KubeScaleUpFailed:
			return c.processOperateFailed(kubernetesCluster)
		case enum.KubeScaleUpFinished:
			return c.processOperateFinished(kubernetesCluster)
		case enum.KubeScaleDownFailed:
			return c.processOperateFailed(kubernetesCluster)
		case enum.KubeScaleDownFinished:
			return c.processOperateFinished(kubernetesCluster)
		default:

		}

	// Terminating
	case enum.Terminating:
		if operation == enum.KubeTerminateFailed {
			return c.processOperateFailed(kubernetesCluster)
		}

	// Failed
	case enum.Failed:
		// delete retry or last job failed
		if kubernetesCluster.DeletionTimestamp != nil {
			if operation == enum.KubeTerminating {
				return c.processClusterTerminating(kubernetesCluster)
			}
		}
		// retry operate
		switch operation {
		case enum.KubeCreating:
			return c.processOperateNew(kubernetesCluster)
		case enum.KubeScalingUp:
			return c.processClusterScaleUp(kubernetesCluster)
		case enum.KubeScalingDown:
			return c.processClusterScaleDown(kubernetesCluster)
		}

	// Prechecking
	case enum.Prechecking:
		return c.processClusterPrecheck(kubernetesCluster)

	default:

	}

	return nil
}

// enqueueKubernetesCluster takes a KubernetesCluster resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than KubernetesCluster.
func (c *Controller) enqueueKubernetesCluster(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

// enqueueKubernetesClusterForDelete takes a deleted KubernetesCluster resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than KubernetesCluster.
func (c *Controller) enqueueKubernetesClusterForDelete(obj interface{}) {
	var key string
	var err error
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}
