package main

import (
	"fmt"
	"os"

	"github.com/StevenCyb/GoCLI/pkg/cli"
)

func main() {
	c := cli.New(
		cli.Name("converter"),
		cli.Argument(
			"input",
			cli.Argument("output",
				cli.Handler(func(ctx *cli.Context) error {
					fmt.Println("Input:", ctx.GetArgument("input"))
					fmt.Println("Output:", ctx.GetArgument("output"))
					return nil
				}),
			),
		),
	)

	_, err := c.RunWith(os.Args)
	if err != nil {
		fmt.Println(err)
		c.PrintHelp()
		os.Exit(1)
	}
}
