package main

import (
	"bytes"
	"testing"
	"time"
)

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
