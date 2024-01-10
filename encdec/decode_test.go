package encdec

import (
	"math/big"
	"testing"
)

func TestDecodeHexToBigInt(t *testing.T) {
	tests := []struct {
		name     string
		inputHex string
		want     *big.Int
	}{
		{
			name:     "HappyPath",
			inputHex: "0x0000000000000000000000000000000000000000000000000000000000467390",
			want:     new(big.Int).SetInt64(4617104),
		},
		{
			name:     "Zero",
			inputHex: "0x00",
			want:     new(big.Int).SetInt64(0),
		},
		{
			name:     "Empty",
			inputHex: "0x",
			want:     new(big.Int).SetInt64(0),
		},
		{
			name:     "NegativeNum",
			inputHex: "0xF843",
			want:     new(big.Int).SetInt64(-1981),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := DecodeHexToBigInt(tc.inputHex)
			if got.Cmp(tc.want) != 0 {
				t.Errorf("Failing Test Name: %q - DecodeHexToBigInt() = %q, want %q", tc.name, got.String(), tc.want.String())
			}
		})
	}
}

func TestDecodeHexToString(t *testing.T) {
	tests := []struct {
		name     string
		inputHex string
		want     string
	}{
		{
			name:     "HappyPath",
			inputHex: "0x476f20466f727468202620436f6e717565722c20486f6d696521",
			want:     "Go Forth & Conquer, Homie!",
		},
		{
			name:     "NumberAsString",
			inputHex: "0x3432",
			want:     "42",
		},
		{
			name:     "EmptyBytes",
			inputHex: "0x",
			want:     "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := DecodeHexToString(tc.inputHex)
			if got != tc.want {
				t.Errorf("DecodeHexToString() = %q, want %q", got, tc.want)
			}
		})
	}
}
