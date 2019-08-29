package main

import (
	"fmt"
	"github.com/code-sleuth/powerplant-monitoring-simulator/src/distributed/coordinator"
)

func main()  {
	ql := coordinator.NewQueueListener()
	go ql.ListenForNewSource()

	var a string
	fmt.Scanln(&a)
}
