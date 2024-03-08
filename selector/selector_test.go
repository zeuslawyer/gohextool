package selector

import (
	"path"
	"strings"
	"testing"
)

func TestSelectorFromSig(t *testing.T) {
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
			panics:      true,
			want:        "not a valid function signature",
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

				SelectorFromSig(tc.functionSig)
			} else {
				got := SelectorFromSig(tc.functionSig)
				if got != tc.want {
					t.Errorf("FunctionSelector(%s) = %s, want %s", tc.functionSig, got, tc.want)
				}
			}
		})
	}
}

func TestSigFromSelector(t *testing.T) {
	tests := []struct {
		name     string
		selector string
		path     string
		url      string
		panics   bool
		want     string // Hex string
	}{
		{
			name:     "abi from file",
			selector: "0xa9059cbb",
			path:     path.Join("testdata", "erc20.abi.json"),
			want:     "transfer(address,uint256)",
		},
		{
			name:     "abi from url",
			selector: "0xa9059cbb",
			url:      "https://gist.githubusercontent.com/veox/8800debbf56e24718f9f483e1e40c35c/raw/f853187315486225002ba56e5283c1dba0556e6f/erc20.abi.json",
			want:     "transfer(address,uint256)",
		},
		{
			name:     "abi from path and url - defaults to path",
			selector: "0xa9059cbb",
			path:     path.Join("testdata", "erc20.abi.json"),
			url:      "https://gist.githubusercontent.com/veox/8800debbf56e24718f9f483e1e40c35c/raw/f853187315486225002ba56e5283c1dba0556e6f/erc20.abi.json",
			want:     "transfer(address,uint256)",
		},
		{
			name:     "non existent selector",
			selector: "0xa3063fba",
			path:     path.Join("testdata", "erc20.abi.json"),
			panics:   true,
			want:     "no method with id",
		},
		{
			name:     "invalid abi path",
			selector: "0xa3063fba",
			path:     path.Join("invalid-testdata-path", "erc20.abi.json"),
			panics:   true,
			want:     "no such file or directory",
		},
		{
			name:   "empty path, empty url",
			panics: true,
			want:   "abiPath and url cannot both be empty",
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

				SigFromSelector(tc.selector, tc.path, tc.url)
			} else {
				got := SigFromSelector(tc.selector, tc.path, tc.url)
				if got != tc.want {
					t.Errorf("abiFromSelector(%s) = %s, want %s", tc.selector, got, tc.want)
				}
			}
		})
	}
}
