package cli

import (
	"io"
	"os"
	"strings"

	"github.com/StevenCyb/GoCLI/pkg/cli/internal/options"
	"github.com/StevenCyb/GoCLI/pkg/cli/internal/restriction"
	"github.com/StevenCyb/GoCLI/pkg/cli/internal/utils"
)

type CLI struct {
	name        *string
	version     *string
	banner      *string
	example     *string
	description *string
	stdout      io.Writer
	stderr      io.Writer
	handler     *HandlerFunc
	command     []*command
	argument    *argument
	options     []*option
}

func New(opts ...restriction.IsCliOption) *CLI {
	cli := &CLI{
		stdout:  os.Stdout,
		stderr:  os.Stderr,
		command: []*command{},
		options: make([]*option, 0),
	}

	for _, opt := range opts {
		switch v := opt.(type) {
		case *options.Name:
			cli.name = &v.Name
		case *options.Version:
			cli.version = &v.Version
		case *options.Banner:
			cli.banner = &v.Banner
		case *options.Example:
			cli.example = &v.Example
		case *options.Description:
			cli.description = &v.Description
		case *options.Stream:
			cli.stdout = v.Stdout
			cli.stderr = v.Stderr
		case *options.Handler:
			if handlerFunc, ok := v.Handler.(HandlerFunc); ok {
				cli.handler = &handlerFunc
			} else {
				panic("Invalid type for Handler option")
			}
		case *command:
			if cli.argument != nil {
				panic(MixOfArgumentAndCommandError(v.name))
			}
			for _, c := range cli.command {
				if c.name == v.name {
					panic(DuplicateCommandError(c.name))
				}
			}
			cli.command = append(cli.command, v)
		case *argument:
			if len(cli.command) > 0 {
				panic(MixOfArgumentAndCommandError(v.name))
			}
			cli.argument = v
		case *option:
			cli.options = append(cli.options, v)
		default:
			panic("Unsupported option type")
		}
	}

	return cli
}

func (cli *CLI) Run() (*Context, error) {
	return cli.RunWith(os.Args)
}

func (cli *CLI) MustRun() *Context {
	ctx, err := cli.Run()
	if err != nil {
		cli.stderr.Write([]byte(err.Error() + "\n"))
		cli.PrintHelp()
		os.Exit(1)
	}

	return ctx
}

func (cli *CLI) MustRunWith(argsRaw []string) *Context {
	ctx, err := cli.RunWith(argsRaw)
	if err != nil {
		cli.stderr.Write([]byte(err.Error() + "\n"))
		cli.PrintHelp()
		os.Exit(1)
	}

	return ctx
}

func (c *CLI) RunWith(argsRaw []string) (*Context, error) {
	ctx := NewContext()
	args := utils.NewAdvancedArray(argsRaw[1:])

	if argValue, exists := args.Next(); exists {
		args.Back()
		if argValue == "--help" || argValue == "-h" {
			c.printHelp(nil)
		}
	}

	if c.argument != nil {
		if err := c.argument.call(args, ctx); err != nil {
			if helpErr, ok := err.(*HelpError); ok {
				c.printHelp(helpErr)
			}
			return nil, err
		}

		return ctx, nil
	}

	if len(c.command) > 0 {
		for _, cmd := range c.command {
			if err := cmd.call(args, ctx); err != nil && err != ErrNotMatched {
				if helpErr, ok := err.(*HelpError); ok {
					c.printHelp(helpErr)
				}
				return nil, err
			} else if err == nil {
				return ctx, nil
			}
		}

		if argValue, exists := args.Next(); exists {
			return nil, UnknownCommandError(argValue)
		}

		return nil, ErrUnexpectedEndCommand
	}

	for _, opt := range c.options {
		if err := opt.call(args, ctx); err != nil && err != ErrNotMatched {
			if _, ok := err.(*HelpError); ok {
				c.printHelp(nil)
			}
			return nil, err
		}
	}

	if c.handler != nil {
		if err := (*c.handler)(ctx); err != nil {
			return nil, err
		}
	}

	return ctx, nil
}

// PrintHelp prints the help message for the CLI application.
// It also exits afterwards with a status code of 0.
func (c *CLI) PrintHelp() {
	c.printHelp(nil)
}

func (c *CLI) printHelp(helpError *HelpError) {
	sb := strings.Builder{}
	var (
		name        = "cli"
		backtrack   string
		description *string
		example     *string
		arg         *argument
		commands    []*command
		options     []*option
	)

	if c.name != nil {
		name = *c.name
	}

	if helpError == nil || helpError.on == nil {
		description = c.description
		arg = c.argument
		commands = c.command
		options = c.options
		example = c.example
	} else {
		switch v := helpError.on.(type) {
		case *argument:
			backtrack = helpError.backtrack
			description = v.description
			arg = v
			commands = v.command
			options = v.options
			example = v.example
		case *command:
			backtrack = helpError.backtrack
			description = v.description
			arg = v.argument
			commands = v.command
			options = v.options
			example = v.example
		}
	}

	if c.banner != nil {
		sb.WriteString(*c.banner + "\n\n")
	}
	if c.version != nil {
		sb.WriteString(*c.version + "\n\n")
	}
	if description != nil {
		sb.WriteString(*description + "\n\n")
	}

	sb.WriteString("Usage: \n\t" + name + backtrack)
	if arg != nil {
		sb.WriteString(" <argument>")
	} else if len(commands) > 0 {
		sb.WriteString(" <command>\n\n")
		sb.WriteString("Commands:\n")
		for _, cmd := range commands {
			sb.WriteString("\t" + cmd.name + "\n")
			if cmd.description != nil {
				sb.WriteString("\t\t" + *cmd.description + "\n")
			}
			if cmd.example != nil {
				sb.WriteString("\t\tExample: " + *cmd.example + "\n")
			}
		}
	}
	if len(options) > 0 {
		sb.WriteString(" [options...]\n\n")
		sb.WriteString("Options:\n")
		for _, opt := range options {
			sb.WriteString("\t--" + opt.long)
			if opt.short != nil {
				sb.WriteString(", -" + string(*opt.short))
			}
			if opt.description != nil {
				sb.WriteString("\n\t\t" + *opt.description + "\n")
			}
		}
	}

	if example != nil {
		sb.WriteString("\n\nExample:\n")
		sb.WriteString("\t" + *example + "\n")
	}

	sb.WriteString("\n\nUse \"" + name + " <command> --help\" for more information about a command.\n\n")

	io.WriteString(c.stdout, sb.String())

	os.Exit(0)
}
