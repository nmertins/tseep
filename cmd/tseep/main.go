package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/nmertins/tseep"
)

const (
	MainLoopPeriod = 10 * time.Second
)

func main() {
	currentConnections := tseep.CurrectConnections{}

	for {
		time.Sleep(MainLoopPeriod)

		data, err := ioutil.ReadFile("/proc/net/tcp")
		if err != nil {
			fmt.Printf("Error reading /proc/net/tcp: %s\n", err.Error())
			continue
		}

		connectionsRaw := string(data)
		tcpConnections := tseep.ParseListOfConnections(connectionsRaw)
		newConnections := currentConnections.Update(tcpConnections)

		tseep.PrintNewConnections(os.Stdout, time.Now(), newConnections)
	}
}
