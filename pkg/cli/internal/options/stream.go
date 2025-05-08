package options

import (
	"io"

	"github.com/StevenCyb/GoCLI/pkg/cli/internal/restriction"
)

type Stream struct {
	restriction.IsCliOption

	Stdout io.Writer
	Stderr io.Writer
}
