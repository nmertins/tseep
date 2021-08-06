package main

import (
	"errors"
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
