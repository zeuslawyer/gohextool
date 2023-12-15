package main

import (
	"fmt"

	"github.com/zeuslawyer/gohextool/encdec"
)

const (
	testStringHex = "0x476f20466f727468202620436f6e717565722c20486f6d696521"
	testBigIntHex = "0xd431"
)

func main() {
	str := encdec.DecodeHexToString(testStringHex)
	fmt.Println("result string : ", str)

	bigInt := encdec.DecodeHexToBigInt(testBigIntHex)
	fmt.Println("result big int : ", bigInt)
}
