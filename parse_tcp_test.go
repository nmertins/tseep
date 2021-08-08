package tseep

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

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
	got, _ := getCurrentConnections("_testdata/sample_input.t0")
	if len(got) != 3 {
		t.Fatalf("expected 3 connections, got %d", len(got))
	}

	want := []TcpConnection{
		{localAddress: "127.0.0.1", localPort: 9843, remoteAddress: "0.0.0.0", remotePort: 0},
		{localAddress: "127.0.0.53", localPort: 53, remoteAddress: "0.0.0.0", remotePort: 0},
		{localAddress: "0.0.0.0", localPort: 22, remoteAddress: "0.0.0.0", remotePort: 0},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestCurrentConnectionContains(t *testing.T) {
	tcpConnections := []TcpConnection{
		{localAddress: "127.0.0.1", localPort: 9843, remoteAddress: "0.0.0.0", remotePort: 0},
		{localAddress: "127.0.0.53", localPort: 53, remoteAddress: "0.0.0.0", remotePort: 0},
		{localAddress: "0.0.0.0", localPort: 22, remoteAddress: "0.0.0.0", remotePort: 0},
	}
	currentConnections := CurrectConnections{
		source:      "",
		connections: tcpConnections,
	}

	existingConnection := TcpConnection{localAddress: "127.0.0.53", localPort: 53, remoteAddress: "0.0.0.0", remotePort: 0}

	got := currentConnections.contains(existingConnection)
	want := true

	if got != want {
		t.Errorf("could not identify existing connection")
	}
}

func TestCurrentConnectionsUpdate(t *testing.T) {
	currentConnections := CurrectConnections{
		source:      "_testdata/sample_input.t0",
		connections: []TcpConnection{},
	}
	t0NewConnections, _ := currentConnections.Update()

	if len(t0NewConnections) != 3 {
		t.Errorf("expected 3 new connections but got %d", len(t0NewConnections))
	}

	if len(currentConnections.connections) != 3 {
		t.Errorf("expected current connections to contain 3 connections but has %d", len(currentConnections.connections))
	}

	currentConnections.source = "_testdata/sample_input.t1"
	t1NewConnections, _ := currentConnections.Update()

	if len(t1NewConnections) != 5 {
		t.Errorf("expected 5 new connections but got %d", len(t1NewConnections))
	}

	if len(currentConnections.connections) != 8 {
		t.Errorf("expected current connections to contain 8 connections but has %d", len(currentConnections.connections))
	}
}

func TestPrintNewConnections(t *testing.T) {
	newConnections := []TcpConnection{
		{localAddress: "10.0.0.5", localPort: 80, remoteAddress: "192.0.2.56", remotePort: 5973},
		{localAddress: "10.0.0.5", localPort: 80, remoteAddress: "203.0.113.105", remotePort: 31313},
		{localAddress: "10.0.0.5", localPort: 80, remoteAddress: "203.0.113.94", remotePort: 9208},
		{localAddress: "10.0.0.5", localPort: 80, remoteAddress: "198.51.100.245", remotePort: 14201},
	}

	buffer := bytes.Buffer{}
	timestamp, _ := time.Parse("2006-01-02T15:04:05.000Z", "2021-04-28T15:28:15.000Z")
	PrintNewConnections(&buffer, timestamp, newConnections)
	want := `2021-04-28 15:28:15: New connection: 192.0.2.56:5973 -> 10.0.0.5:80
2021-04-28 15:28:15: New connection: 203.0.113.105:31313 -> 10.0.0.5:80
2021-04-28 15:28:15: New connection: 203.0.113.94:9208 -> 10.0.0.5:80
2021-04-28 15:28:15: New connection: 198.51.100.245:14201 -> 10.0.0.5:80
`
	if buffer.String() != want {
		t.Errorf("got %q want %q", buffer.String(), want)
	}
}
