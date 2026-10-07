package service

import (
	"context"
	"fmt"
	"regexp"
	"slices"
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

func (s *manualService) List(ctx context.Context, lang string, levels []string, category string) ([]domain.Clause, error) {
	manual, overlay, err := s.load(ctx, lang)
	if err != nil {
		return nil, err
	}
	if category != "" && !slices.Contains(manual.CategoryIDs(), category) {
		return nil, fmt.Errorf("unknown category %s; available: %s", category, strings.Join(manual.CategoryIDs(), ", "))
	}

	clauses := manual.Filter(levels, category)
	resolved := make([]domain.Clause, 0, len(clauses))
	for _, c := range clauses {
		resolved = append(resolved, resolveClauseText(overlay, c))
	}
	return resolved, nil
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
	_, _, _, note := resolveClause(overlay, c)
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
