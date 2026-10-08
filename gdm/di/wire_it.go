package di

import (
	"github.com/google/wire"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/ittest"
)

// ITSet provides the fixture data tree the integration tests run against.
var ITSet = wire.NewSet(
	ittest.ITProvideDataTree,
)
