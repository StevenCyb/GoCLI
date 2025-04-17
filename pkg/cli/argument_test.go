package cli

import (
	"regexp"
	"testing"

	"GoCLI/pkg/cli/internal/options"
	"GoCLI/pkg/cli/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestArgument(t *testing.T) {
	t.Parallel()

	t.Run("BasicCreation", func(t *testing.T) {
		t.Parallel()

		arg := Argument("test")
		assert.Equal(t, "test", arg.name)
		assert.Nil(t, arg.example)
		assert.Nil(t, arg.description)
		assert.Nil(t, arg.handler)
		assert.Empty(t, arg.command)
		assert.Nil(t, arg.argument)
		assert.Empty(t, arg.options)
	})

	t.Run("WithExampleAndDescription", func(t *testing.T) {
		t.Parallel()

		example := "example usage"
		description := "argument description"
		arg := Argument("test",
			&options.Example{Example: example},
			&options.Description{Description: description},
		)

		assert.Equal(t, "test", arg.name)
		assert.NotNil(t, arg.example)
		assert.Equal(t, example, *arg.example)
		assert.NotNil(t, arg.description)
		assert.Equal(t, description, *arg.description)
	})

	t.Run("WithHandler", func(t *testing.T) {
		t.Parallel()

		handler := func(ctx *Context) error { return nil }
		arg := Argument("test", &options.Handler{Handler: HandlerFunc(handler)})

		assert.NotNil(t, arg.handler)
		assert.NotNil(t, *arg.handler)
	})

	t.Run("WithSubcommands", func(t *testing.T) {
		t.Parallel()

		subCmd := Command("sub")
		arg := Argument("test", subCmd)

		assert.Len(t, arg.command, 1)
		assert.Equal(t, "sub", arg.command[0].name)
	})

	t.Run("DuplicateSubcommands", func(t *testing.T) {
		t.Parallel()

		subCmd := Command("sub")
		assert.PanicsWithError(t, "duplicate command: sub", func() {
			Argument("test", subCmd, subCmd)
		})
	})

	t.Run("WithNestedArgument", func(t *testing.T) {
		t.Parallel()

		nestedArg := Argument("nested")
		arg := Argument("test", nestedArg)

		assert.NotNil(t, arg.argument)
		assert.Equal(t, "nested", arg.argument.name)
	})

	t.Run("MixArgumentAndSubcommands", func(t *testing.T) {
		t.Parallel()

		nestedArg := Argument("nested")
		subCmd := Command("sub")
		assert.PanicsWithError(t, "mix of argument and command: sub", func() {
			Argument("test", nestedArg, subCmd)
		})
	})

	t.Run("MixSubcommandsAndArgument", func(t *testing.T) {
		t.Parallel()

		nestedArg := Argument("nested")
		subCmd := Command("sub")
		assert.PanicsWithError(t, MixOfArgumentAndCommandError("nested").Error(), func() {
			Argument("test", subCmd, nestedArg)
		})
	})

	t.Run("WithOptions", func(t *testing.T) {
		t.Parallel()

		opt := &option{long: "opt"}
		arg := Argument("test", opt)

		assert.Len(t, arg.options, 1)
		assert.Equal(t, "opt", arg.options[0].long)
	})

	t.Run("InvalidHandlerType", func(t *testing.T) {
		t.Parallel()

		assert.PanicsWithValue(t, "Invalid type for Handler option", func() {
			Argument("test", &options.Handler{Handler: "invalid"})
		})
	})
}

func TestArgumentCall(t *testing.T) {
	t.Parallel()

	t.Run("UnexpectedEnd", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{})
		ctx := NewContext()
		arg := Argument("test")
		err := arg.call(args, ctx)

		assert.Error(t, err)
		assert.Equal(t, ErrUnexpectedEndCommand, err)
	})

	t.Run("MatchArgument", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"value"})
		ctx := NewContext()
		arg := Argument("test")
		err := arg.call(args, ctx)

		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"test": "value"}, ctx.arguments)
	})

	t.Run("WithOptions", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test", "--opt", "value", "-b", "value2"})
		ctx := NewContext()
		b := 'b'
		opt1 := &option{long: "aa"}
		opt2 := &option{long: "opt"}
		opt3 := &option{long: "bb", short: &b}
		cmd := Argument("test", opt1, opt2, opt3)
		err := cmd.call(args, ctx)

		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"opt": "value", "bb": "value2"}, ctx.options)
	})

	t.Run("WithHandler", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test"})
		ctx := NewContext()
		handlerCalled := false
		cmd := Argument("test", &options.Handler{
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
		cmd := Argument("test", &options.Handler{
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

	t.Run("WithValidation", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"v1"})
		ctx := &Context{arguments: make(map[string]string)}
		opt := Argument("test", Validate(regexp.MustCompile("^v[0-9]+$")))

		err := opt.call(args, ctx)
		assert.NoError(t, err)
		assert.Equal(t, "v1", ctx.arguments["test"])
	})

	t.Run("WithInvalidValidation", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"invalid"})
		ctx := &Context{arguments: make(map[string]string)}
		opt := Argument("test", Validate(regexp.MustCompile("^v[0-9]+$")))

		err := opt.call(args, ctx)
		assert.Error(t, err)
		assert.Equal(t, &InvalidValueError{
			on:    "test",
			value: "invalid",
		}, err)
	})

	t.Run("WithOptionFailingValidation", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"test", "--opt", "value", "-b", "value2"})
		ctx := NewContext()
		b := 'b'
		opt1 := &option{long: "aa"}
		opt2 := &option{long: "opt", validate: regexp.MustCompile("^v[0-9]+$")}
		opt3 := &option{long: "bb", short: &b}
		cmd := Argument("test", opt1, opt2, opt3)
		err := cmd.call(args, ctx)

		assert.Error(t, err)
		assert.Equal(t, &InvalidValueError{
			on:    "opt",
			value: "value",
		}, err)
	})

	t.Run("WithHelpOption", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"--help"})
		ctx := NewContext()
		opt := Argument("help")
		err := opt.call(args, ctx)
		assert.Error(t, err)
		assert.IsType(t, &HelpError{}, err)
	})

	t.Run("WithHelpOptionAndOption", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"-t", "--help"})
		ctx := NewContext()
		opt := Argument("help", &option{long: "test"})
		err := opt.call(args, ctx)
		assert.Error(t, err)
		assert.IsType(t, &HelpError{}, err)
	})

	t.Run("SubCommandExecution", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"a", "sub"})
		ctx := NewContext()
		subCmd := Command("sub", &options.Handler{
			Handler: HandlerFunc(func(ctx *Context) error {
				ctx.arguments["sub"] = "executed"
				return nil
			}),
		})
		arg := Argument("test", subCmd)

		err := arg.call(args, ctx)

		assert.NoError(t, err)
		assert.Equal(t, "executed", ctx.arguments["sub"])
	})

	t.Run("UnknownSubCommand", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"a", "unknown"})
		ctx := NewContext()
		subCmd := Command("sub")
		arg := Argument("test", subCmd)

		err := arg.call(args, ctx)

		assert.Error(t, err)
		assert.IsType(t, UnknownCommandError(""), err)
	})

	t.Run("SubCommandWithOptions", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"a", "sub", "--opt", "value"})
		ctx := NewContext()
		opt := &option{long: "opt"}
		subCmd := Command("sub", opt)
		arg := Argument("test", subCmd)

		err := arg.call(args, ctx)

		assert.NoError(t, err)
		assert.Equal(t, "value", ctx.options["opt"])
	})

	t.Run("SubCommandWithHelpOption", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"a", "sub", "--help"})
		ctx := NewContext()
		subCmd := Command("sub")
		arg := Argument("test", subCmd)

		err := arg.call(args, ctx)

		assert.Error(t, err)
		assert.IsType(t, &HelpError{}, err)
	})
}
