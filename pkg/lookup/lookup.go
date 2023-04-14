package lookup

import (
	"context"
	"fmt"
	"os"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// filter svc by spec.selector
func serviceFilter(svclist *corev1.ServiceList, matchlabels map[string]string) ([]corev1.Service, error) {
	for _, i := range svclist.Items {
		if reflect.DeepEqual(i.Spec.Selector, matchlabels) {
			s := []corev1.Service{i}
			return s, nil
		}
	}
	return nil, nil
}

func getServiceForPod(podName, namespace string) (map[string]string, error) {

	ctx := context.Background()

	// create the pod client
	podClient = clientset.CoreV1().Pods(namespace)

	// create the service client which have access all namespaces
	serviceClient = clientset.CoreV1().Services("")

	// get the pod object
	pod, err := podClient.Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// get the labels for the pod
	podLabels := pod.GetLabels()
	// init the serviceList to contain all matchable svc and matchLabels to cache each matching labels
	matchLabels := make(map[string]string)
	serviceList := make([]corev1.Service, 0)

	// traverse all labels to find all svc matches
	for k, v := range podLabels {
		// cache the current label and create a selector for svc matching
		matchLabels[k] = v

		// list all services with the same labels as the pod for the current labels
		fullist, err := serviceClient.List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		sl, err := serviceFilter(fullist, matchLabels)
		if err != nil {
			return nil, err
		}

		delete(matchLabels, k)

		// append the svc to serviceList
		serviceList = append(serviceList, sl...)
	}

	// no matched svc
	if len(serviceList) == 0 {
		return nil, fmt.Errorf("no services found for pod %s \n", podName)
	}

	// store results into map
	res := make(map[string]string)
	for _, i := range serviceList {
		res[i.GetName()] = i.GetNamespace()
	}

	// return all matched svc
	return res, nil
}

func Find(podname, namespace string) (map[string]string, error) {
	svc, err := getServiceForPod(podname, namespace)
	if err != nil {
		return nil, err
	}
	return svc, nil
}
