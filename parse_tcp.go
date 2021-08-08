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
	tcpConnectionLength     = hexCharactersInIPAddress + 1 + hexCharactersInPort + 1 + hexCharactersInIPAddress + 1 + hexCharactersInPort
	portScanDetectionPeriod = time.Duration(60 * time.Second)
	portScanDetectionCount  = 3
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
	portScans   []PortScan
}

type PortScan struct {
	localAddress  string
	remoteAddress string
	ports         []int
	timestamp     time.Time
}

// Looks through connections for matching local and remote address and port.
// If a match is found the index is returned, otherwise -1.
func (c CurrectConnections) contains(other TcpConnection) int {
	ret := -1
	for i, connection := range c.connections {
		if connection.equals(other) {
			ret = i
			break
		}
	}

	return ret
}

// Update checks for TCP connections, compares them with the existing state and
// returns a list of new connections.
func (c *CurrectConnections) Update() (newConnections []TcpConnection, err error) {
	timestamp := time.Now()
	tcpConnections, err := getCurrentConnections(c.Source, timestamp)
	if err != nil {
		return []TcpConnection{}, err
	}

	for _, connection := range tcpConnections {
		idx := c.contains(connection)
		if idx == -1 {
			newConnections = append(newConnections, connection)
		} else {
			c.connections[idx].timestamp = connection.timestamp
		}
	}

	c.connections = append(c.connections, newConnections...)
	c.checkForPortScans(timestamp)

	return newConnections, nil
}

// checkForPortScans interates over the current TCP connections looking for
// remote addresses that have connected to 3 or port local ports in the last
// 60 seconds.
//
// Due to the amount of iteration going on here, it feels like there is probably
// a better way to store/retrieve this data.
func (c CurrectConnections) checkForPortScans(referenceTime time.Time) []PortScan {
	type portsWithTimestamp struct {
		Ports     []int
		Timestamp time.Time
	}

	scanMap := make(map[string]map[string]portsWithTimestamp, 0)

	// For each local address, build a map of remote addresses and the port they connected to.
	for _, connection := range c.connections {
		// Ignore connections outside the detection period
		if referenceTime.Sub(connection.timestamp) > portScanDetectionPeriod {
			continue
		}
		_, ok := scanMap[connection.localAddress]
		if !ok {
			scanMap[connection.localAddress] = map[string]portsWithTimestamp{
				connection.remoteAddress: {Ports: make([]int, 0), Timestamp: connection.timestamp},
			}
		}
		pwt := scanMap[connection.localAddress][connection.remoteAddress]
		pwt.Ports = append(pwt.Ports, int(connection.localPort))
		// update timestamp to match latest connection attempt
		if connection.timestamp.After(pwt.Timestamp) {
			pwt.Timestamp = connection.timestamp
		}

		scanMap[connection.localAddress][connection.remoteAddress] = pwt
	}

	// Look through the map for local/remote address combinations that have mroe than 3 port connections.
	scans := make([]PortScan, 0)
	for localAddress, remoteAddressMap := range scanMap {
		for remoteAddress, pwt := range remoteAddressMap {
			if (len(pwt.Ports) >= portScanDetectionCount) && (pwt.Timestamp.Equal(referenceTime)) {
				scan := PortScan{
					localAddress: localAddress, remoteAddress: remoteAddress, ports: pwt.Ports, timestamp: referenceTime,
				}
				scans = append(scans, scan)
			}
		}
	}

	return scans
}

// getCurrentConnections reads the source filepath for a list of TCP connections.
//
// Each connection is stamped with the current time. If the source file can't
// be read, an empty list is returned witha n error.
func getCurrentConnections(source string, timestamp time.Time) ([]TcpConnection, error) {
	data, err := ioutil.ReadFile(source)
	if err != nil {
		return []TcpConnection{}, err
	}

	connectionsRaw := string(data)
	tcpConnections := parseListOfConnections(connectionsRaw, timestamp)

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

func formatTimestamp(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

// PrintNewConnections outputs the list of TCP connections to writer.
func PrintNewConnections(writer io.Writer, newConnections []TcpConnection) {
	for _, connection := range newConnections {
		timestamp := formatTimestamp(connection.timestamp)
		fmt.Fprintf(writer, "%s: New connection: %s:%d -> %s:%d\n", timestamp, connection.remoteAddress, connection.remotePort, connection.localAddress, connection.localPort)
	}
}

func portString(ports []int) string {
	ret := ""

	for i, port := range ports {
		ret += fmt.Sprint(port)
		if i != len(ports)-1 {
			ret += ","
		}
	}

	return ret
}

func PrintPortScans(writer io.Writer, portScans []PortScan) {
	for _, scan := range portScans {
		timestamp := formatTimestamp(scan.timestamp)
		fmt.Fprintf(writer, "%s: Port scan detected: %s -> %s on ports %s", timestamp, scan.remoteAddress, scan.localAddress, portString(scan.ports))
	}
}
