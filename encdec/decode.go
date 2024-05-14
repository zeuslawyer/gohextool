package encdec

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
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
// TODO zubin  try doing the opposite..Arguments.Pack()
func AbiDecode(hexInput string, dataTypes string) []any {
	if strings.HasPrefix(hexInput, "0x") {
		// Strip out the "0x" prefix
		hexInput = hexInput[2:]
	}

	var args abi.Arguments = parseDataTypesString(dataTypes)

	b, err := hex.DecodeString(hexInput)
	if err != nil {
		fmt.Printf("Error decoding hex input: %v", err)
	}

	// fmt.Printf("LOOK ZUBIN:  %d && %v\n", len(b), args.NonIndexed())
	// if len(b) < 64 {
	// 	// to avoid this error in toGoType() https://github.com/ethereum/go-ethereum/blob/master/accounts/abi/unpack.go#L224
	// 	// called inside https://github.com/ethereum/go-ethereum/blob/8f7eb9ccd99ee47721a7bfde494d6da28de4cd8e/accounts/abi/argument.go#L184
	// 	// b = common.RightPadBytes(b, 64)
	// }
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

// TODO zubin clean this up and its helper func.
func TODOAbiDecode(hexInput string, dataTypes []string) any {
	// https://pkg.go.dev/github.com/ethereum/go-ethereum/rlp#Decode
	// References: https://ethereum.stackexchange.com/questions/117060/abi-decode-raw-types-with-go
	//https://ethereum.stackexchange.com/questions/29809/how-to-decode-input-data-with-abi-using-golang

	if strings.HasPrefix(hexInput, "0x") {
		// Strip out the "0x" prefix
		hexInput = hexInput[2:]
	}
	decodedBytes, err := hex.DecodeString(hexInput)
	if err != nil {
		fmt.Errorf("Error decoding hex input: %v", err)
	}

	// Convert each string data type to a reflect.Type
	var structTypes []reflect.Type
	for _, typeName := range dataTypes {
		t, _ := abi.NewType(typeName, "", nil)
		structTypes = append(structTypes, t.GetType())
	}

	// Create a new struct type dynamically at runtime
	// to hold the decoded values.
	dynamicallyCreatedStruct := reflect.StructOf(makeStructFields(structTypes))

	// Create a new instance of the struct
	structInstancePtr := reflect.New(dynamicallyCreatedStruct).Elem()
	var resultValues = structInstancePtr.Interface()

	err = rlp.Decode(bytes.NewReader(decodedBytes), &resultValues)
	if err != nil {
		fmt.Printf("Error RLP decoding the hex input: %#v\n", err)
		panic(err)
	} else {
		fmt.Printf("Decoded value: %#v\n", resultValues)
	}
	return resultValues
}

func makeStructFields(fieldTypes []reflect.Type) []reflect.StructField {
	var structFields []reflect.StructField
	for i, fieldType := range fieldTypes {
		field := reflect.StructField{
			Name: fmt.Sprintf("Field%d", i),
			Type: fieldType,
		}
		structFields = append(structFields, field)

	}
	return structFields
}

func parseDataTypesString(dataTypes string) abi.Arguments {
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
