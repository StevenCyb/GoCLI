package options

import "GoCLI/pkg/cli/internal/restriction"

type Short struct {
	restriction.IsOptionOption

	Short rune
}
