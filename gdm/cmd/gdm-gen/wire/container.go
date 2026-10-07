package wire

import (
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/server"
)

// GeneratorContainer is gdm-gen's dependency container, holding only the generator inbound adapter.
type GeneratorContainer struct {
	RootCmd server.GeneratorCommand
}
