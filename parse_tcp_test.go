package main

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func TestParseTcpConnection(t *testing.T) {

	t.Run("valid TCP entry", func(t *testing.T) {
		tcpString := "  12: E10FA20A:DC1A FEA9FEA9:0050 01 00000000:00000000 00:00000000 00000000     0        0 25370 1 0000000000000000 20 0 0 10 -1"
		got, _ := ParseTcpConnection(tcpString)
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
		_, err := ParseTcpConnection(tcpString)
		if err == nil {
			t.Errorf("wanted error but didn't get one")
		}
	})
}

func TestParseListOfConnections(t *testing.T) {
	data, _ := ioutil.ReadFile("_testdata/sample_input.t0")
	connections := string(data)
	got := ParseListOfConnections(connections)
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
	currentConnections := CurrectConnections{
		{localAddress: "127.0.0.1", localPort: 9843, remoteAddress: "0.0.0.0", remotePort: 0},
		{localAddress: "127.0.0.53", localPort: 53, remoteAddress: "0.0.0.0", remotePort: 0},
		{localAddress: "0.0.0.0", localPort: 22, remoteAddress: "0.0.0.0", remotePort: 0},
	}

	existingConnection := TcpConnection{localAddress: "127.0.0.53", localPort: 53, remoteAddress: "0.0.0.0", remotePort: 0}

	got := currentConnections.Contains(existingConnection)
	want := true

	if got != want {
		t.Errorf("could not identify existing connection")
	}
}

func TestCurrentConnectionsUpdate(t *testing.T) {
	t0, _ := ioutil.ReadFile("_testdata/sample_input.t0")
	t0Connections := ParseListOfConnections(string(t0))

	currentConnections := CurrectConnections{}
	t0NewConnections := currentConnections.Update(t0Connections)

	if len(t0NewConnections) != 3 {
		t.Errorf("expected 3 new connections but got %d", len(t0NewConnections))
	}

	if len(currentConnections) != 3 {
		t.Errorf("expected current connections to contain 3 connections but has %d", len(currentConnections))
	}

	t1, _ := ioutil.ReadFile("_testdata/sample_input.t1")
	t1Connections := ParseListOfConnections(string(t1))

	t1NewConnections := currentConnections.Update(t1Connections)

	if len(t1NewConnections) != 5 {
		t.Errorf("expected 5 new connections but got %d", len(t1NewConnections))
	}

	if len(currentConnections) != 8 {
		t.Errorf("expected current connections to contain 8 connections but has %d", len(currentConnections))
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
