package options

import "GoCLI/pkg/cli/internal/restriction"

type Version struct {
	restriction.IsCliOption

	Version string
}
