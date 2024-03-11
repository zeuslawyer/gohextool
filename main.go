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
			Aliases: []string{"selectorFromSig"},
			Usage:   "calculates the function selector from a given function signature",
			Action: func(cliCtx *cli.Context) error {
				fmt.Printf("%v\n", selector.SelectorFromSig(cliCtx.String("sig")))
				return nil
			},
			Flags: []cli.Flag{
				flags.CommandFlags["sig"],
			},
		},
		{
			Name:    "funcsig",
			Aliases: []string{"matchToSig"},
			Usage:   "Look through the provided ABI to find a function signature that matches the given function selector",
			Action: func(cliCtx *cli.Context) error {
				// TODO @zeuslawyer instead of flags the path and url should be arguments?
				fmt.Printf("%v\n", selector.SigFromSelector(cliCtx.String("selector"), cliCtx.String("path"), cliCtx.String("url")))
				return nil
			},
			Flags: []cli.Flag{
				flags.CommandFlags["selector"],
				flags.CommandFlags["path"],
				flags.CommandFlags["url"],
			},
		},

		// TODO zeuslawyer: add selector commands
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
