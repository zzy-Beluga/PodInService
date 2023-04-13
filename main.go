package main

import (
	"PodInService/pkg/lookup"
	"flag"
	"fmt"
)

func main() {
	var namespace string
	var podname string
	flag.StringVar(&namespace, "n", "kube-system", "please input pod namespace")
	flag.StringVar(&podname, "p", "coredns-565d847f94-mzqjr", "please input pod podname")
	flag.Parse()
	ans, _ := lookup.Find(namespace, podname)
	fmt.Println(ans)
}
