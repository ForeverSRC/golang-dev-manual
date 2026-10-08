//go:build wireinject

package wireit

import (
	"github.com/google/wire"

	"github.com/ForeverSRC/golang-dev-manual/gdm/di"
)

// InitializeCLIITTestContainer injects gdm-cli's integration-test dependencies with wire over the fixture data tree.
func InitializeCLIITTestContainer() (*CLIITTestContainer, error) {
	wire.Build(
		di.ITSet,
		di.ProviderSet,
		wire.Struct(new(CLIITTestContainer), "*"),
	)
	return nil, nil
}
