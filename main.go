package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
	"github.com/zeuslawyer/hextool/encdec"
	"github.com/zeuslawyer/hextool/internal/flags"
	"github.com/zeuslawyer/hextool/selector"
)

const (
	testStringHex = "0x22596f75277665206265656e2048657865642122"
	testBigIntHex = "0xd431"
)

func main() {
	app := cli.NewApp()
	app.Name = "hextool"
	app.Description = "A cli devtool to help you encode and decode hex values for Ethereum and EVM based chains."
	app.Commands = []*cli.Command{
		{
			Name:    "tostring",
			Aliases: []string{"getstring"},
			Usage:   "decode a hex string to string",
			Action: func(cliCtx *cli.Context) error {
				fmt.Printf("%v\n", encdec.DecodeHexToString(cliCtx.String("hex")))
				return nil
			},
			Flags: []cli.Flag{
				flags.CommandFlags["hex"],
			},
		},
		{
			Name:    "toint",
			Aliases: []string{"getint"},
			Usage:   "decode a hex string to int",
			Action: func(cliCtx *cli.Context) error {
				fmt.Printf("%v\n", encdec.DecodeHexToBigInt(cliCtx.String("hex")))
				return nil
			},
			Flags: []cli.Flag{
				flags.CommandFlags["hex"],
			},
		},
		{
			Name:    "selector",
			Args:    false,
			Aliases: []string{"selectorFromSig"},
			Usage:   "calculates the function selector from a given function signature.",
			Action: func(cliCtx *cli.Context) error {
				fmt.Printf("%v\n", selector.SelectorFromSig(cliCtx.String("sig")))
				return nil
			},
			Flags: []cli.Flag{
				flags.CommandFlags["sig"],
			},
		},
		{
			Name:    "decodeMethodSelector",
			Aliases: []string{"methodsig"},
			Usage:   "Look through the provided ABI to find a function signature that matches the given function selector",
			Action: func(cliCtx *cli.Context) error {
				fmt.Printf("%v\n", selector.SigFromSelector(
					cliCtx.String("selector"),
					cliCtx.String("path"),
					cliCtx.String("url")),
				)
				return nil
			},
			Flags: []cli.Flag{
				flags.CommandFlags["selector"],
				flags.CommandFlags["path"],
				flags.CommandFlags["url"],
			},
		},
		{
			Name:    "decodeErrorSelector",
			Aliases: []string{"errorSig"},
			Usage:   "Look through the provided ABI to find the error signature that matches the given error selector",
			Action: func(cliCtx *cli.Context) error {
				fmt.Printf("%v\n", selector.ErrorSigFromSelector(
					cliCtx.String("selector"),
					cliCtx.String("path"),
					cliCtx.String("url")),
				)
				return nil
			},
			Flags: []cli.Flag{
				flags.CommandFlags["selector"],
				flags.CommandFlags["path"],
				flags.CommandFlags["url"],
			},
		},
		{
			Name:    "decodeEvent",
			Aliases: []string{"eventsig"},
			Usage:   "Look through the provided ABI to find the event signature that matches the given 32 byte topic hash",
			Action: func(cliCtx *cli.Context) error {
				fmt.Printf("%v\n", selector.EventFromTopicHash(
					cliCtx.String("topic"),
					cliCtx.String("path"),
					cliCtx.String("url")),
				)
				return nil
			},
			Flags: []cli.Flag{
				flags.CommandFlags["topic"],
				flags.CommandFlags["path"],
				flags.CommandFlags["url"],
			},
		},
		{
			Name:    "abi.decode",
			Aliases: []string{"abidecode"},
			Usage:   "abi-decode the given input hex into its corresponding data as per the comma-separated types provided",
			Action: func(cliCtx *cli.Context) error {
				fmt.Printf("%v\n", encdec.AbiDecode(
					cliCtx.String("hex"),
					cliCtx.String("types"),
				))
				return nil
			},
			Flags: []cli.Flag{
				flags.CommandFlags["hex"],
				flags.CommandFlags["types"],
			},
		},
		{
			Name:    "abi.encode",
			Aliases: []string{"abiencode"},
			Usage:   "abi-encode the given input values into hex, as per the data types provided. Input values and data types be comma-separated",
			Action: func(cliCtx *cli.Context) error {
				fmt.Printf("%v\n", encdec.AbiEncode(
					cliCtx.String("values"),
					cliCtx.String("types"),
				))
				return nil
			},
			Flags: []cli.Flag{
				flags.CommandFlags["values"],
				flags.CommandFlags["types"],
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
