package main

import (
	"PodInService/pkg/lookup"
	"flag"
	"fmt"
)

func main() {
	var namespace string
	var podname string
	flag.StringVar(&namespace, "Namespace", "kube-system", "please input pod namespace")
	flag.StringVar(&podname, "Podname", "kube-dns", "please input pod podname")
	flag.Parse()
	fmt.Println(lookup.Find(namespace, podname))
}
