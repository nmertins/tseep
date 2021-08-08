package tseep

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
			{hexInput: "01EE", want: MalformedHexInput},
			{hexInput: "", want: MalformedHexInput},
		}

		for _, tc := range tests {
			_, err := ConvertLittleEndianHexToIP(tc.hexInput)
			if err != tc.want {
				t.Errorf("wanted an error but didn't get one")
			}
		}
	})
}

func TestConvertBigEndianHexToPort(t *testing.T) {

	t.Run("valid hex strings", func(t *testing.T) {
		type test struct {
			hexInput string
			want     uint16
		}

		tests := []test{
			{hexInput: "0050", want: 80},
			{hexInput: "FFFF", want: 65535},
			{hexInput: "01BB", want: 443},
		}

		for _, tc := range tests {
			got, _ := ConvertBigEndianHexToPort(tc.hexInput)

			if got != tc.want {
				t.Errorf("got %d want %d", got, tc.want)
			}
		}
	})

	t.Run("invalid hex strings", func(t *testing.T) {
		type test struct {
			hexInput string
			want     error
		}

		tests := []test{
			{hexInput: "", want: MalformedHexInput},
			{hexInput: "00", want: MalformedHexInput},
			{hexInput: "TEST", want: MalformedHexInput},
		}

		for _, tc := range tests {
			_, err := ConvertBigEndianHexToPort(tc.hexInput)

			if err != tc.want {
				t.Errorf("got %q want %q", err, tc.want)
			}
		}
	})
}
