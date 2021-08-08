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

func getCurrentConnections() (string, error) {
	data, err := ioutil.ReadFile("/proc/net/tcp")
	if err != nil {
		return "", err
	}

	connectionsRaw := string(data)

	return connectionsRaw, nil
}

func main() {
	currentConnections := tseep.CurrectConnections{}

	for {
		connectionsRaw, err := getCurrentConnections()
		if err != nil {
			fmt.Printf("Error reading /proc/net/tcp: %s\n", err.Error())
		} else {
			tcpConnections := tseep.ParseListOfConnections(connectionsRaw)
			newConnections := currentConnections.Update(tcpConnections)
			tseep.PrintNewConnections(os.Stdout, time.Now(), newConnections)
		}
		time.Sleep(MainLoopPeriod)
	}
}
