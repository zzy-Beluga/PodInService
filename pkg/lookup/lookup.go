package lookup

import (
	"PodInService/common"
	"context"
	"fmt"
	"os"

	svcv1 "k8s.io/api/core/v1"
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
}

func getServiceForPod(podName, namespace string) (string, error) {
	// create the pod client
	podClient = clientset.CoreV1().Pods(namespace)

	// create the service client
	serviceClient = clientset.CoreV1().Services(common.SvcNamespace)
	ctx := context.Background()
	// get the pod object
	pod, err := podClient.Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// get the labels for the pod
	podLabels := pod.GetLabels()
	fmt.Println(podLabels)
	appLabels := make(map[string]string)
	serviceList := &svcv1.ServiceList{}
	for k, v := range podLabels {
		appLabels[k] = v
		labelSelector := labels.Set(appLabels).AsSelector()

		// list all services with the same labels as the pod
		serviceList, err = serviceClient.List(ctx, metav1.ListOptions{LabelSelector: labelSelector.String()})
		if err != nil {
			return "", err
		}
		if len(serviceList.Items) == 0 {
			continue
		}
	}
	if len(serviceList.Items) == 0 {
		return "", fmt.Errorf("no services found for pod %s", podName)
	}

	// return the name of the first service in the list
	fmt.Println(len(serviceList.Items))
	return serviceList.Items[0].GetName(), nil
}

func Find(namespace, podname string) (string, error) {
	svc, err := getServiceForPod(podname, namespace)
	if err != nil {
		return "", err
	}
	return svc, nil
}
