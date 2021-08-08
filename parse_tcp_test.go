package tseep

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func assertTcpConnections(t *testing.T, got []TcpConnection, want []TcpConnection) bool {
	ret := true
	for i := range got {
		ret = ret && (got[i].equals(want[i]))
	}

	return ret
}

func TestParseTcpConnection(t *testing.T) {

	t.Run("valid TCP entry", func(t *testing.T) {
		tcpString := "  12: E10FA20A:DC1A FEA9FEA9:0050 01 00000000:00000000 00:00000000 00000000     0        0 25370 1 0000000000000000 20 0 0 10 -1"
		got, _ := parseTcpConnection(tcpString)
		want := TcpConnection{
			localAddress:  "10.162.15.225",
			localPort:     56346,
			remoteAddress: "169.254.169.254",
			remotePort:    80,
		}

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("bad TCP entry", func(t *testing.T) {
		tcpString := ""
		_, err := parseTcpConnection(tcpString)
		if err == nil {
			t.Errorf("wanted error but didn't get one")
		}
	})
}

func TestGetCurrentConnections(t *testing.T) {
	got, _ := getCurrentConnections("_testdata/sample_input.t0", time.Time{})
	if len(got) != 3 {
		t.Fatalf("expected 3 connections, got %d", len(got))
	}

	want := []TcpConnection{
		{localAddress: "127.0.0.1", localPort: 9843, remoteAddress: "0.0.0.0", remotePort: 0},
		{localAddress: "127.0.0.53", localPort: 53, remoteAddress: "0.0.0.0", remotePort: 0},
		{localAddress: "0.0.0.0", localPort: 22, remoteAddress: "0.0.0.0", remotePort: 0},
	}

	if !assertTcpConnections(t, got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCurrentConnectionContains(t *testing.T) {
	tcpConnections := []TcpConnection{
		{localAddress: "127.0.0.1", localPort: 9843, remoteAddress: "0.0.0.0", remotePort: 0},
		{localAddress: "127.0.0.53", localPort: 53, remoteAddress: "0.0.0.0", remotePort: 0},
		{localAddress: "0.0.0.0", localPort: 22, remoteAddress: "0.0.0.0", remotePort: 0},
	}
	currentConnections := CurrectConnections{
		Source:      "",
		connections: tcpConnections,
	}

	existingConnection := TcpConnection{localAddress: "127.0.0.53", localPort: 53, remoteAddress: "0.0.0.0", remotePort: 0}

	got := currentConnections.contains(existingConnection)
	want := 1

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func TestCurrentConnectionsUpdate(t *testing.T) {
	currentConnections := CurrectConnections{
		Source: "_testdata/sample_input.t0",
	}
	t0NewConnections, _ := currentConnections.Update()

	if len(t0NewConnections) != 3 {
		t.Errorf("expected 3 new connections but got %d", len(t0NewConnections))
	}

	if len(currentConnections.connections) != 3 {
		t.Errorf("expected current connections to contain 3 connections but has %d", len(currentConnections.connections))
	}

	currentConnections.Source = "_testdata/sample_input.t1"
	t1NewConnections, _ := currentConnections.Update()

	if len(t1NewConnections) != 5 {
		t.Errorf("expected 5 new connections but got %d", len(t1NewConnections))
	}

	if len(currentConnections.connections) != 8 {
		t.Errorf("expected current connections to contain 8 connections but has %d", len(currentConnections.connections))
	}
}

func TestPrintNewConnections(t *testing.T) {
	timestamp, _ := time.Parse("2006-01-02T15:04:05.000Z", "2021-04-28T15:28:15.000Z")
	newConnections := []TcpConnection{
		{localAddress: "10.0.0.5", localPort: 80, remoteAddress: "192.0.2.56", remotePort: 5973, timestamp: timestamp},
		{localAddress: "10.0.0.5", localPort: 80, remoteAddress: "203.0.113.105", remotePort: 31313, timestamp: timestamp},
		{localAddress: "10.0.0.5", localPort: 80, remoteAddress: "203.0.113.94", remotePort: 9208, timestamp: timestamp},
		{localAddress: "10.0.0.5", localPort: 80, remoteAddress: "198.51.100.245", remotePort: 14201, timestamp: timestamp},
	}

	buffer := bytes.Buffer{}
	PrintNewConnections(&buffer, newConnections)

	got := buffer.String()
	want := `2021-04-28 15:28:15: New connection: 192.0.2.56:5973 -> 10.0.0.5:80
2021-04-28 15:28:15: New connection: 203.0.113.105:31313 -> 10.0.0.5:80
2021-04-28 15:28:15: New connection: 203.0.113.94:9208 -> 10.0.0.5:80
2021-04-28 15:28:15: New connection: 198.51.100.245:14201 -> 10.0.0.5:80
`
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestPortScanDetected(t *testing.T) {
	timestamp := time.Date(2021, 8, 8, 14, 44, 0, 0, time.UTC)

	t.Run("3 connections within port scan detection period", func(t *testing.T) {
		tcpConnections := []TcpConnection{
			{localAddress: "10.0.0.5", localPort: 80, remoteAddress: "192.0.2.56", remotePort: 5973, timestamp: timestamp},
			{localAddress: "10.0.0.5", localPort: 81, remoteAddress: "192.0.2.56", remotePort: 5973, timestamp: timestamp.Add(portScanDetectionPeriod / 2)},
			{localAddress: "10.0.0.5", localPort: 82, remoteAddress: "192.0.2.56", remotePort: 5973, timestamp: timestamp.Add(portScanDetectionPeriod)},
		}

		currentConnections := CurrectConnections{
			Source:      "",
			connections: tcpConnections,
		}

		got := currentConnections.checkForPortScans(timestamp.Add(portScanDetectionPeriod))
		want := []PortScan{
			{localAddress: "10.0.0.5", remoteAddress: "192.0.2.56", ports: []int{80, 81, 82}, timestamp: timestamp.Add(portScanDetectionPeriod)},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected to detect port scan")
		}
	})

	t.Run("3 connections outside port scan detection period", func(t *testing.T) {
		tcpConnections := []TcpConnection{
			{localAddress: "10.0.0.5", localPort: 80, remoteAddress: "192.0.2.56", remotePort: 5973, timestamp: timestamp},
			{localAddress: "10.0.0.5", localPort: 81, remoteAddress: "192.0.2.56", remotePort: 5973, timestamp: timestamp.Add(portScanDetectionPeriod / 2)},
			{localAddress: "10.0.0.5", localPort: 82, remoteAddress: "192.0.2.56", remotePort: 5973, timestamp: timestamp.Add(portScanDetectionPeriod + 10*time.Second)},
		}

		currentConnections := CurrectConnections{
			Source:      "",
			connections: tcpConnections,
		}

		got := currentConnections.checkForPortScans(timestamp.Add(portScanDetectionPeriod + 10*time.Second))

		if len(got) != 0 {
			t.Errorf("did not expect to detect a port scan")
		}
	})

	t.Run("don't repeat port scan alerts", func(t *testing.T) {
		// Connection attempts that meet the port scanning criteria
		tcpConnections := []TcpConnection{
			{localAddress: "10.0.0.5", localPort: 80, remoteAddress: "192.0.2.56", remotePort: 5973, timestamp: timestamp},
			{localAddress: "10.0.0.5", localPort: 81, remoteAddress: "192.0.2.56", remotePort: 5973, timestamp: timestamp.Add(portScanDetectionPeriod / 10)},
			{localAddress: "10.0.0.5", localPort: 82, remoteAddress: "192.0.2.56", remotePort: 5973, timestamp: timestamp.Add(portScanDetectionPeriod / 5)},
		}

		currentConnections := CurrectConnections{
			Source:      "",
			connections: tcpConnections,
		}

		got := currentConnections.checkForPortScans(timestamp.Add(portScanDetectionPeriod / 5))

		if len(got) == 0 {
			t.Errorf("expected to detect a port scan")
		}

		// Simulating the next call to Update. We already reported the port scan on the last one,
		// so we don't expect to find one here.
		got = currentConnections.checkForPortScans(timestamp.Add(portScanDetectionPeriod / 5).Add(10 * time.Second))

		if len(got) != 0 {
			t.Errorf("did not expect to detect a port scan")
		}

		currentConnections.connections = append(tcpConnections, TcpConnection{localAddress: "10.0.0.5", localPort: 82, remoteAddress: "192.0.2.56", remotePort: 5973, timestamp: timestamp.Add(portScanDetectionPeriod / 5).Add(10 * time.Second)})
		// Now that there's been another connection attempt, we should report another port scan.
		got = currentConnections.checkForPortScans(timestamp.Add(portScanDetectionPeriod / 5).Add(10 * time.Second))

		if len(got) == 0 {
			t.Errorf("expected to detect a port scan")
		}
	})
}

func TestPrintPortScans(t *testing.T) {
	timestamp, _ := time.Parse("2006-01-02T15:04:05.000Z", "2021-04-28T15:28:05.000Z")
	portScans := []PortScan{
		{localAddress: "10.0.0.5", remoteAddress: "192.0.2.56", ports: []int{80, 81, 82, 83}, timestamp: timestamp},
	}

	buffer := bytes.Buffer{}

	PrintPortScans(&buffer, portScans)

	got := buffer.String()
	want := `2021-04-28 15:28:05: Port scan detected: 192.0.2.56 -> 10.0.0.5 on ports 80,81,82,83`

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
