package cli

import (
	"GoCLI/pkg/cli/internal/restriction"
	"errors"
)

var ErrNotMatched = errors.New("not matched")
var ErrUnexpectedEndCommand = errors.New("unexpected end of command")

type DuplicateCommandError string

func (e DuplicateCommandError) Error() string {
	return "duplicate command: " + string(e)
}

type MixOfArgumentAndCommandError string

func (e MixOfArgumentAndCommandError) Error() string {
	return "mix of argument and command: " + string(e)
}

type UnknownCommandError string

func (e UnknownCommandError) Error() string {
	return "unknown command: " + string(e)
}

type UnknownArgumentError string

func (e UnknownArgumentError) Error() string {
	return "unknown argument: " + string(e)
}

type InvalidValueError struct {
	on    string
	value string
}

func (e InvalidValueError) Error() string {
	return "invalid value for " + e.on + ": " + e.value
}

type HelpError struct {
	on        restriction.IsCliOption
	backtrack string
}

func (e HelpError) Error() string {
	return "help"
}
