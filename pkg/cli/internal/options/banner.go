package options

import "GoCLI/pkg/cli/internal/restriction"

type Banner struct {
	restriction.IsCliOption

	Banner string
}
