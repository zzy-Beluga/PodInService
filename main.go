package main

import (
	"PodInService/pkg/lookup"
	"flag"
	"fmt"
)

func main() {
	var namespace string
	var podname string
	flag.StringVar(&namespace, "n", "default", "namespace")
	flag.StringVar(&podname, "p", "common-nginx-vm-574bb74b46-86m9d", "podname")
	flag.Parse()
	res, err := lookup.Find(podname, namespace)
	if err != nil {
		fmt.Print(err)
	}

	// print svc results
	fmt.Println("")
	for k, v := range res {
		fmt.Printf("The Pod Matches the Following Service:\nService Name: %s \nService NameSapce: %s \n", k, v)
	}
}
