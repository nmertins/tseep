package tseep

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type ConversionError string

const (
	HexCharactersInIPAddress = 8
	HexCharactersInPort      = 4
	MalformedHexInput        = ConversionError("hexadecimal string is malformed")
)

func (e ConversionError) Error() string {
	return string(e)
}

// ConvertLittleEndianHexToIP returns the IPv4 address represented by the hexadecimal stirng s.
//
// ConvertLittleEndianHexToIP expects the following to be true:
//   - string s to be EXACTLY 8 hexadecimal characters
//   - string s to be formatted in little endian
// If s is not of length 8 or contains non-hexadecimal characters, an empty string is returned with an error.
func ConvertLittleEndianHexToIP(s string) (string, error) {
	if len(s) != HexCharactersInIPAddress {
		return "", MalformedHexInput
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
func ConvertBigEndianHexToPort(s string) (uint16, error) {
	if len(s) != HexCharactersInPort {
		return 0, MalformedHexInput
	}

	portBytes, err := hex.DecodeString(s)
	if err != nil {
		return 0, MalformedHexInput
	}

	port := binary.BigEndian.Uint16(portBytes)

	return port, nil
}
