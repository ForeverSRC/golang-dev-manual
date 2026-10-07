package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/service"
)

// Writer writes the manual files into the target directory.
type Writer struct{}

var _ service.ManualWriter = (*Writer)(nil)

// NewWriter creates the manual writing adapter.
func NewWriter() *Writer { return &Writer{} }

// Write writes the "filename -> content" map into dir and returns the full paths sorted by name.
func (w *Writer) Write(_ context.Context, dir string, files map[string]string) ([]string, error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, err
	}

	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	paths := make([]string, 0, len(names))
	for _, name := range names {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(files[name]), 0o600); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}
