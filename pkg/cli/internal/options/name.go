package options

import "GoCLI/pkg/cli/internal/restriction"

type Name struct {
	restriction.IsCliOption

	Name string
}
