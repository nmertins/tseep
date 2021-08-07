package main

import (
	"fmt"
	"io/ioutil"
	"time"
)

const (
	MainLoopPeriod = 10 * time.Second
)

func main() {
	for {
		time.Sleep(MainLoopPeriod)

		data, err := ioutil.ReadFile("/proc/net/tcp")
		if err != nil {
			fmt.Printf("Error reading /proc/net/tcp: %s\n", err.Error())
			continue
		}

		connections := string(data)

		fmt.Printf("%s", connections)
	}
}
