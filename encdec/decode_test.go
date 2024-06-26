package encdec

import (
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestDecodeHexToBigInt(t *testing.T) {
	wantedBigInt, ok := new(big.Int).SetString("139000000000000000000", 10)
	if !ok {
		t.Fatalf("Cannot set big.Int from string %q", "139000000000000000000")
	}

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
			inputHex: "0x-7bd", // "0x" + strconv.FormatInt(-1981, 16), // "0x-7bd"
			want:     new(big.Int).SetInt64(-1981),
		},
		{
			name:     "BigInt",
			inputHex: "0x078903338be34c0000",
			want:     wantedBigInt,
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
	bigNegativeNum, ok := new(big.Int).SetString("-139000000000000000000", 10)
	if !ok {
		t.Fatalf("Cannot set big.Int from string %q", "-139000000000000000000")

	}
	tests := []struct {
		name      string
		inputHex  string
		dataTypes string
		want      []any
	}{
		// see https://adibas03.github.io/online-ethereum-abi-encoder-decoder/#/decode to get decoded scalar values
		{
			name:      "HappyPath#1-single negative int",
			inputHex:  "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff843",
			dataTypes: "int",
			want:      []any{big.NewInt(-1981)},
		},
		{
			name:      "HappyPath#1-0x prefix",
			inputHex:  "0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000a676f2d686578746f6f6c00000000000000000000000000000000000000000000",
			dataTypes: "string",
			want:      []any{"go-hextool"}, // see https://adibas03.github.io/online-ethereum-abi-encoder-decoder/#/decode to get decoded scalar values
		},
		{
			name:      "HappyPath#2-multiple scalar types",
			inputHex:  "0x000000000000000000000000000000000000000000000000000000000000000afffffffffffffffffffffffffffffffffffffffffffffff876fccc741cb400000000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000d68617070792074657374696e6700000000000000000000000000000000000000",
			dataTypes: "uint, int256, string",
			want:      []any{big.NewInt(10), bigNegativeNum, "happy testing"},
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
