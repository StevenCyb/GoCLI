package cli

import (
	"GoCLI/pkg/cli/internal/options"
	"io"
	"regexp"
)

// CLI
func Name(name string) *options.Name {
	return &options.Name{
		Name: name,
	}
}

// CLI
func Version(version string) *options.Version {
	return &options.Version{
		Version: version,
	}
}

// CLI
func Stream(stdout, stderr io.Writer) *options.Stream {
	return &options.Stream{
		Stdout: stdout,
		Stderr: stderr,
	}
}

// CLI
func Banner(banner string) *options.Banner {
	return &options.Banner{
		Banner: banner,
	}
}

// CLI, Command, argument
func Example(example string) *options.Example {
	return &options.Example{
		Example: example,
	}
}

// CLI, Command, argument
func Description(description string) *options.Description {
	return &options.Description{
		Description: description,
	}
}

type HandlerFunc func(ctx *Context) error

// CLI, Command, argument
func Handler(fn HandlerFunc) *options.Handler {
	return &options.Handler{
		Handler: fn,
	}
}

// Option
func Required() *options.Required {
	return &options.Required{}
}

// Option
func Short(short rune) *options.Short {
	return &options.Short{
		Short: short,
	}
}

// Option
func Default(defaultValue string) *options.Default {
	return &options.Default{
		DefaultValue: defaultValue,
	}
}

// Option, argument
func Validate(reg *regexp.Regexp) *options.Validate {
	return &options.Validate{
		Validate: reg,
	}
}
