package cli

import (
	"regexp"
	"testing"

	"GoCLI/pkg/cli/internal/options"
	"GoCLI/pkg/cli/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestOption(t *testing.T) {
	t.Parallel()

	t.Run("LongNameOnly", func(t *testing.T) {
		t.Parallel()

		opt := Option("test")
		assert.Equal(t, "test", opt.long)
		assert.Nil(t, opt.short)
		assert.False(t, opt.required)
		assert.Nil(t, opt.defaultValue)
	})

	t.Run("ShortName", func(t *testing.T) {
		t.Parallel()

		short := 't'
		opt := Option("test", &options.Short{Short: short})
		assert.Equal(t, "test", opt.long)
		assert.NotNil(t, opt.short)
		assert.Equal(t, short, *opt.short)
	})

	t.Run("Required", func(t *testing.T) {
		t.Parallel()

		opt := Option("test", &options.Required{})
		assert.True(t, opt.required)
	})

	t.Run("DefaultValue", func(t *testing.T) {
		t.Parallel()

		defaultValue := "default"
		opt := Option("test", &options.Default{DefaultValue: defaultValue})
		assert.NotNil(t, opt.defaultValue)
		assert.Equal(t, defaultValue, *opt.defaultValue)
	})
}

func TestOptionCall(t *testing.T) {
	t.Parallel()

	t.Run("MatchLongOption", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"--test", "value"})
		ctx := &Context{options: make(map[string]string)}
		opt := Option("test")

		err := opt.call(args, ctx)
		assert.NoError(t, err)
		assert.Equal(t, "value", ctx.options["test"])
	})

	t.Run("MatchShortOption", func(t *testing.T) {
		t.Parallel()

		short := 't'
		args := utils.NewAdvancedArray([]string{"-t", "value"})
		ctx := &Context{options: make(map[string]string)}
		opt := Option("test", &options.Short{Short: short})

		err := opt.call(args, ctx)
		assert.NoError(t, err)
		assert.Equal(t, "value", ctx.options["test"])
	})

	t.Run("EmptyValueIfNone", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"--test"})
		ctx := &Context{options: make(map[string]string)}
		opt := Option("test")

		err := opt.call(args, ctx)
		assert.NoError(t, err)
		assert.Equal(t, "", ctx.options["test"])
	})

	t.Run("MatchSecondOption", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"--test", "--test2"})
		ctx := &Context{options: make(map[string]string)}
		opt := Option("test")

		err := opt.call(args, ctx)
		assert.NoError(t, err)
		assert.Equal(t, "", ctx.options["test"])
	})

	t.Run("ReturnErrNotMatched", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"--other"})
		ctx := &Context{options: make(map[string]string)}
		opt := Option("test")

		err := opt.call(args, ctx)
		assert.ErrorIs(t, err, ErrNotMatched)
	})

	t.Run("WithValidation", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"--test", "v1"})
		ctx := &Context{options: make(map[string]string)}
		opt := Option("test", Validate(regexp.MustCompile("^v[0-9]+$")))

		err := opt.call(args, ctx)
		assert.NoError(t, err)
		assert.Equal(t, "v1", ctx.options["test"])
	})

	t.Run("WithInvalidValidation", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"--test", "invalid"})
		ctx := &Context{options: make(map[string]string)}
		opt := Option("test", Validate(regexp.MustCompile("^v[0-9]+$")))

		err := opt.call(args, ctx)
		assert.Error(t, err)
		assert.Equal(t, &InvalidValueError{
			on:    "test",
			value: "invalid",
		}, err)
	})

	t.Run("WithHelpOption", func(t *testing.T) {
		t.Parallel()

		args := utils.NewAdvancedArray([]string{"--help"})
		ctx := &Context{options: make(map[string]string)}
		opt := Option("help")
		err := opt.call(args, ctx)
		assert.Error(t, err)
		assert.IsType(t, &HelpError{}, err)
	})
}
