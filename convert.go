package tseep

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type conversionError string

const (
	hexCharactersInIPAddress = 8
	hexCharactersInPort      = 4
	malformedHexInput        = conversionError("hexadecimal string is malformed")
)

func (e conversionError) Error() string {
	return string(e)
}

// ConvertLittleEndianHexToIP returns the IPv4 address represented by the hexadecimal stirng s.
//
// ConvertLittleEndianHexToIP expects the following to be true:
//   - string s to be EXACTLY 8 hexadecimal characters
//   - string s to be formatted in little endian
// If s is not of length 8 or contains non-hexadecimal characters, an empty string is returned with an error.
func convertLittleEndianHexToIP(s string) (string, error) {
	if len(s) != hexCharactersInIPAddress {
		return "", malformedHexInput
	}

	octets, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}

	var ip string

	for i := len(octets) - 1; i >= 0; i-- {
		ip += fmt.Sprint(octets[i])
		if i > 0 {
			ip += "."
		}
	}

	return ip, nil
}

// ConvertBigEndianHexToPort returns the TCP port represented by the hexadecimal stirng s.
//
// ConvertBigEndianHexToPort expects the following to be true:
//   - string s to be EXACTLY 4 hexadecimal characters
//   - string s to be formatted in big endian
// If s is not of length 4 or contains non-hexadecimal characters, 0 will be returned with an error.
func convertBigEndianHexToPort(s string) (uint16, error) {
	if len(s) != hexCharactersInPort {
		return 0, malformedHexInput
	}

	portBytes, err := hex.DecodeString(s)
	if err != nil {
		return 0, malformedHexInput
	}

	port := binary.BigEndian.Uint16(portBytes)

	return port, nil
}
