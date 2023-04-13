package controller

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// Event indicate the informerEvent
/*
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
*/

var clientset *kubernetes.Clientset
var podClient v1.PodInterface
var serviceClient v1.ServiceInterface

func init() {
	// use the current context in kubeconfig
	path := os.Getenv("HOME") + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// create the pod client
	podClient = clientset.CoreV1().Pods("kube-system")

	// create the service client
	serviceClient = clientset.CoreV1().Services("kube-system")
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
	appLabels := make(map[string]string)
	appLabels["k8s-app"] = podLabels["k8s-app"]
	fmt.Println(podLabels)
	// create a label selector for the pod's labels
	labelSelector := labels.Set(appLabels).AsSelector()

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
	fmt.Println(len(serviceList.Items))
	return serviceList.Items[0].GetName(), nil
}

func Start() (string, error) {
	svc, err := getServiceForPod("coredns-565d847f94-mzqjr")
	if err != nil {
		return "", err
	}
	return svc, nil
}
