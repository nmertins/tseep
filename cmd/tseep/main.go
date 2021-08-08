package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nmertins/tseep"
)

const (
	MainLoopPeriod = 10 * time.Second
)

func main() {
	currentConnections := tseep.CurrectConnections{
		Source: "/proc/net/tcp",
	}

	for {
		newConnections, err := currentConnections.Update()
		if err != nil {
			fmt.Printf("Error reading TCP connections: %s\n", err.Error())
		} else {
			tseep.PrintNewConnections(os.Stdout, newConnections)
		}
		time.Sleep(MainLoopPeriod)
	}
}
