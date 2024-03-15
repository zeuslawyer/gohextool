package flags

import "github.com/urfave/cli/v2"

// Create a map of flags with keys as the flag name and values as the cli.Flag type
var CommandFlags = make(map[string]cli.Flag)

func init() {
	CommandFlags["hex"] = &cli.StringFlag{
		Name:  "hex",
		Value: "0x",
		Usage: "hex string to convert",
	}

	CommandFlags["selector"] = &cli.StringFlag{
		Name:  "selector",
		Usage: "Function Selector hex string",
	}
	CommandFlags["path"] = &cli.StringFlag{
		Name:  "path",
		Usage: "absolute path to the ABI file",
	}
	CommandFlags["url"] = &cli.StringFlag{
		Name:  "url",
		Usage: "public API endpoint from where to fetch the object containing the abi property",
	}
	CommandFlags["sig"] = &cli.StringFlag{
		Name:  "sig",
		Usage: "Function signature in quotes. Exclude the the 'function' keyword. Must follow the ABI spec e.g.  'function foo(uint32 a, int b)' = 'foo(uint32,int256)'",
	}
}
