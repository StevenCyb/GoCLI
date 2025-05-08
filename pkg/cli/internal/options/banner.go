package options

import "github.com/StevenCyb/GoCLI/pkg/cli/internal/restriction"

type Banner struct {
	restriction.IsCliOption

	Banner string
}
