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
)

const controllerAgentName = "ecs-controller"

//type KubernetesOperatorPhase string

const (
	// kubernetes cluster phase,"None,Creating,Running,Failed,Scaling"
	// Active is the create kubernetes job is running
	None        ecsv1.KubernetesOperatorPhase = ""
	Creating    ecsv1.KubernetesOperatorPhase = "Creating"
	Running     ecsv1.KubernetesOperatorPhase = "Running"
	Failed      ecsv1.KubernetesOperatorPhase = "Failed"
	Scaling     ecsv1.KubernetesOperatorPhase = "Scaling"
	Terminating ecsv1.KubernetesOperatorPhase = "Terminating"

	// SuccessSynced is used as part of the Event 'reason' when a KubernetesCluster is synced
	SuccessSynced = "Synced"

	// MessageResourceSynced is the message used for an Event fired when a KubernetesCluster
	// is synced successfully
	MessageResourceSynced = "ecs synced successfully"

	SuccessCreated = "SuccessCreated"
	FailCreated    = "FailCreated"

	CreateKubernetesClusterJobSuccess = "create kubernetes cluster job success"
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
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

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
	switch {
	case errors.IsNotFound(err):
		// The KubernetesCluster resource may no longer exist, in which case we stop
		// processing.
		err = c.processKubernetesClusterDeletion(key)
	case err != nil:
		runtime.HandleError(fmt.Errorf("Unable to retrieve service %v from store: %v", key, err))
	default:
		err = c.processKubernetesClusterCreateOrUpdate(kubernetesCluster, key)
	}
	return err
}

func (c *Controller) processKubernetesClusterCreateOrUpdate(kubernetesCluster *ecsv1.KubernetesCluster, key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	fmt.Println("%+v", *kubernetesCluster)

	glog.Infof("[Neutron] Try to process kubernetesCluster: %#v ...", kubernetesCluster)

	// FIX ME: Do diff().
	//
	// actualKubernetesCluster, exists := neutron.Get(namespace, name)
	//
	// if !exists {
	// 	neutron.Create(namespace, name)
	// } else if !reflect.DeepEqual(actualKubernetesCluster, kubernetesCluster) {
	// 	neutron.Update(namespace, name)
	// }

	// TODO: handle previous events when operator restart
	// when started,diff current and expect obj
	//if !reflect.DeepEqual(kubernetesCluster, currentKubernetesCluster) {
	//case "SCALE":
	//case "CREATE":
	//}

	switch kubernetesCluster.Status.Phase {
	// phase is "" express create new kubernetesCluster
	case "":
		// update phase
		kubernetesCluster.Status.Phase = Creating
		_, err := c.kubernetesClusterClientset.EcsV1().KubernetesClusters(namespace).UpdateStatus(kubernetesCluster)
		if err != nil {
			glog.Errorf("update status failed with:%v", err)
			return err
		}
		// create kubernetes cluster
		batchJob := newCreateKubernetesClusterBatchJob(kubernetesCluster)

		fmt.Println("%+v", batchJob)

		_, err = c.kubeclientset.BatchV1().Jobs(namespace).Create(batchJob)
		if err != nil {
			createKubernetesClusterJobFailed := fmt.Sprintf("create kubernetes cluster job failed with:%v", err)
			c.recorder.Event(kubernetesCluster, corev1.EventTypeNormal, FailCreated, createKubernetesClusterJobFailed)
			glog.Errorf("create %s/%s kubernetes cluster job failed with:%v", namespace, name, err)
			return err
		}
		c.recorder.Event(kubernetesCluster, corev1.EventTypeNormal, SuccessCreated, CreateKubernetesClusterJobSuccess)

		// callback controller to ensure create success

	case "CREATING":
		// check job

	// update or delete
	default:

	}

	c.recorder.Event(kubernetesCluster, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) processKubernetesClusterDeletion(key string) error {

	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	glog.Warningf("KubernetesCluster: %s/%s does not exist in local cache, will delete it from Neutron ...",
		namespace, name)

	glog.Infof("[Neutron] Deleting kubernetesCluster: %s/%s ...", namespace, name)

	// FIX ME: call Neutron API to delete this kubernetesCluster by name.
	//
	// neutron.Delete(namespace, name)
	// delete job and kubernetes cluster
	// update DeletionTimestamp
	deleteClusterBatchJob := newDeleteKubernetesClusterBatchJob(name, namespace)
	_, err = c.kubeclientset.BatchV1().Jobs(namespace).Create(deleteClusterBatchJob)
	if err != nil {
		glog.Errorf("delete %s/%s kubernetes cluster job failed with:%v", namespace, name, err)
		return err
	}

	// call back and delete deleteJob and createJob
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
