package tseep

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	// Initial offset in /proc/net/tcp containing indent and connection entry number
	tcpConnectionStartIndex = 6
	// Connection details are given in the format: `<local address>:<local port> <remote address>:<remote port>`
	tcpConnectionEndIndex = tcpConnectionStartIndex + hexCharactersInIPAddress + 1 + hexCharactersInPort + 1 + hexCharactersInIPAddress + 1 + hexCharactersInPort
)

type TcpConnection struct {
	localAddress  string
	localPort     uint16
	remoteAddress string
	remotePort    uint16
}

type CurrectConnections []TcpConnection

func (c CurrectConnections) contains(other TcpConnection) bool {
	ret := false
	for _, connection := range c {
		if connection == other {
			ret = true
		}
	}

	return ret
}

func (c *CurrectConnections) Update(tcpConnections []TcpConnection) (newConnections []TcpConnection) {
	for _, connection := range tcpConnections {
		if !c.contains(connection) {
			newConnections = append(newConnections, connection)
		}
	}

	*c = tcpConnections

	return newConnections
}

func ParseListOfConnections(connections string) []TcpConnection {
	lines := strings.Split(connections, "\n")
	var parsedConnections []TcpConnection

	for _, line := range lines[1:] {
		tcp, _ := parseTcpConnection(line)
		parsedConnections = append(parsedConnections, tcp)
	}

	return parsedConnections
}

// parseTcpConnection extracts the local and remote IPv4 address and TCP port from single lines of /proc/net/tcp
// passed in as a string s.
//
// parseTcpConnection only handles single lines, the caller must handle iterating over multiple entries.
// If there is an error parsing the stirng s, an empty TcpConnection struct with the error will be returned.
func parseTcpConnection(s string) (TcpConnection, error) {
	if len(s) < tcpConnectionEndIndex {
		return TcpConnection{}, errors.New("TCP connection string is malformed")
	}
	connection := s[tcpConnectionStartIndex:tcpConnectionEndIndex]
	split := strings.Split(connection, " ")
	local := strings.Split(split[0], ":")
	remote := strings.Split(split[1], ":")

	localAddress, errLocalAddress := convertLittleEndianHexToIP(local[0])
	if errLocalAddress != nil {
		return TcpConnection{}, errLocalAddress
	}
	localPort, errLocalPort := convertBigEndianHexToPort(local[1])
	if errLocalPort != nil {
		return TcpConnection{}, errLocalAddress
	}

	remoteAddress, errRemoteAddress := convertLittleEndianHexToIP(remote[0])
	if errRemoteAddress != nil {
		return TcpConnection{}, errRemoteAddress
	}
	remotePort, errRemotePort := convertBigEndianHexToPort(remote[1])
	if errRemotePort != nil {
		return TcpConnection{}, errRemotePort
	}

	return TcpConnection{
		localAddress:  localAddress,
		localPort:     localPort,
		remoteAddress: remoteAddress,
		remotePort:    remotePort,
	}, nil
}

func PrintNewConnections(writer io.Writer, t time.Time, newConnections []TcpConnection) {
	timestamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	for _, connection := range newConnections {
		fmt.Fprintf(writer, "%s: New connection: %s:%d -> %s:%d\n", timestamp, connection.remoteAddress, connection.remotePort, connection.localAddress, connection.localPort)
	}
}
