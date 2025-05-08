package cli

import (
	"regexp"
	"testing"

	"github.com/StevenCyb/GoCLI/pkg/cli/internal/options"
	"github.com/StevenCyb/GoCLI/pkg/cli/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	t.Run("BasicCreation", func(t *testing.T) {
		t.Parallel()

		cmd := Command("test")
		assert.Equal(t, "test", cmd.name)
		assert.Nil(t, cmd.example)
		assert.Nil(t, cmd.description)
		assert.Nil(t, cmd.handler)
		assert.Empty(t, cmd.command)
		assert.Nil(t, cmd.argument)
		assert.Empty(t, cmd.options)
	})

	t.Run("WithExampleAndDescription", func(t *testing.T) {
		t.Parallel()

		example := "example usage"
		description := "command description"
		cmd := Command("test",
			&options.Example{Example: example},
			&options.Description{Description: description},
		)

		assert.Equal(t, "test", cmd.name)
		assert.NotNil(t, cmd.example)
		assert.Equal(t, example, *cmd.example)
		assert.NotNil(t, cmd.description)
		assert.Equal(t, description, *cmd.description)
	})

	t.Run("WithHandler", func(t *testing.T) {
		t.Parallel()

		handler := func(ctx *Context) error { return nil }
		cmd := Command("test", &options.Handler{Handler: HandlerFunc(handler)})

		assert.NotNil(t, cmd.handler)
		assert.NotNil(t, *cmd.handler)
	})

	t.Run("WithSubcommands", func(t *testing.T) {
		t.Parallel()

		subCmd := Command("sub")
		cmd := Command("test", subCmd)

		assert.Len(t, cmd.command, 1)
		assert.Equal(t, "sub", cmd.command[0].name)
	})

	t.Run("DuplicateSubcommands", func(t *testing.T) {
		t.Parallel()

		subCmd := Command("sub")
		assert.PanicsWithError(t, "duplicate command: sub", func() {
			Command("test", subCmd, subCmd)
		})
	})

	t.Run("WithArgument", func(t *testing.T) {
		t.Parallel()

		arg := &argument{name: "arg"}
		cmd := Command("test", arg)

		assert.NotNil(t, cmd.argument)
		assert.Equal(t, "arg", cmd.argument.name)
	})

	t.Run("MixArgumentAndSubcommands", func(t *testing.T) {
		t.Parallel()

		arg := &argument{name: "arg"}
		subCmd := Command("sub")
		assert.PanicsWithError(t, "mix of argument and command: sub", func() {
			Command("test", arg, subCmd)
		})
	})

	t.Run("WithOptions", func(t *testing.T) {
		t.Parallel()

		opt := &option{long: "opt"}
		cmd := Command("test", opt)

		assert.Len(t, cmd.options, 1)
		assert.Equal(t, "opt", cmd.options[0].long)
	})

	t.Run("MixOfArgumentAndCommand", func(t *testing.T) {
		t.Parallel()

		arg := &argument{name: "arg"}
		subCmd := Command("sub")
		assert.PanicsWithError(t, "mix of argument and command: arg", func() {
			Command("test", subCmd, arg)
		})
	})

	t.Run("InvalidHandlerType", func(t *testing.T) {
		t.Parallel()

		assert.PanicsWithValue(t, "Invalid type for Handler option", func() {
			Command("test", &options.Handler{Handler: "invalid"})
		})
	})
}

func TestCommandCall(t *testing.T) {
	t.Parallel()

	t.Run("NotMatching", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test"})
		ctx := NewContext()
		cmd := Command("nope")
		err := cmd.call(args, ctx)

		assert.Error(t, err)
		assert.Equal(t, ErrNotMatched, err)
	})

	t.Run("MatchWithoutTrilling", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test"})
		ctx := NewContext()
		cmd := Command("test")
		err := cmd.call(args, ctx)

		assert.NoError(t, err)
		assert.Equal(t, []string{"test"}, ctx.commands)
	})

	t.Run("MatchWithFollowingArgument", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test", "hello"})
		ctx := NewContext()
		cmd := Command("test", Argument("arg"))
		err := cmd.call(args, ctx)

		assert.NoError(t, err)
		assert.Equal(t, []string{"test"}, ctx.commands)
		assert.Equal(t, map[string]string{"arg": "hello"}, ctx.arguments)
	})

	t.Run("MatchWithSubcommand", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test", "sub"})
		ctx := NewContext()
		subCmd := Command("sub")
		cmd := Command("test", subCmd)
		err := cmd.call(args, ctx)

		assert.NoError(t, err)
		assert.Equal(t, []string{"test", "sub"}, ctx.commands)
	})

	t.Run("MatchWithNonMatchingSubcommand", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test", "sub"})
		ctx := NewContext()
		subCmd := Command("nope")
		cmd := Command("test", subCmd)
		err := cmd.call(args, ctx)

		assert.Error(t, err)
		assert.Equal(t, UnknownCommandError("sub"), err)
	})

	t.Run("MatchWithSubcommandAndUnexpectedEnd", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test", "sub"})
		ctx := NewContext()
		subCmd2 := Command("sub2")
		subCmd1 := Command("sub", subCmd2)
		cmd := Command("test", subCmd1)
		err := cmd.call(args, ctx)

		assert.Error(t, err)
		assert.Equal(t, ErrUnexpectedEndCommand, err)
	})

	t.Run("WithOptions", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test", "--opt", "value", "-b", "value2"})
		ctx := NewContext()
		b := 'b'
		opt1 := &option{long: "aa"}
		opt2 := &option{long: "opt"}
		opt3 := &option{long: "bb", short: &b}
		cmd := Command("test", opt1, opt2, opt3)
		err := cmd.call(args, ctx)

		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"opt": "value", "bb": "value2"}, ctx.options)
	})

	t.Run("WithHandler", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test"})
		ctx := NewContext()
		handlerCalled := false
		cmd := Command("test", &options.Handler{
			Handler: HandlerFunc(func(ctx *Context) error {
				handlerCalled = true
				return nil
			}),
		})
		err := cmd.call(args, ctx)

		assert.NoError(t, err)
		assert.True(t, handlerCalled)
	})

	t.Run("WithHandlerError", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test"})
		ctx := NewContext()
		handlerCalled := false
		cmd := Command("test", &options.Handler{
			Handler: HandlerFunc(func(ctx *Context) error {
				handlerCalled = true
				return ErrNotMatched
			}),
		})
		err := cmd.call(args, ctx)

		assert.Error(t, err)
		assert.Equal(t, ErrNotMatched, err)
		assert.True(t, handlerCalled)
	})

	t.Run("WithOptionFailingValidation", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test", "--opt", "value", "-b", "value2"})
		ctx := NewContext()
		b := 'b'
		opt1 := &option{long: "aa"}
		opt2 := &option{long: "opt", validate: regexp.MustCompile("^v[0-9]+$")}
		opt3 := &option{long: "bb", short: &b}
		cmd := Command("test", opt1, opt2, opt3)
		err := cmd.call(args, ctx)

		assert.Error(t, err)
		assert.Equal(t, &InvalidValueError{
			on:    "opt",
			value: "value",
		}, err)
	})

	t.Run("SubCommandWithHelpOption", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test", "sub", "--help"})
		ctx := NewContext()
		subCmd := Command("sub")
		cmd := Command("test", subCmd)

		err := cmd.call(args, ctx)

		assert.Error(t, err)
		assert.IsType(t, &HelpError{}, err)
	})
}
