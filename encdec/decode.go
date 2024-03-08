package encdec

import (
	"fmt"
	"math/big"
	"strconv"

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
 * Decodes `hex` to a Big Int64. `hex` must be prexifed with 0x.
 */
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
