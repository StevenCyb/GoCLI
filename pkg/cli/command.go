package cli

import (
	"github.com/StevenCyb/GoCLI/pkg/cli/internal/options"
	"github.com/StevenCyb/GoCLI/pkg/cli/internal/restriction"
	"github.com/StevenCyb/GoCLI/pkg/cli/internal/utils"
)

type command struct {
	restriction.IsCliOption
	restriction.IsCommandOption
	restriction.IsArgumentOption

	name        string
	example     *string
	description *string
	handler     *HandlerFunc
	command     []*command
	argument    *argument
	options     []*option
}

// CLI, Command, Argument
func Command(name string, opts ...restriction.IsCommandOption) *command {
	o := &command{
		name: name,
	}

	for _, opt := range opts {
		switch v := opt.(type) {
		case *options.Example:
			o.example = &v.Example
		case *options.Description:
			o.description = &v.Description
		case *options.Handler:
			if handlerFunc, ok := v.Handler.(HandlerFunc); ok {
				o.handler = &handlerFunc
			} else {
				panic("Invalid type for Handler option")
			}
		case *command:
			if o.argument != nil {
				panic(MixOfArgumentAndCommandError(v.name))
			}
			for _, c := range o.command {
				if c.name == v.name {
					panic(DuplicateCommandError(c.name))
				}
			}
			o.command = append(o.command, v)
		case *argument:
			if len(o.command) > 0 {
				panic(MixOfArgumentAndCommandError(v.name))
			}
			o.argument = v
		case *option:
			o.options = append(o.options, v)
		default:
			panic("unsupported option type")
		}
	}

	return o
}

func (c *command) call(args *utils.AdvancedArray[string], ctx *Context) error {
	if argValue, exists := args.Next(); exists {
		if argValue != c.name {
			args.Back()
			return ErrNotMatched
		}

		if argValue, exists := args.Next(); exists {
			args.Back()
			if argValue == "--help" || argValue == "-h" {
				helpErr := &HelpError{on: c, backtrack: " " + c.name}
				return helpErr
			}
		}

		ctx.commands = append(ctx.commands, c.name)

		if c.argument != nil {
			if err := c.argument.call(args, ctx); err != nil {
				if helpErr, ok := err.(*HelpError); ok {
					helpErr.backtrack = " " + c.name + helpErr.backtrack
					return helpErr
				}
				return err
			}
			return nil
		}

		if len(c.command) > 0 {
			for _, cmd := range c.command {
				if err := cmd.call(args, ctx); err != nil && err != ErrNotMatched {
					if helpErr, ok := err.(*HelpError); ok {
						helpErr.backtrack = " " + c.name + helpErr.backtrack
						return helpErr
					}
					return err
				} else if err == nil {
					return nil
				}
			}

			arg, _ := args.Next()

			return UnknownCommandError(arg)
		}

		for _, opt := range c.options {
			if err := opt.call(args, ctx); err != nil && err != ErrNotMatched {
				if helpErr, ok := err.(*HelpError); ok {
					helpErr.on = c
					return helpErr
				}

				return err
			}
		}

		if c.handler != nil {
			if err := (*c.handler)(ctx); err != nil {
				return err
			}
		}

		return nil
	}

	return ErrUnexpectedEndCommand
}
