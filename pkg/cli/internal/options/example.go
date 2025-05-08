package options

import "github.com/StevenCyb/GoCLI/pkg/cli/internal/restriction"

type Example struct {
	restriction.IsCliOption
	restriction.IsCommandOption
	restriction.IsArgumentOption

	Example string
}
