package controller

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

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

func getServiceForPod(ctx context.Context,podName string) (string, error) {
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
