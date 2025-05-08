package main

import (
	"fmt"
	"os"

	"github.com/StevenCyb/GoCLI/pkg/cli"
)

func main() {
	c := cli.New(
		cli.Name("converter"),
		cli.Option(
			"timezone",
			cli.Short('t'),
			cli.Required(),
		),
		cli.Handler(func(ctx *cli.Context) error {
			fmt.Println("Current timezone:", ctx.GetOption("timezone"))
			// ...
			return nil
		}),
	)

	_, err := c.RunWith(os.Args)
	if err != nil {
		fmt.Println(err)
		c.PrintHelp()
		os.Exit(1)
	}
}
