package encdec

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

/*
  - Given a tuple of data values and a tuple of their corresponding types,
  - ABI-encode the data values according to their provided types
  - `input` is a comma-separated string of values. Eg: "hello, 123, true, 456".
    `dataTypes` is a comma-separated string of types. Eg: "string, uint, bool, uint".
  - The sequence of types in `dataTypes` must match the sequence of values in `input`.
*/
func AbiEncode(inputValues string, dataTypes string) (res string) {
	if len(inputValues) == 0 {
		res = "0x"
		return res
	}

	// Split  each input into a slice of string values
	inputAsSlice := strings.Split(inputValues, ",")
	typesAsSlice := strings.Split(dataTypes, ",")
	if len(inputAsSlice) != len(typesAsSlice) {
		panic(fmt.Sprintf("Number of input values does not match number of types - %d inputs to  %d types", len(inputAsSlice), len(typesAsSlice)))

	}

	// convert strings to the appropriate types
	// TODO zubin resume here. Add other types. Complete the tests.
	typedInputValuesSlice := make([]any, len(typesAsSlice))
	for idx, ty := range typesAsSlice {
		_ty := strings.TrimSpace(ty)
		inp := inputAsSlice[idx]

		switch _ty {
		case "string":
			typedInputValuesSlice[idx] = inputAsSlice[idx]
		case "address":
			typedInputValuesSlice[idx] = common.HexToAddress(inputAsSlice[idx])
		case "uint", "uint256":
			typedValue, ok := new(big.Int).SetString(inp, 10)
			if !ok {
				panic(fmt.Sprintf("Error converting %q  of type %s to big.Int", inp, _ty))
			}
			typedInputValuesSlice[idx] = typedValue
		case "uint8", "uint16", "uint32", "uint64":
			var bitsize int
			if _ty == "uint8" {
				bitsize = 8
			} else {
				// convert the last two characters into int and assign it to bitsize
				bsize, err := strconv.Atoi(_ty[4:])
				if err != nil {
					panic(fmt.Sprintf("Unsupported bitsize %q in %s", bsize, _ty))
				}
				bitsize = bsize
			}

			typedValue, err := strconv.ParseUint(inp, 10, bitsize)
			if err != nil {
				panic(fmt.Sprintf("Error converting %q  of type %s to big.Int", inp, _ty))
			}
			if bitsize == 8 {
				typedInputValuesSlice[idx] = uint8(typedValue)
			}
			if bitsize == 16 {
				typedInputValuesSlice[idx] = uint16(typedValue)
			}
			if bitsize == 32 {
				typedInputValuesSlice[idx] = uint32(typedValue)
			}
			if bitsize == 64 {
				typedInputValuesSlice[idx] = uint64(typedValue)
			}
		default:
			panic(fmt.Sprintf("Unsupported type %q", _ty))
		}

	}

	var args abi.Arguments = dataTypesToAbiArgs(dataTypes)
	values, err := args.PackValues(typedInputValuesSlice)
	if err != nil {
		panic(fmt.Sprintf("Error packing input values: %v", err))
	}

	res = hexutil.Encode(values)
	return res
}
