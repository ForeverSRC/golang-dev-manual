// Package jsonfile implements the clause data reading ports, loading manual data from JSON files.
package jsonfile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"path"
	"strings"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/service"
)

// Repository reads the manual metadata, the per-chapter clause data, and the language overlays from a JSON data tree.
type Repository struct {
	tree fs.FS
}

var (
	_ service.ClauseRepository = (*Repository)(nil)
	_ service.I18nRepository   = (*Repository)(nil)
)

// New creates the JSON file clause repository over the given data tree.
func New(tree fs.FS) *Repository { return &Repository{tree: tree} }

// Load reads the manual metadata, then concatenates each chapter's clause file in toc order,
// and checks that the clauses' category ids resolve.
func (r *Repository) Load(_ context.Context) (*domain.Manual, error) {
	manual := new(domain.Manual)
	if err := readJSON(r.tree, "manual.json", manual); err != nil {
		return nil, err
	}

	for _, ch := range manual.ToC {
		var clauses []domain.Clause
		if err := readJSON(r.tree, path.Join("clauses", ch.ID+".json"), &clauses); err != nil {
			return nil, err
		}
		manual.Clauses = append(manual.Clauses, clauses...)
	}
	if err := validate(manual); err != nil {
		return nil, err
	}
	return manual, nil
}

// validate checks the data's internal references: ids are unique inside the table of contents and across the clauses,
// every clause points at a chapter and a section that exist, and the clause tags stay within the vocabulary.
func validate(manual *domain.Manual) error {
	sections := make(map[string]map[string]bool, len(manual.ToC))
	for _, ch := range manual.ToC {
		if _, ok := sections[ch.ID]; ok {
			return fmt.Errorf("duplicate chapter id %s", ch.ID)
		}
		own := make(map[string]bool, len(ch.Sections))
		for _, s := range ch.Sections {
			if own[s.ID] {
				return fmt.Errorf("duplicate section id %s in chapter %s", s.ID, ch.ID)
			}
			own[s.ID] = true
		}
		sections[ch.ID] = own
	}

	vocabulary := make(map[string]bool, len(manual.Tags))
	for _, t := range manual.Tags {
		vocabulary[t.ID] = true
	}
	used := make(map[string]bool, len(vocabulary))
	clauseIDs := make(map[string]bool, len(manual.Clauses))

	for _, c := range manual.Clauses {
		if clauseIDs[c.ID] {
			return fmt.Errorf("duplicate clause id %s", c.ID)
		}
		clauseIDs[c.ID] = true

		own, ok := sections[c.Chapter]
		if !ok {
			return fmt.Errorf("clause %s references unknown chapter %s", c.ID, c.Chapter)
		}
		if c.Section != "" && !own[c.Section] {
			return fmt.Errorf("clause %s references unknown section %s in chapter %s", c.ID, c.Section, c.Chapter)
		}
		seen := make(map[string]bool, len(c.Tags))
		for _, t := range c.Tags {
			if !vocabulary[t] {
				return fmt.Errorf("clause %s references unknown tag %s", c.ID, t)
			}
			if seen[t] {
				return fmt.Errorf("duplicate tag %s in clause %s", t, c.ID)
			}
			seen[t] = true
			used[t] = true
		}
	}
	for _, t := range manual.Tags {
		if !used[t.ID] {
			return fmt.Errorf("tag %s is not used by any clause", t.ID)
		}
	}
	return nil
}

// LoadI18n reads a language's overlay: the copy and table of contents from i18n/<lang>/manual.json, plus the clause translations under i18n/<lang>/clauses/.
// When the clauses directory is absent, the language is treated as having no clause translations, and the missing-translation check is left to the caller.
func (r *Repository) LoadI18n(_ context.Context, lang string) (*domain.I18n, error) {
	overlayDir := path.Join("i18n", lang)
	overlay := new(domain.I18n)
	if err := readJSON(r.tree, path.Join(overlayDir, "manual.json"), overlay); err != nil {
		return nil, err
	}

	overlay.Clauses = make(map[string]domain.I18nClause)
	clausesDir := path.Join(overlayDir, "clauses")
	entries, err := fs.ReadDir(r.tree, clausesDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return overlay, nil
		}
		return nil, fmt.Errorf("read %s: %w", clausesDir, err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		part := make(map[string]domain.I18nClause)
		if err := readJSON(r.tree, path.Join(clausesDir, e.Name()), &part); err != nil {
			return nil, err
		}
		maps.Copy(overlay.Clauses, part)
	}
	return overlay, nil
}

func readJSON(fsys fs.FS, name string, v any) error {
	raw, err := fs.ReadFile(fsys, name)
	if err != nil {
		return fmt.Errorf("read %s: %w", name, err)
	}
	if err := json.Unmarshal(raw, v); err != nil {
		return fmt.Errorf("parse %s: %w", name, err)
	}
	return nil
}
