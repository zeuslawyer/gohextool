package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
	"github.com/zeuslawyer/hextool/encdec"
	"github.com/zeuslawyer/hextool/internal/flags"
)

const (
	testStringHex = "0x22596f75277665206265656e2048657865642122"
	testBigIntHex = "0xd431"
)

func main() {
	app := cli.NewApp()
	app.Name = "hextool"
	app.Description = "A cli tool to help you encode and decode hex strings"
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
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
