package service

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
)

type manualService struct {
	repo    ClauseRepository
	i18n    I18nRepository
	scanner GrepScanner
}

// NewManualService creates the manual application service.
func NewManualService(repo ClauseRepository, i18n I18nRepository, scanner GrepScanner) ManualService {
	return &manualService{repo: repo, i18n: i18n, scanner: scanner}
}

func (s *manualService) Labels(ctx context.Context, lang string) (domain.I18nCli, error) {
	overlay, err := s.i18n.LoadI18n(ctx, lang)
	if err != nil {
		return domain.I18nCli{}, err
	}
	return overlay.Cli, nil
}

func (s *manualService) List(ctx context.Context, q ClauseQuery) ([]ClauseView, error) {
	manual, overlay, err := s.load(ctx, q.Lang)
	if err != nil {
		return nil, err
	}
	if err := validateFilters(manual, q.Category, q.Tags); err != nil {
		return nil, err
	}

	clauses := manual.Filter(q.Levels, q.Category, q.Tags)
	views := make([]ClauseView, 0, len(clauses))
	for _, c := range clauses {
		views = append(views, clauseView(manual, overlay, c))
	}
	return views, nil
}

func (s *manualService) Tags(ctx context.Context, lang string) ([]domain.Tag, error) {
	manual, overlay, err := s.load(ctx, lang)
	if err != nil {
		return nil, err
	}
	tags := make([]domain.Tag, 0, len(manual.Tags))
	for _, t := range manual.Tags {
		tags = append(tags, domain.Tag{ID: t.ID, Description: tagDescription(overlay, t)})
	}
	return tags, nil
}

func (s *manualService) Search(ctx context.Context, q ClauseQuery) ([]SearchResult, error) {
	manual, overlay, err := s.load(ctx, q.Lang)
	if err != nil {
		return nil, err
	}
	if err := validateFilters(manual, q.Category, q.Tags); err != nil {
		return nil, err
	}

	terms := domain.Tokenize(q.Query)
	results := make([]SearchResult, 0, len(manual.Clauses))
	for _, c := range manual.Filter(q.Levels, q.Category, q.Tags) {
		view := clauseView(manual, overlay, c)
		score, matched := view.Clause.MatchScore(terms)
		if score == 0 {
			continue
		}
		results = append(results, SearchResult{ClauseView: view, Score: score, Matched: matched})
	}
	sort.SliceStable(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		if results[i].Clause.Level.Rank() != results[j].Clause.Level.Rank() {
			return results[i].Clause.Level.Rank() < results[j].Clause.Level.Rank()
		}
		return results[i].Clause.ID < results[j].Clause.ID
	})
	return results, nil
}

func (s *manualService) Explain(ctx context.Context, lang, id string) (ClauseView, error) {
	manual, overlay, err := s.load(ctx, lang)
	if err != nil {
		return ClauseView{}, err
	}
	c, ok := manual.ByID(id)
	if !ok {
		return ClauseView{}, fmt.Errorf("clause %s not found", id)
	}
	return clauseView(manual, overlay, c), nil
}

func (s *manualService) Check(ctx context.Context, lang, root string) (*CheckReport, error) {
	manual, overlay, err := s.load(ctx, lang)
	if err != nil {
		return nil, err
	}
	rules, uncovered, err := buildGrepRules(manual, overlay)
	if err != nil {
		return nil, err
	}
	hits, err := s.scanner.Scan(ctx, root, rules)
	if err != nil {
		return nil, err
	}

	views := make([]ClauseView, 0, len(uncovered))
	for _, c := range uncovered {
		views = append(views, clauseView(manual, overlay, c))
	}
	return &CheckReport{Hits: hits, Uncovered: views}, nil
}

// validateFilters rejects a category or a tag the manual does not define, listing the available values.
func validateFilters(manual *domain.Manual, category string, tags []string) error {
	if category != "" && !slices.Contains(manual.CategoryIDs(), category) {
		return fmt.Errorf("unknown category %s; available: %s", category, strings.Join(manual.CategoryIDs(), ", "))
	}
	for _, t := range tags {
		if !slices.Contains(manual.TagIDs(), t) {
			return fmt.Errorf("unknown tag %s; available: %s", t, strings.Join(manual.TagIDs(), ", "))
		}
	}
	return nil
}

// load reads the manual and the requested language's overlay.
func (s *manualService) load(ctx context.Context, lang string) (*domain.Manual, *domain.I18n, error) {
	manual, err := s.repo.Load(ctx)
	if err != nil {
		return nil, nil, err
	}
	overlay, err := s.i18n.LoadI18n(ctx, lang)
	if err != nil {
		return nil, nil, err
	}
	return manual, overlay, nil
}

// clauseView prepares one clause for display in the overlay's language.
func clauseView(manual *domain.Manual, overlay *domain.I18n, c domain.Clause) ClauseView {
	note := resolveClause(overlay, c).note
	return ClauseView{
		Clause:   resolveClauseText(overlay, c),
		Category: categoryName(manual, overlay, c),
		Detect:   detectText(overlay.Render, c.Detect, note),
	}
}

// buildGrepRules compiles the grep checks in the clauses into rules and collects the clauses without an automated check.
func buildGrepRules(manual *domain.Manual, overlay *domain.I18n) ([]GrepRule, []domain.Clause, error) {
	var rules []GrepRule
	var uncovered []domain.Clause
	for _, c := range manual.Clauses {
		if c.Detect == nil || c.Detect.Tool != "grep" || c.Detect.Pattern == "" {
			uncovered = append(uncovered, c)
			continue
		}
		re, err := regexp.Compile(c.Detect.Pattern)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid grep pattern for clause %s: %w", c.ID, err)
		}
		rules = append(rules, GrepRule{ClauseID: c.ID, Level: c.Level, Summary: resolveClauseText(overlay, c).Summary, Pattern: re})
	}
	return rules, uncovered, nil
}
