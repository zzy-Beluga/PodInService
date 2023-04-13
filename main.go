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
	flag.StringVar(&podname, "p", "", "podname")
	flag.Parse()
	res, err := lookup.Find(podname, namespace)
	if err != nil {
		fmt.Print(err)
	}

	// print svc results
	fmt.Println("----The Pod Matches the Following Services Resource----")
	for k, v := range res {
		fmt.Printf("Service Name: %v, Service NameSapce: %v \n", k, v)
	}
}
