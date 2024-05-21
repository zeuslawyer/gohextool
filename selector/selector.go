package selector

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Calculates the function selector given function signature `funcSig`.
// The function signature should be in the form of `functionName(type1,type2,...)`.
// Eg: "transfer(address,uint256)"
func SelectorFromSig(funcSig string) string {
	if funcSig == "" {
		fmt.Println("Error: function signature cannot be empty. Pass in the '--sig' flag with the function signature.")
		return ""
	}
	funcSig = strings.ReplaceAll(funcSig, " ", "")
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
	selector := funcSigHash.String()[:10] // first 4 bytes ==8 characters, plus "0x"
	return selector
}

// Given a function selector, returns the function signature from provided ABI file and path
// or from a URL.  If both are provided it will default to using the file path.
func SigFromSelector(selector string, _abiPath string, abiUrl string) string {
	return fromSelector(false, selector, _abiPath, abiUrl)
}

func ErrorSigFromSelector(selector string, _abiPath string, abiUrl string) string {
	return fromSelector(true, selector, _abiPath, abiUrl)
}

// Given an Events Topic Hash (32 bytes), returns the event's signature from provided ABI file and path
// or from a URL.  If both are provided it will default to using the file path.
func EventFromTopicHash(topicHex string, _abiPath string, abiUrl string) string {
	if _abiPath == "" && abiUrl == "" {
		panic(fmt.Errorf("abiPath and url cannot both be empty"))
	}

	var abiJsonStr string
	var abiPath string
	if _abiPath != "" {
		err := validateUriExtension(_abiPath)
		if err != nil {
			panic(err)
		}

		abiPath = _abiPath
		fileBytes, err := os.ReadFile(abiPath)
		if err != nil {
			fmt.Println("Error reading file")
			panic(err)
		}

		abiJsonStr = bytesToJsonString(fileBytes, abiPath)

	} else { // reading from URL instead of file
		err := validateUriExtension(abiUrl)
		if err != nil {
			panic(err)
		}

		abiPath = abiUrl
		resp, err := http.Get(abiPath)
		if err != nil {
			fmt.Printf("Error fetching ABI file from url %s", abiPath)
			panic(err)
		}
		defer resp.Body.Close()

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading file from http response")
			panic(err)
		}

		abiJsonStr = bytesToJsonString(respBytes, abiPath)
	}

	topicBytes := hexutil.MustDecode(topicHex)
	topicHash := common.BytesToHash(topicBytes)

	parsedAbi, err := abi.JSON(strings.NewReader(abiJsonStr))
	if err != nil { // @zeuslawyer TODO check if this is the correct way to check for this error
		fmt.Printf("Error parsing ABI from : %s. \nABI provided must be an array.\n", abiPath)
		panic(err)
	}

	ev, err := parsedAbi.EventByID(topicHash)
	if err != nil {
		fmt.Println("Error looking up event by its topics hash")
		panic(err)
	}
	if ev == nil {
		return fmt.Sprintf("Method not found in file at %s", _abiPath)
	}

	return ev.Sig
}

func validateUriExtension(uri string) error {
	ext := filepath.Ext(uri)
	if ext != ".json" {
		return fmt.Errorf("invalid file/url extension: %s, must be .json", ext)
	}
	return nil
}

func bytesToJsonString(b []byte, abiSourceUri string) string {
	var data any // TODO: zeuslawyer refactor to handle array or map. Eg of array https://raw.githubusercontent.com/Cyfrin/ccip-contracts/main/contracts-ccip/abi/v0.8/Router.json

	if err := json.Unmarshal(b, &data); err != nil {
		panic(fmt.Errorf("error parsing JSON from file at %s : %s", abiSourceUri, err))
	}

	// Check if the data is a slice (array) or a map (object)
	var abiData any
	switch v := data.(type) {
	case []interface{}:
		fmt.Println("Data is an array\n")
		// You can work with v as a []interface{}
		abiData = v
	case map[string]interface{}:
		fmt.Println("Data is an object\n")
		d, ok := v["abi"]
		if !ok {
			panic(fmt.Errorf("Property 'abi' not found in unmarshalled JSON data. Check the file at %s", abiSourceUri))
		}

		// check that the abi property is an array
		_, ok = d.([]any)
		if !ok {
			fmt.Printf("Value of property 'abi' in supplied file is not an array")
		}
		abiData = d
	default:
		panic(fmt.Errorf("Data in file at %s is neither an array nor an object", abiSourceUri))
	}

	jsonBytes, err := json.Marshal(abiData)
	if err != nil {
		fmt.Printf("Error marshalling ABI data to JSON bytes: %s", err)
		return ""
	}

	return string(jsonBytes)
}

func fromSelector(isErrorSelector bool, selector string, _abiPath string, abiUrl string) string {
	if _abiPath == "" && abiUrl == "" {
		panic(fmt.Errorf("abiPath and url cannot both be empty"))
	}
	var abiJsonStr string
	var abiPath string

	if _abiPath != "" {
		err := validateUriExtension(_abiPath)
		if err != nil {
			panic(err)
		}

		abiPath = _abiPath
		fileBytes, err := os.ReadFile(abiPath)
		if err != nil {
			fmt.Println("Error reading file")
			panic(err)
		}

		abiJsonStr = bytesToJsonString(fileBytes, abiPath)
	} else { // reading from URL instead of file
		err := validateUriExtension(abiUrl)
		if err != nil {
			panic(err)
		}

		abiPath = abiUrl
		resp, err := http.Get(abiPath)
		if err != nil {
			fmt.Printf("Error fetching ABI file from url %s", abiPath)
			panic(err)
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading file from http response")
			panic(err)
		}

		abiJsonStr = bytesToJsonString(b, abiPath)
	}

	parsedAbi, err := abi.JSON(strings.NewReader(abiJsonStr))
	if err != nil { // @zeuslawyer TODO check if this is the correct way to check for this error
		fmt.Printf("Error parsing ABI from : %s. \nABI provided must be an array.\n", abiPath)
		panic(err)
	}

	selectorBytes := hexutil.MustDecode(selector)

	var sig string
	if isErrorSelector {
		var first4Bytes [4]byte
		copy(first4Bytes[:], selectorBytes[:4])
		errorSig, e := parsedAbi.ErrorByID(first4Bytes)
		if e != nil {
			panic(fmt.Errorf("Error looking up error signature by its selector: %s", e))
		}
		if errorSig == nil {
			panic(fmt.Errorf("Error signature for selector %s not found in file at %s", selector, _abiPath))
		}
		sig = errorSig.Sig
	} else {
		methodSig, e := parsedAbi.MethodById(selectorBytes)
		if e != nil {
			panic(fmt.Errorf("Error looking up method signature by its selector: %s", e))
		}
		if methodSig == nil {
			panic(fmt.Errorf("Method signature for selector %s not found in file at %s", selector, _abiPath))
		}
		sig = methodSig.Sig
	}

	return sig
}
