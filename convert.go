package main

import (
	"encoding/hex"
	"fmt"
)

type ConversionError string

const (
	HexCharactersIn32Bits    = 8
	ErrInvalidCharacterCount = ConversionError("hexadecimal string does not have 8 characters")
)

func (e ConversionError) Error() string {
	return string(e)
}

// ConvertLittleEndianHexToIP returns the IPv4 address represented by the hexadecimal stirng s.
//
// ConvertLittleEndianHexToIP expects the following to be true:
//   - string s to be EXACTLY 8  hexadecimal characters
//   - string s to be formatted in little endian
// If either is not true, an empty string is returned with an error.
func ConvertLittleEndianHexToIP(s string) (string, error) {
	if len(s) != HexCharactersIn32Bits {
		return "", ErrInvalidCharacterCount
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
