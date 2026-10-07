package service_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/adapter/filesystem"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/repository/jsonfile"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/service"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		wantErr  string
		wantFile string
	}{
		{
			name:     "should name chapter file by chapter id when translation is complete",
			lang:     "en",
			wantFile: "01-programming-conventions.md",
		},
		{
			name:    "should list clause id when translation misses a clause",
			lang:    "fr",
			wantErr: "NAMING-001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := t.TempDir()
			tree := os.DirFS(writeTestData(t))
			gen := service.NewGenerator(jsonfile.New(tree), jsonfile.New(tree), filesystem.NewWriter())

			paths, err := gen.Generate(t.Context(), out, []string{tt.lang})
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Contains(t, paths, filepath.Join(out, tt.lang, tt.wantFile))
		})
	}
}

// writeTestData builds a minimal data directory with Chinese base data and en/fr overlays, returning the data directory.
func writeTestData(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	writeJSON(t, filepath.Join(dir, "manual.json"), domain.Manual{
		Version:    "3.0",
		GoBaseline: "1.27.1",
		ToC: []domain.Chapter{{
			ID:       "programming-conventions",
			Name:     "编程规约",
			Sections: []domain.Section{{ID: "naming", Name: "命名规约"}},
		}},
	})
	writeJSON(t, filepath.Join(dir, "clauses", "programming-conventions.json"), []domain.Clause{
		{ID: "NAMING-001", Level: domain.LevelMust, Chapter: "programming-conventions", Section: "naming", SinceGo: "1.0", Summary: "包名使用小写单词连写。"},
	})
	// The en overlay is complete; fr translates the chapter but leaves the section and the clause untranslated, to verify the missing list.
	writeJSON(t, filepath.Join(dir, "i18n", "en", "manual.json"), domain.I18n{
		ToC: map[string]domain.I18nChapter{
			"programming-conventions": {Name: "Programming Conventions", Sections: map[string]string{"naming": "Naming"}},
		},
		Appendix: domain.I18nAppendix{Slug: "appendix"},
	})
	writeJSON(t, filepath.Join(dir, "i18n", "en", "clauses", "programming-conventions.json"),
		map[string]domain.I18nClause{"NAMING-001": {Summary: "s", Details: "d", Rationale: "r"}})
	writeJSON(t, filepath.Join(dir, "i18n", "fr", "manual.json"), domain.I18n{
		ToC: map[string]domain.I18nChapter{
			"programming-conventions": {Name: "Conventions de programmation"},
		},
	})
	return dir
}

func writeJSON(t *testing.T, path string, v any) {
	t.Helper()
	raw, err := json.Marshal(v)
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o750))
	require.NoError(t, os.WriteFile(path, raw, 0o600))
}
