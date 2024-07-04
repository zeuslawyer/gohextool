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
func AbiEncode(inputValuesStr string, dataTypesStr string) (res string) {
	if len(inputValuesStr) == 0 {
		res = "0x"
		return res
	}

	// Split  each input into a slice of string values
	inputValuesSlice := strings.Split(inputValuesStr, ",")
	dataTypesSlice := strings.Split(dataTypesStr, ",")
	if len(inputValuesSlice) != len(dataTypesSlice) {
		panic(fmt.Sprintf("Number of input values does not match number of types - %d inputs to  %d types", len(inputValuesSlice), len(dataTypesSlice)))

	}

	// convert strings to the appropriate types
	// TODO zubin resume here. Add other types. Complete the tests.
	typedInputValuesSlice := make([]any, len(dataTypesSlice))

	for idx, ty := range dataTypesSlice {
		_ty := strings.TrimSpace(ty)
		inpValue := inputValuesSlice[idx]

		switch _ty {
		case "string":
			typedInputValuesSlice[idx] = inputValuesSlice[idx]
		case "address":
			typedInputValuesSlice[idx] = common.HexToAddress(inputValuesSlice[idx])

		// ints and  uints greater than 64 bits need special treatment using BigInt.
		case "uint", "uint128", "uint256", "int", "int128", "int256":
			typedValue, ok := new(big.Int).SetString(inpValue, 10)
			if !ok {
				panic(fmt.Sprintf("Error converting %q  of type %s to big.Int", inpValue, _ty))
			}
			typedInputValuesSlice[idx] = typedValue

		case "uint8", "uint16", "uint32", "uint64":
			var bitsize int

			// convert the last two characters into int and assign it to bitsize
			bsize, err := strconv.Atoi(_ty[4:])
			if err != nil {
				panic(fmt.Sprintf("Unsupported bitsize %q in %s", bsize, _ty))
			}
			bitsize = bsize

			typedValue, err := strconv.ParseUint(inpValue, 10, bitsize)
			if err != nil {
				panic(fmt.Sprintf("Error converting %q  of type %s to uint%d:  %s", inpValue, _ty, bitsize, err))
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
		case "int8", "int16", "int32", "int64":
			var bitsize int

			// convert the last two characters into int and assign it to bitsize
			bsize, err := strconv.Atoi(_ty[3:])
			if err != nil {
				panic(fmt.Sprintf("Unsupported bitsize %q in %s", bsize, _ty))
			}
			bitsize = bsize

			typedValue, err := strconv.ParseInt(inpValue, 10, bitsize)
			if err != nil {
				panic(fmt.Sprintf("Error converting %q  of type %s to int%d", inpValue, _ty, bitsize))
			}
			if bitsize == 8 {
				typedInputValuesSlice[idx] = int8(typedValue)
			}
			if bitsize == 16 {
				typedInputValuesSlice[idx] = int16(typedValue)
			}
			if bitsize == 32 {
				typedInputValuesSlice[idx] = int32(typedValue)
			}
			if bitsize == 64 {
				typedInputValuesSlice[idx] = int64(typedValue)
			}

		default:
			panic(fmt.Sprintf("Unsupported type %q", _ty))
		}

	}

	var args abi.Arguments = dataTypesToAbiArgs(dataTypesStr)
	values, err := args.PackValues(typedInputValuesSlice)
	if err != nil {
		panic(fmt.Sprintf("Error packing input values: %v", err))
	}

	res = hexutil.Encode(values)
	return res
}
