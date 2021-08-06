package main

import "testing"

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
