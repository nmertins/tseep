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
		type test struct {
			hexInput string
			want     error
		}

		tests := []test{
			{hexInput: "01EE", want: ErrInvalidCharacterCount},
			{hexInput: "", want: ErrInvalidCharacterCount},
		}

		for _, tc := range tests {
			_, err := ConvertLittleEndianHexToIP(tc.hexInput)
			if err != tc.want {
				t.Errorf("wanted an error but didn't get one")
			}
		}
	})
}
