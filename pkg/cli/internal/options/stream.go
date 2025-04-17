package options

import (
	"GoCLI/pkg/cli/internal/restriction"
	"io"
)

type Stream struct {
	restriction.IsCliOption

	Stdout io.Writer
	Stderr io.Writer
}
