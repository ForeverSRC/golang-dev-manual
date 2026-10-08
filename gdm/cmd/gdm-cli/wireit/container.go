package wireit

import (
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/api/clihandler"
)

// CLIITTestContainer is gdm-cli's integration-test container, holding the handler each case builds its command tree from.
type CLIITTestContainer struct {
	Handler *clihandler.CommandLineHandler
}
