package options

import "GoCLI/pkg/cli/internal/restriction"

type Default struct {
	restriction.IsOptionOption

	DefaultValue string
}
