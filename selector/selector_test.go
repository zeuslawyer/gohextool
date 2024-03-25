package selector

import (
	"path"
	"strings"
	"testing"
)

const (
	// "Raw" Github Gists
	badJsonUrl  = "https://gist.githubusercontent.com/zeuslawyer/ecec03ff3f50311e510c201de4c076d5/raw/54b14fbfb686e5605e79a4a950031ecaff279d4a/bad-data-erc20.json"
	goodJsonUrl = "https://gist.githubusercontent.com/zeuslawyer/ecec03ff3f50311e510c201de4c076d5/raw/f096531942e922cb3f1d5daa2132f0e476356ced/good-data-erc20.json"
	badUrlExt   = "https://gist.githubusercontent.com/zeuslawyer/ecec03ff3f50311e510c201de4c076d5/raw/54b14fbfb686e5605e79a4a950031ecaff279d4a/bad-data-erc20.NotJson"
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

func TestFuncFromSelector(t *testing.T) {
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
			url:      goodJsonUrl,
			want:     "transfer(address,uint256)",
		},
		{
			name:     "abi from path and url - defaults to path",
			selector: "0xa9059cbb",
			path:     path.Join("testdata", "erc20.abi.json"),
			url:      goodJsonUrl,
			want:     "transfer(address,uint256)",
		},
		{
			name:     "file - invalid JSON",
			selector: "0xa9059cbb",
			path:     path.Join("testdata", "erc20.bad-abi.json"),
			panics:   true,
			want:     "error parsing JSON from file",
		},
		{
			name:     "URL - invalid JSON",
			selector: "0xa9059cbb",
			url:      badJsonUrl,
			panics:   true,
			want:     "error parsing JSON from file",
		},
		{
			name:     "file - not .json extension",
			selector: "0xa9059cbb",
			path:     path.Join("testdata", "erc20.abi.NotJson"),
			panics:   true,
			want:     "invalid file/url extension",
		},
		{
			name:     "URL - not .json file extension",
			selector: "0xa9059cbb",
			url:      badUrlExt,
			panics:   true,
			want:     "invalid file/url extension",
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

				FuncFromSelector(tc.selector, tc.path, tc.url)
			} else {
				got := FuncFromSelector(tc.selector, tc.path, tc.url)
				if got != tc.want {
					t.Errorf("abiFromSelector(%s) = %s, want %s", tc.selector, got, tc.want)
				}
			}
		})
	}
}

func TestEventFromTopicHash(t *testing.T) {
	tests := []struct {
		name      string
		topicHash string
		path      string
		url       string
		panics    bool
		want      string // Hex string
	}{
		// TODO @zeuslawyer do remaining test cases
		{
			name:      "Event Signature from ABI file",
			topicHash: "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", // "0x8c5be1e5",
			path:      path.Join("testdata", "erc20.abi.json"),
			want:      "Approval(address,address,uint256)",
		},
		{
			name:      "Event Signature from ABI file",
			topicHash: "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			path:      path.Join("testdata", "erc20.abi.json"),
			want:      "Transfer(address,address,uint256)",
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

				EventFromTopicHash(tc.topicHash, tc.path, tc.url)
			} else {
				got := EventFromTopicHash(tc.topicHash, tc.path, tc.url)
				if got != tc.want {
					t.Errorf("abiFromSelector(%s) = %s, want %s", tc.topicHash, got, tc.want)
				}
			}
		})
	}
}
