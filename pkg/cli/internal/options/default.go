package options

import "github.com/StevenCyb/GoCLI/pkg/cli/internal/restriction"

type Default struct {
	restriction.IsOptionOption

	DefaultValue string
}
