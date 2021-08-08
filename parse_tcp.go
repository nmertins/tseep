package tseep

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

const (
	// Connection details are given in the format: `<local address>:<local port> <remote address>:<remote port>`
	tcpConnectionLength = hexCharactersInIPAddress + 1 + hexCharactersInPort + 1 + hexCharactersInIPAddress + 1 + hexCharactersInPort
)

type TcpConnection struct {
	localAddress  string
	localPort     uint16
	remoteAddress string
	remotePort    uint16
	timestamp     time.Time
}

func (t TcpConnection) equals(other TcpConnection) bool {
	ret := true
	ret = ret && t.localAddress == other.localAddress
	ret = ret && (t.localPort == other.localPort)
	ret = ret && (t.remoteAddress == other.remoteAddress)
	ret = ret && (t.remotePort == other.remotePort)

	return ret
}

type CurrectConnections struct {
	Source      string
	connections []TcpConnection
}

func (c CurrectConnections) contains(other TcpConnection) bool {
	ret := false
	for _, connection := range c.connections {
		ret = ret || connection.equals(other)
	}

	return ret
}

// Update checks for TCP connections, compares them with the existing state and
// returns a list of new connections.
func (c *CurrectConnections) Update() (newConnections []TcpConnection, err error) {
	tcpConnections, err := getCurrentConnections(c.Source)
	if err != nil {
		return []TcpConnection{}, err
	}

	for _, connection := range tcpConnections {
		if !c.contains(connection) {
			newConnections = append(newConnections, connection)
		}
	}

	c.connections = tcpConnections

	return newConnections, nil
}

// getCurrentConnections reads the source filepath for a list of TCP connections.
//
// Each connection is stamped with the current time. If the source file can't
// be read, an empty list is returned witha n error.
func getCurrentConnections(source string) ([]TcpConnection, error) {
	data, err := ioutil.ReadFile(source)
	if err != nil {
		return []TcpConnection{}, err
	}

	connectionsRaw := string(data)
	tcpConnections := parseListOfConnections(connectionsRaw, time.Now())

	return tcpConnections, nil
}

// parseListOfConnections translates a multiline string of connections into
// TcpConnection objects stamped with the given timestamp.
func parseListOfConnections(connections string, timestamp time.Time) []TcpConnection {
	lines := strings.Split(connections, "\n")
	var parsedConnections []TcpConnection

	for _, line := range lines[1:] {
		tcp, _ := parseTcpConnection(line)
		tcp.timestamp = timestamp
		parsedConnections = append(parsedConnections, tcp)
	}

	return parsedConnections
}

// parseTcpConnection extracts the local and remote IPv4 address and TCP port from
// single lines of /proc/net/tcp passed in as a string s.
//
// parseTcpConnection only handles single lines, the caller must handle iterating over
// multiple entries. If there is an error parsing the stirng s, an empty TcpConnection
// struct with the error will be returned.
func parseTcpConnection(s string) (TcpConnection, error) {
	idx := strings.Index(s, ":")
	if idx < 0 {
		return TcpConnection{}, errors.New("TCP connection string is malformed")
	}
	startIdx := idx + 2
	endIdx := startIdx + tcpConnectionLength
	connection := s[startIdx:endIdx]
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

// PrintNewConnections outputs the list of TCP connections to writer.
func PrintNewConnections(writer io.Writer, t time.Time, newConnections []TcpConnection) {
	timestamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	for _, connection := range newConnections {
		fmt.Fprintf(writer, "%s: New connection: %s:%d -> %s:%d\n", timestamp, connection.remoteAddress, connection.remotePort, connection.localAddress, connection.localPort)
	}
}
