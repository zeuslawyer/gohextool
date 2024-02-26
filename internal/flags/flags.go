package flags

import "github.com/urfave/cli/v2"

// Create a map of flags with keys as the flag name and values as the cli.Flag type
var CommandFlags = make(map[string]cli.Flag)

func init() {
	CommandFlags["lang"] = &cli.StringFlag{
		Name:  "lang",
		Value: "english",
		Usage: "language for the greeting",
	}

	CommandFlags["hex"] = &cli.StringFlag{
		Name:  "hex",
		Value: "0x",
		Usage: "hex string to convert",
	}
}
