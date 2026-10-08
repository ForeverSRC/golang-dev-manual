// Package ittest provides the fixture data integration tests run against, so their expectations stay independent of the shipped data.
package ittest

import (
	"embed"
	"io/fs"
)

//go:embed testdata
var fixture embed.FS

// ITProvideDataTree returns the fixture data tree: manual.json, clauses/, and i18n/.
func ITProvideDataTree() (fs.FS, error) {
	return fs.Sub(fixture, "testdata")
}
