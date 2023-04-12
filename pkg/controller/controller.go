package controller

import (
	"context"
	"fmt"

	"github.com/bitnami-labs/kubewatch/pkg/handlers"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

// Event indicate the informerEvent
type Event struct {
	key          string
	eventType    string
	namespace    string
	resourceType string
}

// Controller object
type Controller struct {
	logger       *logrus.Entry
	clientset    kubernetes.Interface
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	eventHandler handlers.Handler
}

var clientset *kubernetes.Clientset
var podClient v1.PodInterface
var serviceClient v1.ServiceInterface

func init() {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// create the pod client
	podClient = clientset.CoreV1().Pods("default")

	// create the service client
	serviceClient = clientset.CoreV1().Services("default")
}

func getServiceForPod(podName string) (string, error) {
	ctx := context.Background()
	// get the pod object
	pod, err := podClient.Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// get the labels for the pod
	podLabels := pod.GetLabels()

	// create a label selector for the pod's labels
	labelSelector := labels.Set(podLabels).AsSelector()

	// list all services with the same labels as the pod
	serviceList, err := serviceClient.List(ctx, metav1.ListOptions{LabelSelector: labelSelector.String()})
	if err != nil {
		return "", err
	}

	// if there are no services, return an error
	if len(serviceList.Items) == 0 {
		return "", fmt.Errorf("no services found for pod %s", podName)
	}

	// return the name of the first service in the list
	return serviceList.Items[0].GetName(), nil
}

func Start() {
	getServiceForPod("mypod")
}
