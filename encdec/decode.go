package encdec

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
 * Decodes `hex` to a string. `hex` must be prexifed with 0x.
 */
func DecodeHexToString(hex string) string {
	decodedBytes := hexutil.MustDecode(hex)

	return string(decodedBytes) + "\n" // concat newlines so that returned output in terminal pushes terminal prompt "%" to new line.
}

func DecodeHexToBigInt(hex string) *big.Int {
	if hex == "0x" || len(hex) == 0 {
		panic(fmt.Sprintf("%q provided as --hex input", hex))
	}

	if (len(hex) == 0) || (hex == "0x00") {
		return new(big.Int).SetInt64(0)
	}

	hexWithoutPrefix := hex[2:]
	num, err := strconv.ParseInt(hexWithoutPrefix, 16, 64)
	if err != nil {
		panic(err)
	}
	return new(big.Int).SetInt64(num)
}

// TODO: @zeuslawyer resume here.
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
	selector := funcSigHash.String()[:10] //4 bytes ==8 characters, plus "0x"
	return selector
}
