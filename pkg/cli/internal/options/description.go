package options

import "github.com/StevenCyb/GoCLI/pkg/cli/internal/restriction"

type Description struct {
	restriction.IsCliOption
	restriction.IsOptionOption
	restriction.IsCommandOption
	restriction.IsArgumentOption

	Description string
}
