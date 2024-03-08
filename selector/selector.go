package selector

import (
	"fmt"
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
func FunctionSelector(funcSig string) string {

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

func abiFromSelector(selector string, path string) string {
	selectorBytes := hexutil.MustDecode(selector)

	abiJson, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file")
		panic(err)
	}

	parsedAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		fmt.Println("Error parsing ABI")
		panic(err)
	}

	method, err := parsedAbi.MethodById(selectorBytes)
	if err != nil {
		fmt.Println("Error looking up method by its ID")
		panic(err)
	}
	if method == nil {
		return fmt.Sprintf("Method not found in file at %s", path)
	}

	return method.Sig
}
