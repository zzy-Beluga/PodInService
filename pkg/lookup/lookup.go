package lookup

import (
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
	serviceList := &svcv1.ServiceList{}

	// traverse all labels to find all svc matches
	for k, v := range podLabels {

		// cache the current label and create a selector for svc matching
		matchLabels[k] = v
		fmt.Println(matchLabels)
		labelSelector := labels.Set(matchLabels).AsSelector()

		// list all services with the same labels as the pod for the current labels
		fmt.Println(labelSelector.String())
		sl, err := serviceClient.List(ctx, metav1.ListOptions{LabelSelector: labelSelector.String()})
		fmt.Println(sl)
		if err != nil {
			return nil, err
		}

		delete(matchLabels, k)

		//no matching svc for this label
		if len(sl.Items) == 0 {
			fmt.Printf("The Label %v has no matched service \n", v)
			continue
		}

		// append the svc to serviceList
		serviceList.Items = append(serviceList.Items, sl.Items...)
	}

	// no matched svc
	if len(serviceList.Items) == 0 {
		return nil, fmt.Errorf("no services found for pod %s", podName)
	}

	// store results into map
	res := make(map[string]string)
	for _, i := range serviceList.Items {
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
