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
			Aliases: []string{""},
			Usage:   "decode a hex string to string",
			Action: func(cliCtx *cli.Context) error {
				fmt.Print(encdec.DecodeHexToString(cliCtx.String("hex")))
				return nil
			},
			Flags: []cli.Flag{
				flags.CommandFlags["hex"],
			},
		},
	}

	// app := &cli.App{
	// 	Name:        "hextool - A hex string encoding and decoding tool",
	// 	Description: "A cli tool to help you encode and decode hex strings",
	// 	Commands: []*cli.Command{
	// {
	// 	Name:    "tostring",
	// 	Aliases: []string{""},
	// 	Usage:   "decode a hex string to string",
	// 	Action: func(cliCtx *cli.Context) error {
	// 		if cliCtx.String("hex") == "0x" {
	// 			log.Print("No hex string provide. Please provide a hex string using the --hex flag")
	// 		}

	// 		fmt.Print(encdec.DecodeHexToString(cliCtx.String("hex")))
	// 		return nil
	// 	},
	// 	Flags: []cli.Flag{
	// 		flags.CommandFlags["hex"],
	// 	},
	// },
	// 	},
	// }

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
