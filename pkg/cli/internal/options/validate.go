package options

import (
	"regexp"

	"github.com/StevenCyb/GoCLI/pkg/cli/internal/restriction"
)

type Validate struct {
	restriction.IsOptionOption
	restriction.IsArgumentOption

	Validate *regexp.Regexp
}
