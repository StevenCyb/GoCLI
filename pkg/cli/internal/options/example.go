package options

import "GoCLI/pkg/cli/internal/restriction"

type Example struct {
	restriction.IsCliOption
	restriction.IsCommandOption
	restriction.IsArgumentOption

	Example string
}
