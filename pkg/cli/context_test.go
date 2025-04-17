package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext_VisitedCommand(t *testing.T) {
	t.Parallel()

	ctx := NewContext()
	ctx.commands = []string{"test"}

	assert.True(t, ctx.VisitedCommand("test"))
	assert.False(t, ctx.VisitedCommand("nonexistent"))
}

func TestContext_VisitedArgument(t *testing.T) {
	t.Parallel()

	ctx := NewContext()
	ctx.arguments["arg"] = "value"

	assert.True(t, ctx.VisitedArgument("arg"))
	assert.False(t, ctx.VisitedArgument("nonexistent"))
}

func TestContext_UsedOption(t *testing.T) {
	t.Parallel()

	ctx := NewContext()
	ctx.options["opt"] = "value"

	assert.True(t, ctx.UsedOption("opt"))
	assert.False(t, ctx.UsedOption("nonexistent"))
}

func TestContext_GetCommand(t *testing.T) {
	t.Parallel()

	ctx := NewContext()
	ctx.commands = []string{"cmd"}

	cmd := ctx.GetCommand("cmd")
	assert.NotNil(t, cmd)
	assert.Equal(t, "cmd", *cmd)

	nonexistent := ctx.GetCommand("nonexistent")
	assert.Nil(t, nonexistent)
}

func TestContext_GetArgument(t *testing.T) {
	t.Parallel()

	ctx := NewContext()
	ctx.arguments["arg"] = "value"

	arg := ctx.GetArgument("arg")
	assert.NotNil(t, arg)
	assert.Equal(t, "value", *arg)

	nonexistent := ctx.GetArgument("nonexistent")
	assert.Nil(t, nonexistent)
}

func TestContext_GetOption(t *testing.T) {
	t.Parallel()

	ctx := NewContext()
	ctx.options["opt"] = "value"

	opt := ctx.GetOption("opt")
	assert.NotNil(t, opt)
	assert.Equal(t, "value", *opt)

	nonexistent := ctx.GetOption("nonexistent")
	assert.Nil(t, nonexistent)
}
