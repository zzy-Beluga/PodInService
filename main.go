package main

import (
	"PodInService/pkg/lookup"
	"flag"
	"fmt"
)

func main() {

	// cmdline flags, n for namespace, p for podnames
	var namespace string
	var podname string
	flag.StringVar(&namespace, "n", "default", "namespace")
	flag.StringVar(&podname, "p", "", "podname")
	flag.Parse()

	//look for services
	res, err := lookup.Find(podname, namespace)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print svc results
	fmt.Println("The Pod Matches the Following Service:")
	for k, v := range res {
		fmt.Printf("Service Name: %s Service NameSapce: %s \n", k, v)
	}
}
