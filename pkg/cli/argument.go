package cli

import (
	"regexp"

	"github.com/StevenCyb/GoCLI/pkg/cli/internal/options"
	"github.com/StevenCyb/GoCLI/pkg/cli/internal/restriction"
	"github.com/StevenCyb/GoCLI/pkg/cli/internal/utils"
)

type argument struct {
	restriction.IsCliOption
	restriction.IsCommandOption
	restriction.IsArgumentOption
	validate *regexp.Regexp

	name        string
	example     *string
	description *string
	handler     *HandlerFunc
	command     []*command
	argument    *argument
	options     []*option
}

// CLI, Command
func Argument(name string, opts ...restriction.IsArgumentOption) *argument {
	a := &argument{
		name: name,
	}

	for _, opt := range opts {
		switch v := opt.(type) {
		case *options.Example:
			a.example = &v.Example
		case *options.Description:
			a.description = &v.Description
		case *options.Handler:
			if handlerFunc, ok := v.Handler.(HandlerFunc); ok {
				a.handler = &handlerFunc
			} else {
				panic("Invalid type for Handler option")
			}
		case *command:
			if a.argument != nil {
				panic(MixOfArgumentAndCommandError(v.name))
			}
			for _, c := range a.command {
				if c.name == v.name {
					panic(DuplicateCommandError(c.name))
				}
			}
			a.command = append(a.command, v)
		case *argument:
			if len(a.command) > 0 {
				panic(MixOfArgumentAndCommandError(v.name))
			}
			a.argument = v
		case *option:
			a.options = append(a.options, v)
		case *options.Validate:
			a.validate = v.Validate
		default:
			panic("unsupported option type")
		}
	}

	return a
}

func (c *argument) call(args *utils.AdvancedArray[string], ctx *Context) error {
	if argValue, exists := args.Next(); exists {

		if argValue == "--help" || argValue == "-h" {
			return &HelpError{on: c, backtrack: " " + c.name}
		}

		if c.validate != nil && !c.validate.MatchString(argValue) {
			return &InvalidValueError{
				on:    c.name,
				value: argValue,
			}
		}
		ctx.arguments[c.name] = argValue

		if argValue == "--help" || argValue == "-h" {
			return &HelpError{on: c, backtrack: " " + c.name}
		}

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
