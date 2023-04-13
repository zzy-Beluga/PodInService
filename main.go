package main

import (
	"PodInService/pkg/lookup"
	"fmt"
)

func main() {
	svc, _ := lookup.Find()
	fmt.Println(svc)
}
