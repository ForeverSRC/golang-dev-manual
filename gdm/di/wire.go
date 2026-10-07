// Package di centralizes the dependency-injection ProviderSet.
// Port-to-implementation bindings and each layer's constructors are gathered here, and wire generates the wiring code at compile time.
package di

import (
	"io/fs"

	"github.com/google/wire"

	"github.com/ForeverSRC/golang-dev-manual/gdm/data"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/adapter/filesystem"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/api/clihandler"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/repository/jsonfile"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/server"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/service"
)

// ProviderSet provides the implementations of each layer and binds the ports.
var ProviderSet = wire.NewSet(
	wire.InterfaceValue(new(fs.FS), data.FS),

	jsonfile.New,
	wire.Bind(new(service.ClauseRepository), new(*jsonfile.Repository)),
	wire.Bind(new(service.I18nRepository), new(*jsonfile.Repository)),

	filesystem.NewScanner,
	wire.Bind(new(service.GrepScanner), new(*filesystem.Scanner)),
	filesystem.NewWriter,
	wire.Bind(new(service.ManualWriter), new(*filesystem.Writer)),

	service.NewManualService,
	service.NewGenerator,

	clihandler.New,

	server.ProvideCLICommand,
	server.ProvideGeneratorCommand,
)
