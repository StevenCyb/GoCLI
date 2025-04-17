package cli

import (
	"io"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	t.Parallel()

	name := "test-cli"
	result := Name(name)

	assert.NotNil(t, result)
	assert.Equal(t, name, result.Name)
}

func TestVersion(t *testing.T) {
	t.Parallel()

	version := "1.0.0"
	result := Version(version)

	assert.NotNil(t, result)
	assert.Equal(t, version, result.Version)
}

func TestStream(t *testing.T) {
	t.Parallel()

	stdout := io.Discard
	stderr := io.Discard
	result := Stream(stdout, stderr)

	assert.NotNil(t, result)
	assert.Equal(t, stdout, result.Stdout)
	assert.Equal(t, stderr, result.Stderr)
}

func TestBanner(t *testing.T) {
	t.Parallel()

	banner := "Welcome to CLI"
	result := Banner(banner)

	assert.NotNil(t, result)
	assert.Equal(t, banner, result.Banner)
}

func TestExample(t *testing.T) {
	t.Parallel()

	example := "cli command --help"
	result := Example(example)

	assert.NotNil(t, result)
	assert.Equal(t, example, result.Example)
}

func TestDescription(t *testing.T) {
	t.Parallel()

	description := "This is a test CLI"
	result := Description(description)

	assert.NotNil(t, result)
	assert.Equal(t, description, result.Description)
}

func TestHandler(t *testing.T) {
	t.Parallel()

	handlerFunc := func(ctx *Context) error { return nil }
	result := Handler(handlerFunc)

	assert.NotNil(t, result)
	assert.NotNil(t, result.Handler)
}

func TestRequired(t *testing.T) {
	t.Parallel()

	result := Required()

	assert.NotNil(t, result)
}

func TestShort(t *testing.T) {
	t.Parallel()

	short := 'h'
	result := Short(short)

	assert.NotNil(t, result)
	assert.Equal(t, short, result.Short)
}

func TestDefault(t *testing.T) {
	t.Parallel()

	defaultValue := "default"
	result := Default(defaultValue)

	assert.NotNil(t, result)
	assert.Equal(t, defaultValue, result.DefaultValue)
}

func TestValidate(t *testing.T) {
	t.Parallel()

	regex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	result := Validate(regex)

	assert.NotNil(t, result)
	assert.Equal(t, regex, result.Validate)
}
