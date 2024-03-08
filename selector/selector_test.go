package selector

import (
	"path"
	"strings"
	"testing"
)

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

func TestAbiFromSelector(t *testing.T) {
	tests := []struct {
		name     string
		selector string
		path     string
		panics   bool
		want     string // Hex string
	}{
		{
			name:     "erc20 transfer",
			selector: "0xa9059cbb",
			path:     path.Join("testdata", "erc20.abi.json"),
			want:     "transfer(address,uint256)",
		},
		{
			name:     "non existent selector",
			selector: "0xa3063fba",
			path:     path.Join("testdata", "erc20.abi.json"),
			want:     "no method with id",
			panics:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics {
				defer func() {
					if r := recover(); r != nil {
						// Check if the panic value is as expected
						errorString := r.(error).Error()
						wantErrorSubString := tc.want
						if strings.Contains(errorString, wantErrorSubString) == false {
							t.Errorf("Expected panic message to contain: %s, got: %v", wantErrorSubString, errorString)
						}
					} else {
						// The function did not panic as expected
						t.Error("Expected the function to panic, but it did not")
					}
				}()

				abiFromSelector(tc.selector, tc.path)
			} else {
				got := abiFromSelector(tc.selector, tc.path)
				if got != tc.want {
					t.Errorf("abiFromSelector(%s) = %s, want %s", tc.selector, got, tc.want)
				}
			}
		})
	}
}
