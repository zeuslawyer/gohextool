package encdec

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

/*
 * Decodes `hex` to a string. `hex` must be prexifed with 0x.
 */
func DecodeHexToString(hex string) string {
	bytes := hexutil.MustDecode(hex)
	return string(bytes)
}

func DecodeHexToBigInt(hex string) *big.Int {
	if (len(hex) == 0) || (hex == "0x") || (hex == "0x00") {
		return new(big.Int).SetInt64(0)
	}
	// hexutil requires that integers are encoded using the least amount of digits (no leading zero digits).
	hexWithoutPrefix := hex[2:]
	trimmed := strings.TrimLeft(hexWithoutPrefix, "0")

	bigInt := hexutil.MustDecodeBig("0x" + trimmed)
	return bigInt
}
