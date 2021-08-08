package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

const (
	MainLoopPeriod = 10 * time.Second
)

func PrintNewConnections(writer io.Writer, t time.Time, newConnections []TcpConnection) {
	timestamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	for _, connection := range newConnections {
		fmt.Fprintf(writer, "%s: New connection: %s:%d -> %s:%d\n", timestamp, connection.remoteAddress, connection.remotePort, connection.localAddress, connection.localPort)
	}
}

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
