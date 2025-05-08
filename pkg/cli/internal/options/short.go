package options

import "github.com/StevenCyb/GoCLI/pkg/cli/internal/restriction"

type Short struct {
	restriction.IsOptionOption

	Short rune
}
