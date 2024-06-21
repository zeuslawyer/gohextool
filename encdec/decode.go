package encdec

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

/*
 * Decodes `hex` to a string. `hex` must be prexifed with 0x.
 */
func DecodeHexToString(hex string) string {
	decodedBytes := hexutil.MustDecode(hex)

	return string(decodedBytes) + "\n" // concat newlines so that returned output in terminal pushes terminal prompt "%" to new line.
}

/*
 * Decodes `hex` to a Big Int. `hex` must be prexifed with 0x.
 */
func DecodeHexToBigInt(hex string) *big.Int {
	if hex == "0x" || hex == "" {
		panic(fmt.Sprintf("%q provided as --hex input", hex))
	}

	if (len(hex) == 0) || (hex == "0x00") {
		return new(big.Int).SetInt64(0)
	}

	hexWithoutPrefix := hex[2:]
	bi, _ := new(big.Int).SetString(hexWithoutPrefix, 16)
	return bi
}

/*
 * Given a hex string and a tuple of types, decode the hex string to the corresponding types.
 * `dataTypes` is a comma-separated string of types. Eg: "string, uint, bool, uint".
 * The sequence of types in `dataTypes` must match the sequence of values in the hex string.
 */

func AbiDecode(hexInput string, dataTypes string) []any {
	if strings.HasPrefix(hexInput, "0x") {
		// Strip out the "0x" prefix
		hexInput = hexInput[2:]
	}

	b, err := hex.DecodeString(hexInput)
	if err != nil {
		fmt.Printf("Error decoding hex input: %v", err)
	}

	if len(b) == 0 {
		return []any{} // Empty.
	}

	var args abi.Arguments = dataTypesToAbiArgs(dataTypes)
	values, err := args.Unpack(b)
	if err != nil {
		fmt.Printf("Error Unpacking hex input: %v", err)

		return nil
	}
	for _, val := range values {
		// print type of the value
		fmt.Printf("Decoded value '%v' of type %T\n", val, val)
	}
	return values
}
