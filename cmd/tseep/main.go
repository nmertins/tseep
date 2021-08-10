package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/nmertins/tseep"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	MainLoopPeriod = 10 * time.Second
)

func main() {
	currentConnections := tseep.CurrectConnections{
		Source: "/proc/net/tcp",
	}

	connectionsCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "tseep_new_connections",
		Help: "The total number of TCP connections received",
	})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	for {
		newConnections, portScans, err := currentConnections.Update()
		if err != nil {
			fmt.Printf("Error reading TCP connections: %s\n", err.Error())
		} else {
			tseep.PrintNewConnections(os.Stdout, newConnections)
			tseep.PrintPortScans(os.Stdout, portScans)
			for _, portScan := range portScans {
				err := tseep.BlockPortScanSource(portScan)
				if err != nil {
					fmt.Printf("Attempted to block remote address but failed: %s\n", err.Error())
				}
			}

			connectionsCounter.Add(float64(len(newConnections)))
		}
		time.Sleep(MainLoopPeriod)
	}
}
