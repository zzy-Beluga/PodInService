package main

import (
	"PodInService/pkg/controller"
	"fmt"
)

func main() {
	svc, _ := controller.Start()
	fmt.Println(svc)
}
