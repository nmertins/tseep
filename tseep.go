package main

import (
	"fmt"
	"time"
)

const (
	MainLoopPeriod = 10 * time.Second
)

func main() {
	for {
		fmt.Printf("Hello, world!\n")
		time.Sleep(MainLoopPeriod)
	}
}
