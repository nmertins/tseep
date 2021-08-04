package main

import (
	"testing"
)

func TestConvertLittleEndianHexToIP(t *testing.T) {

	t.Run("valid hex string", func(t *testing.T) {
		hexInput := "0100007F"
		got, _ := ConvertLittleEndianHexToIP(hexInput)
		want := "127.0.0.1"

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("invalid hex string", func(t *testing.T) {
		hexInput := "01EE"
		_, err := ConvertLittleEndianHexToIP(hexInput)

		if err != ErrInvalidCharacterCount {
			t.Errorf("wanted an error but didn't get one")
		}
	})
}
