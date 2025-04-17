package cli

import (
	"GoCLI/pkg/cli/internal/options"
	"GoCLI/pkg/cli/internal/restriction"
	"GoCLI/pkg/cli/internal/utils"
	"regexp"
	"strings"
)

type option struct {
	long         string
	short        *rune
	required     bool
	defaultValue *string
	description  *string
	validate     *regexp.Regexp

	restriction.IsCliOption
	restriction.IsCommandOption
	restriction.IsArgumentOption
}

// CLI, Command
func Option(long string, opts ...restriction.IsOptionOption) *option {
	o := &option{
		long: long,
	}

	for _, opt := range opts {
		switch v := opt.(type) {
		case *options.Short:
			o.short = &v.Short
		case *options.Default:
			o.defaultValue = &v.DefaultValue
		case *options.Description:
			o.description = &v.Description
		case *options.Required:
			o.required = true
		case *options.Validate:
			o.validate = v.Validate
		default:
			panic("unsupported option type")
		}
	}

	return o
}

func (o *option) call(args *utils.AdvancedArray[string], ctx *Context) error {
	if argValue, exists := args.Next(); exists {
		if argValue == "--help" || argValue == "-h" {
			return &HelpError{on: o}
		}

		if argValue == "--"+o.long || o.short != nil && argValue == "-"+string(*o.short) {
			argValue, exists := args.Next()
			if exists && !strings.HasPrefix(argValue, "-") {
				if o.validate != nil && !o.validate.MatchString(argValue) {
					return &InvalidValueError{
						on:    o.long,
						value: argValue,
					}
				}
				ctx.options[o.long] = argValue
			} else {
				if exists {
					args.Back()
				}
				ctx.options[o.long] = ""
			}

			return nil
		}

		args.Back()
	}

	return ErrNotMatched
}
