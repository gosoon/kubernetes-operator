package kuberesource

import (
	"os"

	"github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned/scheme"
	"github.com/gosoon/kubernetes-operator/pkg/controller"

	"github.com/gosoon/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

func NewResourceLock(kubeClient *kubernetes.Clientset) (resourcelock.Interface, error) {
	// init eventRecorder
	eventBroadcaster := record.NewBroadcaster()
	eventRecorder := eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: controller.ComponentName})
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})

	// init host identity
	id, err := os.Hostname()
	if err != nil {
		glog.Errorf("get hostname error: %v", err)
		return nil, err
	}
	id = id + "_" + string(uuid.NewUUID())

	rl, err := resourcelock.New("endpoints",
		"kube-system",
		controller.ComponentName,
		kubeClient.CoreV1(),
		kubeClient.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: eventRecorder,
		})

	if err != nil {
		glog.Errorf("error creating lock: %v", err)
		return nil, err
	}
	return rl, nil
}
