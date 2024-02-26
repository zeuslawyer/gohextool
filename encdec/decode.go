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
	bytes := hexutil.MustDecode(hex)
	return string(bytes) + "\n" // concat newlines so that returned output in terminal pushes terminal prompt "%" to new line.
}

// func DecodeHexToBigInt(hex string) *big.Int {
// 	if (len(hex) == 0) || (hex == "0x") || (hex == "0x00") {
// 		return new(big.Int).SetInt64(0)
// 	}
// 	// hexutil requires that integers are encoded using the least amount of digits (no leading zero digits).
// 	hexWithoutPrefix := hex[2:]
// 	trimmed := strings.TrimLeft(hexWithoutPrefix, "0")

// 	bigInt := hexutil.MustDecodeBig("0x" + trimmed)
// 	return bigInt
// }

func DecodeHexToBigInt(hex string) *big.Int {
	if hex == "0x" {
		panic(fmt.Sprintf("%q provided as --hex input", hex))
	}

	if (len(hex) == 0) || (hex == "0x00") {
		return new(big.Int).SetInt64(0)
	}

	hexWithoutPrefix := hex[2:]
	// hexutil requires that integers are encoded using the least amount of digits (no leading zero digits).
	// trimmed := strings.TrimLeft(hexWithoutPrefix, "0")

	num, err := strconv.ParseInt(hexWithoutPrefix, 16, 64)
	if err != nil {
		panic(err)
	}
	return new(big.Int).SetInt64(num)

	// bigInt := hexutil.MustDecodeBig("0x" + trimmed)
	// return bigInt
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
