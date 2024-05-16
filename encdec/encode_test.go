package encdec

import (
	"strings"
	"testing"
)

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
