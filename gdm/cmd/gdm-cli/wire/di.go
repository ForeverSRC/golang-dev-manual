//go:build wireinject

package wire

import (
	"github.com/google/wire"

	"github.com/ForeverSRC/golang-dev-manual/gdm/di"
)

// InitializeCLI injects gdm's dependencies with wire.
func InitializeCLI() (*CLIContainer, error) {
	wire.Build(
		di.ProviderSet,
		wire.Struct(new(CLIContainer), "*"),
	)
	return nil, nil
}
