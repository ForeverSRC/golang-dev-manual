//go:build wireinject

package wire

import (
	"github.com/google/wire"

	"github.com/ForeverSRC/golang-dev-manual/gdm/di"
)

// InitializeGenerator injects gdm-gen's dependencies with wire.
func InitializeGenerator() (*GeneratorContainer, error) {
	wire.Build(
		di.ProviderSet,
		wire.Struct(new(GeneratorContainer), "*"),
	)
	return nil, nil
}
