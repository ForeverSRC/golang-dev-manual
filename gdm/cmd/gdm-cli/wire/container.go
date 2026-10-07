package wire

import (
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/server"
)

// CLIContainer is gdm's dependency container, holding only the CLI inbound adapter.
type CLIContainer struct {
	RootCmd server.CLICommand
}
