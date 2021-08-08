package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	MainLoopPeriod = 10 * time.Second
)

func main() {
	currentConnections := CurrectConnections{}

	for {
		time.Sleep(MainLoopPeriod)

		data, err := ioutil.ReadFile("/proc/net/tcp")
		if err != nil {
			fmt.Printf("Error reading /proc/net/tcp: %s\n", err.Error())
			continue
		}

		connectionsRaw := string(data)
		tcpConnections := ParseListOfConnections(connectionsRaw)
		newConnections := currentConnections.Update(tcpConnections)

		PrintNewConnections(os.Stdout, time.Now(), newConnections)
	}
}
