package main

import (
	"errors"
	"strings"
)

const (
	// Initial offset in /proc/net/tcp containing indent and connection entry number
	TcpConnectionStartIndex = 6
	// Connection details are given in the format: `<local address>:<local port> <remote address>:<remote port>`
	TcpConnectionEndIndex = TcpConnectionStartIndex + HexCharactersInIPAddress + 1 + HexCharactersInPort + 1 + HexCharactersInIPAddress + 1 + HexCharactersInPort
)

type TcpConnection struct {
	localAddress  string
	localPort     uint16
	remoteAddress string
	remotePort    uint16
}

type CurrectConnections []TcpConnection

func (c CurrectConnections) Contains(other TcpConnection) bool {
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
		if !c.Contains(connection) {
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
		tcp, _ := ParseTcpConnection(line)
		parsedConnections = append(parsedConnections, tcp)
	}

	return parsedConnections
}

// ParseTcpConnection extracts the local and remote IPv4 address and TCP port from single lines of /proc/net/tcp
// passed in as a string s.
//
// ParseTcpConnection only handles single lines, the caller must handle iterating over multiple entries.
// If there is an error parsing the stirng s, an empty TcpConnection struct with the error will be returned.
func ParseTcpConnection(s string) (TcpConnection, error) {
	if len(s) < TcpConnectionEndIndex {
		return TcpConnection{}, errors.New("TCP connection string is malformed")
	}
	connection := s[TcpConnectionStartIndex:TcpConnectionEndIndex]
	split := strings.Split(connection, " ")
	local := strings.Split(split[0], ":")
	remote := strings.Split(split[1], ":")

	localAddress, errLocalAddress := ConvertLittleEndianHexToIP(local[0])
	if errLocalAddress != nil {
		return TcpConnection{}, errLocalAddress
	}
	localPort, errLocalPort := ConvertBigEndianHexToPort(local[1])
	if errLocalPort != nil {
		return TcpConnection{}, errLocalAddress
	}

	remoteAddress, errRemoteAddress := ConvertLittleEndianHexToIP(remote[0])
	if errRemoteAddress != nil {
		return TcpConnection{}, errRemoteAddress
	}
	remotePort, errRemotePort := ConvertBigEndianHexToPort(remote[1])
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
