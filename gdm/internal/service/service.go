// Package service implements the application service layer, orchestrating the domain models and external capabilities.
// The ports (ClauseRepository / I18nRepository / GrepScanner / ManualWriter) are defined here by the consumer,
// and the unexported manualService and generatorService structs implement ManualService and Generator.
package service

import (
	"context"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
)

// ManualService is the application service for manual search and checking.
type ManualService interface {
	// Labels returns the CLI's field labels for the requested language.
	Labels(ctx context.Context, lang string) (domain.I18nCli, error)
	// List filters clauses by level, category, and tags, returning clause text in the requested language.
	List(ctx context.Context, lang string, levels []string, category string, tags []string) ([]ClauseView, error)
	// Tags returns the tag vocabulary with each tag's description in the requested language.
	Tags(ctx context.Context, lang string) ([]domain.Tag, error)
	// Search scores the clauses against a query in the requested language, returning the hits in relevance order.
	Search(ctx context.Context, lang, query string, levels []string, category string, tags []string) ([]SearchResult, error)
	// Explain renders the clause with the given id in the requested language.
	Explain(ctx context.Context, lang, id string) (ClauseView, error)
	// Check runs the grep checks from the clauses against the Go files under root.
	Check(ctx context.Context, lang, root string) (*CheckReport, error)
}

// ClauseView is one clause prepared for display: its text in the requested language, plus the display category and check method.
type ClauseView struct {
	Clause   domain.Clause
	Category string
	Detect   string
}

// Generator produces the bilingual markdown manual from the clause data.
type Generator interface {
	// Generate writes the manual for each language under outDir/<language> and returns the written file paths.
	Generate(ctx context.Context, outDir string, langs []string) ([]string, error)
}
