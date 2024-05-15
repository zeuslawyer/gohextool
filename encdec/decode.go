package encdec

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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

// TODO zubin  Abi.encode with,..Arguments.Pack()
/*
 * Given a tuple of data values and a tuple of their corresponding types,
 * ABI-encode the data values according to their provided types
* `input` is a comma-separated string of values. Eg: "hello, 123, true, 456".
 `dataTypes` is a comma-separated string of types. Eg: "string, uint, bool, uint".
 *  The sequence of types in `dataTypes` must match the sequence of values in `input`.
*/
func AbiEncode(inputValues string, dataTypes string) (res string) {
	if len(inputValues) == 0 {
		res = "0x"
		return res
	}

	// Split  each input into a slice of string values
	inputAsSlice := strings.Split(inputValues, ",")
	typesAsSlice := strings.Split(dataTypes, ",")

	// convert strings to the appropriate types
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

// Parse  comma-separated string containing a list of 1 or more
// data types and return Abi.Arguments.
func dataTypesToAbiArgs(dataTypes string) abi.Arguments {
	typesSlice := strings.Split(dataTypes, ",")

	// convert each input type into an Abi.Arg.
	abiArgs := make([]abi.Argument, len(typesSlice))

	for idx, typeName := range typesSlice {
		_typeName := strings.TrimSpace(typeName)
		// uints throw an error when creating an Abi.Type, so convert them to uint256.
		if _typeName == "uint" {
			fmt.Printf("...type %q converted to uint256\n", _typeName)
			_typeName = "uint256"
		}
		abiType, err := abi.NewType(_typeName, "", nil)
		if err != nil {
			panic(fmt.Sprintf("Error creating Abi.Type for %s: %v", _typeName, err))
		}
		abiArgs[idx] = abi.Argument{
			Name:    fmt.Sprintf("arg%d", idx),
			Type:    abiType,
			Indexed: false,
		}
	}

	return abiArgs
}
