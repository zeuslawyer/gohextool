package encdec

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

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
		if _typeName == "int" {
			fmt.Printf("...type %q converted to int256\n", _typeName)
			_typeName = "int256"
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
