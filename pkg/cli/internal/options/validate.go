package options

import (
	"GoCLI/pkg/cli/internal/restriction"
	"regexp"
)

type Validate struct {
	restriction.IsOptionOption
	restriction.IsArgumentOption

	Validate *regexp.Regexp
}
