package encdec

import (
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

func TestAbiDecode(t *testing.T) {
	tests := []struct {
		name      string
		inputHex  string
		dataTypes string
		want      []any
	}{
		{
			name:      "HappyPath#1-0x prefix",
			inputHex:  "0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000a676f2d686578746f6f6c00000000000000000000000000000000000000000000",
			dataTypes: "string",
			want:      []any{"go-hextool"}, // see https://adibas03.github.io/online-ethereum-abi-encoder-decoder/#/decode to get decoded scalar values
		},
		{
			name:      "HappyPath#2-multiple scalar types",
			inputHex:  "0x000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000166f4b60000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000d68617070792074657374696e6700000000000000000000000000000000000000",
			dataTypes: "uint, uint256, string",
			want:      []any{big.NewInt(10), big.NewInt(23524534), "happy testing"},
		},
		{
			name:      "HappyPath#3-Prefix Gets Added",
			inputHex:  "000000000000000000000000000000000000000000000000000000000000002b",
			dataTypes: "uint16",
			want:      []any{uint16(43)},
		},
		{
			name:      "HappyPath#4-multiple scalars including address",
			inputHex:  "00000000000000000000000000000000000000000000000000000000000003e90000000000000000000000000000000000000000000000000000000000000060000000000000000000000000208aa722aca42399eac5192ee778e4d42f4e5de300000000000000000000000000000000000000000000000000000000000000057a7562696e000000000000000000000000000000000000000000000000000000",
			dataTypes: "uint16, string, address",
			want:      []any{uint16(1001), "zubin", common.HexToAddress("0x208aa722aca42399eac5192ee778e4d42f4e5de3")},
		},
		{
			name:      "HappyPath#5-empty input",
			inputHex:  "0x",
			dataTypes: "address",
			want:      []any{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := AbiDecode(tc.inputHex, tc.dataTypes)
			if len(got) != len(tc.want) {
				t.Errorf("%s failing because returned slice have unequal length, %d & %d", tc.name, len(got), len(tc.want))
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("AbiDecode() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestAbiEncode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		dataTypes string
		want      string
		panics    bool
	}{
		{
			name:      "HappyPath#1-emtpy string",
			input:     "",
			dataTypes: "uint256",
			want:      "0x",
		},
		{
			name:      "HappyPath#2-multiple scalar types",
			input:     "8,1234567,hextool",
			dataTypes: "uint16, uint, string",
			want:      "0x0000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000012d68700000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000007686578746f6f6c00000000000000000000000000000000000000000000000000",
		},
		{
			name:      "HappyPath#3-with address",
			input:     "8,0x208AA722Aca42399eaC5192EE778e4D42f4E5De3,hextool is rad",
			dataTypes: "uint32,address,string",
			want:      "0x0000000000000000000000000000000000000000000000000000000000000008000000000000000000000000208aa722aca42399eac5192ee778e4d42f4e5de30000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000e686578746f6f6c20697320726164000000000000000000000000000000000000",
		},
		{
			name:      "HappyPath#3-with address",
			input:     "8,0x208AA722Aca42399eaC5192EE778e4D42f4E5De3,hextool is rad",
			dataTypes: "uint32,address,string",
			want:      "0x0000000000000000000000000000000000000000000000000000000000000008000000000000000000000000208aa722aca42399eac5192ee778e4d42f4e5de30000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000e686578746f6f6c20697320726164000000000000000000000000000000000000",
		},
		{
			name:      "HappyPath#3-panics with mismatched input length",
			panics:    true,
			input:     "hextool is rad",
			dataTypes: "uint32,address,string",
			want:      "Number of input values does not match number of types",
		},
		{
			name:      "HappyPath#4-panics with unknown type",
			panics:    true,
			input:     "8,0x208AA722Aca42399eaC5192EE778e4D42f4E5De3,hextool is rad",
			dataTypes: "uint999,address,string",
			want:      "Unsupported type",
		},
		// {
		// 	name:      "HappyPath#4-multiple scalars including address",
		// 	inputHex:  "00000000000000000000000000000000000000000000000000000000000003e90000000000000000000000000000000000000000000000000000000000000060000000000000000000000000208aa722aca42399eac5192ee778e4d42f4e5de300000000000000000000000000000000000000000000000000000000000000057a7562696e000000000000000000000000000000000000000000000000000000",
		// 	dataTypes: "uint16, string, address",
		// 	want:      []any{uint16(1001), "zubin", common.HexToAddress("0x208aa722aca42399eac5192ee778e4d42f4e5de3")},
		// },
		// {
		// 	name:      "HappyPath#5-empty input",
		// 	inputHex:  "0x",
		// 	dataTypes: "address",
		// 	want:      []any{},
		// },
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics {
				defer func() {
					if r := recover(); r != nil {
						// Check if the panic value is as expected
						errorString := r.(string)
						wantErrorSubString := tc.want

						if strings.Contains(errorString, wantErrorSubString) == false {
							t.Errorf("Expected panic message to contain: %s, got: %v", wantErrorSubString, errorString)
						}
					} else {
						// The function did not panic as expected
						t.Error("Expected the function to panic, but it did not")
					}
				}()

				AbiEncode(tc.input, tc.dataTypes)
			} else {
				got := AbiEncode(tc.input, tc.dataTypes)
				// if got not equal to want fail test
				if got != tc.want {
					t.Errorf("AbiEncode() = %v, want %v", got, tc.want)
				}
			}

		})
	}
}
