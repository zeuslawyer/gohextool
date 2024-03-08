package encdec

import (
	"math/big"
	"strconv"
	"strings"
	"testing"
)

func TestDecodeHexToBigInt(t *testing.T) {
	tests := []struct {
		name     string
		inputHex string
		want     *big.Int
		panics   bool
	}{
		{
			name:     "HappyPath",
			inputHex: "0x0000000000000000000000000000000000000000000000000000000000467390",
			want:     new(big.Int).SetInt64(4617104),
		},
		{
			name:     "0x panics",
			inputHex: "0x",
			panics:   true,
		},
		{
			name:     "Empty String",
			inputHex: "0x",
			panics:   true,
		},
		{
			name:     "Zero",
			inputHex: "0x00",
			want:     new(big.Int).SetInt64(0),
		},
		{
			name:     "NegativeNum",
			inputHex: "0x" + strconv.FormatInt(-1981, 16), // "0x-7bd"
			want:     new(big.Int).SetInt64(-1981),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics {
				defer func() {
					if r := recover(); r != nil {
						// Check if the panic value is as expected
						errorString := r.(string)
						wantErrorSubString := "\"0x\" provided as --hex input"

						if strings.Contains(errorString, wantErrorSubString) == false {
							t.Errorf("Expected panic message to contain: %s, got: %v", wantErrorSubString, errorString)
						}
					} else {
						// The function did not panic as expected
						t.Error("Expected the function to panic, but it did not")
					}
				}()

				DecodeHexToBigInt(tc.inputHex)
			} else {
				got := DecodeHexToBigInt(tc.inputHex)
				if got.Cmp(tc.want) != 0 {
					t.Errorf("Failing Test Name: %q - DecodeHexToBigInt() = %q, want %q", tc.name, got.String(), tc.want.String())
				}
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
			want:     "Go Forth & Conquer, Homie!" + "\n",
		},
		{
			name:     "NumberAsString",
			inputHex: "0x3432",
			want:     "42" + "\n",
		},
		{
			name:     "EmptyBytes",
			inputHex: "0x",
			want:     "" + "\n",
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

func TestFunctionSelector(t *testing.T) {
	tests := []struct {
		name        string
		functionSig string
		panics      bool
		want        string // Hex string
	}{
		{
			name:        "greet",
			functionSig: "greet(string)",
			want:        "0xead710c4", // https://www.evm-function-selector.click/
		},
		{
			name:        "basic transfer",
			functionSig: "transfer(address,uint256)",
			want:        "0xa9059cbb", // https://www.evm-function-selector.click/
		},
		{
			name:        "bad function",
			functionSig: "gibberish",
			want:        "0xa9059cbb",
			panics:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics {
				defer func() {
					if r := recover(); r != nil {
						// Check if the panic value is as expected
						errorString := r.(error).Error()
						wantErrorSubString := "not a valid function signature"

						if strings.Contains(errorString, wantErrorSubString) == false {
							t.Errorf("Expected panic message to contain: %s, got: %v", wantErrorSubString, errorString)
						}
					} else {
						// The function did not panic as expected
						t.Error("Expected the function to panic, but it did not")
					}
				}()

				FunctionSelector(tc.functionSig)
			} else {
				got := FunctionSelector(tc.functionSig)
				if got != tc.want {
					t.Errorf("FunctionSelector(%s) = %s, want %s", tc.functionSig, got, tc.want)
				}
			}
		})
	}
}
