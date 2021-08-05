package main

import (
	"strings"
)

type TcpConnection struct {
	localAddress  string
	localPort     uint16
	remoteAddress string
	remotePort    uint16
}

const (
	// Initial offset in /proc/net/tcp containing indent and connection entry number
	TcpConnectionStartIndex = 6
	// Connection details are given in the format: `<local address>:<local port> <remote address>:<remote port>`
	TcpConnectionEndIndex = TcpConnectionStartIndex + HexCharactersInIPAddress + 1 + HexCharactersInPort + 1 + HexCharactersInIPAddress + 1 + HexCharactersInPort
)

func ParseTcpConnection(s string) TcpConnection {
	connection := s[TcpConnectionStartIndex:TcpConnectionEndIndex]
	split := strings.Split(connection, " ")
	local := strings.Split(split[0], ":")
	remote := strings.Split(split[1], ":")

	localAddress, _ := ConvertLittleEndianHexToIP(local[0])
	localPort, _ := ConvertBigEndianHexToPort(local[1])

	remoteAddress, _ := ConvertLittleEndianHexToIP(remote[0])
	remotePort, _ := ConvertBigEndianHexToPort(remote[1])

	return TcpConnection{
		localAddress:  localAddress,
		localPort:     localPort,
		remoteAddress: remoteAddress,
		remotePort:    remotePort,
	}
}
