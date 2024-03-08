package selector

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Calculates the function selector given function signature `funcSig`.
// The function signature should be in the form of `functionName(type1,type2,...)`.
// Eg: "transfer(address,uint256)"
func SelectorFromSig(funcSig string) string {
	// TODO @zeuslawyer change to use the abi packages Method type
	validateInput := func(sig string) error {
		re := regexp.MustCompile(`^(\w+)`) // match the first word in a given string
		matches := re.FindStringSubmatch(sig)

		if len(matches) < 2 {
			return fmt.Errorf("unable to extract function name from signature: %s", sig)
		}

		// validate signature format
		signatureRegex := regexp.MustCompile(`^\w+\([^\)]*\)$`)
		if !signatureRegex.MatchString(sig) {
			return fmt.Errorf("\n%q is not a valid function signature", sig)
		}

		return nil
	}

	if err := validateInput(funcSig); err != nil {
		panic(err)
	}

	funcSigHash := crypto.Keccak256Hash([]byte(funcSig))
	selector := funcSigHash.String()[:10] // first 4 bytes ==8 characters, plus "0x"
	return selector
}

// Given a function selector, returns the function signature from provided ABI file and path
// or from a URL.  If both are provided it will default to using the file path.
func SigFromSelector(selector string, abiPath string, abiUrl string) string {
	if abiPath == "" && abiUrl == "" {
		panic(fmt.Errorf("abiPath and url cannot both be empty"))
	}

	selectorBytes := hexutil.MustDecode(selector)

	var abiJson []byte
	var abiSource string
	if abiPath != "" {
		abiSource = abiPath
		file, err := os.ReadFile(abiSource)
		if err != nil {
			fmt.Println("Error reading file")
			panic(err)
		}

		abiJson = file
	} else {
		abiSource = abiUrl
		resp, err := http.Get(abiUrl)
		if err != nil {
			fmt.Printf("Error fetching ABI file from url %s", abiUrl)
			panic(err)
		}
		defer resp.Body.Close()

		abiJson, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading ABI from http response")
			panic(err)
		}
	}

	parsedAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		fmt.Println("Error parsing ABI from : ", abiSource)
		panic(err)
	}

	method, err := parsedAbi.MethodById(selectorBytes)
	if err != nil {
		fmt.Println("Error looking up method by its ID")
		panic(err)
	}
	if method == nil {
		return fmt.Sprintf("Method not found in file at %s", abiPath)
	}

	return method.Sig
}
