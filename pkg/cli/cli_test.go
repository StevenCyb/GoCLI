package cli

import (
	"errors"
	"regexp"
	"testing"

	"GoCLI/pkg/cli/internal/options"

	"github.com/stretchr/testify/assert"
)

func TestCLI(t *testing.T) {
	t.Parallel()

	t.Run("BasicCreation", func(t *testing.T) {
		t.Parallel()

		cli := New()
		assert.Nil(t, cli.name)
		assert.Nil(t, cli.version)
		assert.Nil(t, cli.banner)
		assert.Nil(t, cli.example)
		assert.Nil(t, cli.description)
		assert.NotNil(t, cli.stdout)
		assert.NotNil(t, cli.stderr)
		assert.Nil(t, cli.handler)
		assert.Empty(t, cli.command)
		assert.Nil(t, cli.argument)
		assert.Empty(t, cli.options)
	})

	t.Run("WithMeta", func(t *testing.T) {
		t.Parallel()

		banner := "test banner"
		example := "test example"
		description := "test description"
		cli := New(
			&options.Banner{Banner: banner},
			&options.Example{Example: example},
			&options.Description{Description: description},
			&options.Stream{Stdout: nil, Stderr: nil},
		)

		assert.NotNil(t, cli.banner)
		assert.Equal(t, banner, *cli.banner)
		assert.NotNil(t, cli.example)
		assert.Equal(t, example, *cli.example)
		assert.NotNil(t, cli.description)
		assert.Equal(t, description, *cli.description)
		assert.Nil(t, cli.stdout)
		assert.Nil(t, cli.stderr)
	})

	t.Run("WithNameAndVersion", func(t *testing.T) {
		t.Parallel()

		name := "test-cli"
		version := "1.0.0"
		cli := New(
			&options.Name{Name: name},
			&options.Version{Version: version},
		)

		assert.NotNil(t, cli.name)
		assert.Equal(t, name, *cli.name)
		assert.NotNil(t, cli.version)
		assert.Equal(t, version, *cli.version)
	})

	t.Run("WithHandler", func(t *testing.T) {
		t.Parallel()

		handler := func(ctx *Context) error { return nil }
		cli := New(&options.Handler{Handler: HandlerFunc(handler)})

		assert.NotNil(t, cli.handler)
		assert.NotNil(t, *cli.handler)
	})

	t.Run("WithOptions", func(t *testing.T) {
		t.Parallel()
		opt1 := &option{long: "opt1"}
		opt2 := &option{long: "opt2"}
		cli := New(opt1, opt2)
		assert.Len(t, cli.options, 2)
		assert.Equal(t, opt1, cli.options[0])
		assert.Equal(t, opt2, cli.options[1])
	})

	t.Run("WithCommands", func(t *testing.T) {
		t.Parallel()

		cmd := Command("test")
		cli := New(cmd)

		assert.Len(t, cli.command, 1)
		assert.Equal(t, "test", cli.command[0].name)
	})

	t.Run("DuplicateCommands", func(t *testing.T) {
		t.Parallel()

		cmd := Command("test")
		assert.PanicsWithError(t, "duplicate command: test", func() {
			New(cmd, cmd)
		})
	})

	t.Run("WithArgument", func(t *testing.T) {
		t.Parallel()

		arg := &argument{name: "arg"}
		cli := New(arg)

		assert.NotNil(t, cli.argument)
		assert.Equal(t, "arg", cli.argument.name)
	})

	t.Run("MixArgumentAndCommands", func(t *testing.T) {
		t.Parallel()

		arg := &argument{name: "arg"}
		cmd := Command("test")
		assert.PanicsWithError(t, "mix of argument and command: test", func() {
			New(arg, cmd)
		})
	})

	t.Run("MixArgumentAndCommands", func(t *testing.T) {
		t.Parallel()

		arg := &argument{name: "arg"}
		cmd := Command("test")
		assert.PanicsWithError(t, "mix of argument and command: arg", func() {
			New(cmd, arg)
		})
	})

	t.Run(("MixArgumentAndCommandsWithOptions"), func(t *testing.T) {
		t.Parallel()

		arg := &argument{name: "arg"}
		cmd := Command("test")
		assert.PanicsWithError(t, "mix of argument and command: test", func() {
			New(arg, cmd, &options.Name{Name: "test"})
		})
	})
}

func TestCLIRun(t *testing.T) {
	t.Parallel()

	t.Run("RunWithNoCommandsOrArguments", func(t *testing.T) {
		t.Parallel()

		cli := New()
		ctx, err := cli.RunWith([]string{"test"})

		assert.NoError(t, err)
		assert.NotNil(t, ctx)
	})

	t.Run("RunWithCommand", func(t *testing.T) {
		t.Parallel()

		cmd1 := Command("test")
		cmd2 := Command("test2")
		cli := New(cmd1, cmd2)
		ctx, err := cli.RunWith([]string{"test2"})

		assert.NoError(t, err)
		assert.NotNil(t, ctx)
		assert.Equal(t, []string{"test2"}, ctx.commands)
	})

	t.Run("RunWithUnknownCommand", func(t *testing.T) {
		t.Parallel()

		cmd := Command("test")
		cli := New(cmd)
		ctx, err := cli.RunWith([]string{"unknown"})

		assert.Error(t, err)
		assert.Nil(t, ctx)
		assert.Equal(t, UnknownCommandError("unknown"), err)
	})

	t.Run("RunWithHandler", func(t *testing.T) {
		t.Parallel()

		handlerCalled := false
		handler := func(ctx *Context) error {
			handlerCalled = true
			return nil
		}
		cli := New(&options.Handler{Handler: HandlerFunc(handler)})
		ctx, err := cli.RunWith([]string{"test"})

		assert.NoError(t, err)
		assert.NotNil(t, ctx)
		assert.True(t, handlerCalled)
	})

	t.Run("RunWithHandlerError", func(t *testing.T) {
		t.Parallel()

		handler := func(ctx *Context) error {
			return errors.New("handler error")
		}
		cli := New(&options.Handler{Handler: HandlerFunc(handler)})
		ctx, err := cli.RunWith([]string{"test"})

		assert.Error(t, err)
		assert.Nil(t, ctx)
		assert.Equal(t, "handler error", err.Error())
	})

	t.Run("RunWithArgument", func(t *testing.T) {
		t.Parallel()

		arg := &argument{name: "arg"}
		cli := New(arg)
		ctx, err := cli.RunWith([]string{"test"})

		assert.NoError(t, err)
		assert.NotNil(t, ctx)
		assert.Equal(t, map[string]string{"arg": "test"}, ctx.arguments)
	})

	t.Run("WithOptions", func(t *testing.T) {
		t.Parallel()

		opt := &option{long: "opt"}
		cli := New(opt)
		ctx, err := cli.RunWith([]string{"--opt", "value"})

		assert.NoError(t, err)
		assert.NotNil(t, ctx)
		assert.Equal(t, map[string]string{"opt": "value"}, ctx.options)
	})

	t.Run("WithOptionFailingValidation", func(t *testing.T) {
		t.Parallel()

		args := []string{"--opt", "value", "-b", "value2"}
		b := 'b'
		opt1 := &option{long: "aa"}
		opt2 := &option{long: "opt", validate: regexp.MustCompile("^v[0-9]+$")}
		opt3 := &option{long: "bb", short: &b}
		cli := New(opt1, opt2, opt3)
		ctx, err := cli.RunWith(args)

		assert.Error(t, err)
		assert.Nil(t, ctx)
		assert.Equal(t, &InvalidValueError{
			on:    "opt",
			value: "value",
		}, err)
	})

	t.Run("WithArgumentWithOptionFailingValidation", func(t *testing.T) {
		t.Parallel()

		args := []string{"x", "--opt", "value", "-b", "value2"}
		b := 'b'
		opt1 := &option{long: "aa"}
		opt2 := &option{long: "opt", validate: regexp.MustCompile("^v[0-9]+$")}
		opt3 := &option{long: "bb", short: &b}
		cli := New(Argument("test", opt1, opt2, opt3))
		ctx, err := cli.RunWith(args)

		assert.Error(t, err)
		assert.Nil(t, ctx)
		assert.Equal(t, &InvalidValueError{
			on:    "opt",
			value: "value",
		}, err)
	})
}
