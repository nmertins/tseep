package main

import (
	"io/ioutil"
	"reflect"
	"testing"
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
