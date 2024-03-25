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
func FuncFromSelector(selector string, abiPath string, abiUrl string) string {
	if abiPath == "" && abiUrl == "" {
		panic(fmt.Errorf("abiPath and url cannot both be empty"))
	}
	var abiJsonStr string
	var abiSource string

	if abiPath != "" {
		err := validateUriExtension(abiPath)
		if err != nil {
			panic(err)
		}

		abiSource = abiPath
		fileBytes, err := os.ReadFile(abiSource)
		if err != nil {
			fmt.Println("Error reading file")
			panic(err)
		}

		abiJsonStr = bytesToJsonString(fileBytes, abiSource)
	} else { // reading from URL instead of file
		err := validateUriExtension(abiUrl)
		if err != nil {
			panic(err)
		}

		abiSource = abiUrl
		resp, err := http.Get(abiSource)
		if err != nil {
			fmt.Printf("Error fetching ABI file from url %s", abiSource)
			panic(err)
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading file from http response")
			panic(err)
		}

		abiJsonStr = bytesToJsonString(b, abiSource)
	}

	parsedAbi, err := abi.JSON(strings.NewReader(abiJsonStr))
	if err != nil { // @zeuslawyer TODO check if this is the correct way to check for this error
		fmt.Printf("Error parsing ABI from : %s. \nABIs provided must be an array.\n", abiSource)
		panic(err)
	}

	selectorBytes := hexutil.MustDecode(selector)
	method, err := parsedAbi.MethodById(selectorBytes)
	if err != nil {
		fmt.Println("Error looking up method by its ID")
		panic(err)
	}
	if method == nil {
		return fmt.Sprintf("Method not found in file at %s", abiPath)
	}

	return method.Sig
}

// Given an Events Topic Hash (32 bytes), returns the event's signature from provided ABI file and path
// or from a URL.  If both are provided it will default to using the file path.
func EventFromTopicHash(topicHash string, abiPath string, abiUrl string) string {
	if abiPath == "" && abiUrl == "" {
		panic(fmt.Errorf("abiPath and url cannot both be empty"))
	}

	var abiJsonStr string
	var abiSource string
	if abiPath != "" { // TODO @zeuslawyer abstract logic from both the function and event decoders.
		err := validateUriExtension(abiPath)
		if err != nil {
			panic(err)
		}

		abiSource = abiPath
		fileBytes, err := os.ReadFile(abiSource)
		if err != nil {
			fmt.Println("Error reading file")
			panic(err)
		}

		abiJsonStr = bytesToJsonString(fileBytes, abiSource)

	} else { // reading from URL instead of file
		err := validateUriExtension(abiUrl)
		if err != nil {
			panic(err)
		}

		abiSource = abiUrl
		resp, err := http.Get(abiSource)
		if err != nil {
			fmt.Printf("Error fetching ABI file from url %s", abiSource)
			panic(err)
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading file from http response")
			panic(err)
		}

		abiJsonStr = bytesToJsonString(b, abiSource)
	}

	parsedAbi, err := abi.JSON(strings.NewReader(abiJsonStr))
	if err != nil { // @zeuslawyer TODO check if this is the correct way to check for this error
		fmt.Printf("Error parsing ABI from : %s. \nABIs provided must be an array.\n", abiSource)
		panic(err)
	}

	selectorBytes := hexutil.MustDecode(topicHash)
	selectorHash := common.BytesToHash(selectorBytes)
	ev, err := parsedAbi.EventByID(selectorHash)
	if err != nil {
		fmt.Println("Error looking up event by its topics hash")
		panic(err)
	}
	if ev == nil {
		return fmt.Sprintf("Method not found in file at %s", abiPath)
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
	var data map[string]interface{}

	if err := json.Unmarshal(b, &data); err != nil {
		panic(fmt.Errorf("error parsing JSON from file at %s", abiSourceUri))
	}

	abiData, ok := data["abi"]
	if !ok {
		fmt.Printf("Property 'abi' not found in unmarshalled JSON data. Check the file at %s", abiSourceUri)
		return ""
	}

	// isSlice := validateIsSlice(abiSlice)
	// if !isSlice {
	// 	fmt.Printf("Value of property 'abi' in supplied file is not an array")
	// 	return ""
	// }

	// arrData := abiSlice.([]interface{}) // @zeuslawyer TODO since this type assertion returns OK do we need validateIsSlice?
	// jsonBytes, err := json.Marshal(arrData)
	// if err != nil {
	// 	fmt.Printf("Error marshalling ABI's array data to JSON bytes")
	// 	return ""
	// }

	abiSlice, ok := abiData.([]interface{}) // @zeuslawyer TODO since this type assertion returns OK do we need valudateIsSlice?
	if !ok {
		fmt.Printf("Value of property 'abi' in supplied file is not an array")
		return ""
	}

	jsonBytes, err := json.Marshal(abiSlice)
	if err != nil {
		fmt.Printf("Error marshalling ABI's array data to JSON bytes")
		return ""
	}

	return string(jsonBytes)
}
