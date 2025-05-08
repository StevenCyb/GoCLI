package options

import "github.com/StevenCyb/GoCLI/pkg/cli/internal/restriction"

type Version struct {
	restriction.IsCliOption

	Version string
}
