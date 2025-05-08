package options

import (
	"github.com/StevenCyb/GoCLI/pkg/cli/internal/restriction"
)

type Handler struct {
	restriction.IsCliOption
	restriction.IsCommandOption
	restriction.IsArgumentOption

	Handler interface{}
}
