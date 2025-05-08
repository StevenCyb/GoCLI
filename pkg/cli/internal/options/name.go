package options

import "github.com/StevenCyb/GoCLI/pkg/cli/internal/restriction"

type Name struct {
	restriction.IsCliOption

	Name string
}
