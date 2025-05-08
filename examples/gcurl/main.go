package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/StevenCyb/GoCLI/pkg/cli"
)

func main() {
	c := cli.New(
		cli.Name("gurl"),
		cli.Banner(`╔═╗╔═╗┬ ┬┬─┐┬
║ ╦║  │ │├┬┘│
╚═╝╚═╝└─┘┴└─┴─┘`),
		cli.Description("A simple CLI for making HTTP requests"),
		cli.Version("1.0.0"),
		cli.Command(
			"get", cli.Argument(
				"url",
				cli.Validate(regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)),
				cli.Description("The URL to get"),
				cli.Handler(
					func(ctx *cli.Context) error {
						url := ctx.GetArgument("url")
						fmt.Printf("Perform [GET] %s\n", *url)
						return nil
					},
				),
				cli.Option(
					"verbose",
					cli.Short('v'),
					cli.Default(""),
				),
			),
			cli.Description("Get a resource"),
			cli.Example("cli get http://example.com"),
		),
		cli.Command(
			"post", cli.Argument(
				"url",
				cli.Validate(regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)),
				cli.Description("The URL to post"),
				cli.Handler(
					func(ctx *cli.Context) error {
						url := ctx.GetArgument("url")
						bodyFile := ctx.GetOption("verbose")
						fmt.Printf("Perform [POST] %s\n", *url)
						if bodyFile != nil {
							fmt.Printf("\t With %s\n", *bodyFile)
						}
						return nil
					},
				),
				cli.Option(
					"body_file",
					cli.Short('b'),
				),
				cli.Option(
					"verbose",
					cli.Short('v'),
					cli.Default(""),
				),
			),
			cli.Description("Get a resource"),
			cli.Example("cli get http://example.com"),
		),
		// ...
		cli.Command(
			"version", cli.Handler(
				func(_ *cli.Context) error {
					fmt.Println("1.0.0")
					return nil
				},
			),
			cli.Description("Get the version of the CLI"),
		),
	)

	_, err := c.RunWith(os.Args)
	if err != nil {
		fmt.Println(err)
		c.PrintHelp()
		os.Exit(1)
	}
}
