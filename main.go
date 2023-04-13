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
	ans, err := lookup.Find(podname, namespace)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(ans)
}
