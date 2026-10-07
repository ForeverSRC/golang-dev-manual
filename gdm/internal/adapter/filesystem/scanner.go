// Package filesystem implements two outbound adapters: source scanning and manual writing.
package filesystem

import (
	"bufio"
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/service"
)

// Scanner walks the .go files under a directory and applies grep rules.
type Scanner struct{}

var _ service.GrepScanner = (*Scanner)(nil)

// NewScanner creates the source scanning adapter.
func NewScanner() *Scanner { return &Scanner{} }

// Scan applies each rule in turn to the .go files under root and returns all matching lines.
func (s *Scanner) Scan(_ context.Context, root string, rules []service.GrepRule) ([]service.GrepHit, error) {
	var hits []service.GrepHit
	for _, r := range rules {
		ruleHits, err := scanGoFiles(root, r)
		if err != nil {
			return nil, err
		}
		hits = append(hits, ruleHits...)
	}
	return hits, nil
}

// scanGoFiles walks the .go files under root and returns the lines matching the rule.
func scanGoFiles(root string, rule service.GrepRule) ([]service.GrepHit, error) {
	var found []service.GrepHit
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if name := d.Name(); name == "vendor" || name == ".git" {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		return grepFile(path, rule, &found)
	})
	return found, err
}

func grepFile(path string, rule service.GrepRule, found *[]service.GrepHit) error {
	f, err := os.Open(path) //nolint:gosec // the path is given explicitly by the caller
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	line := 0
	for scanner.Scan() {
		line++
		text := scanner.Text()
		if rule.Pattern.MatchString(text) {
			*found = append(*found, service.GrepHit{
				Path: path,
				Line: line,
				Text: strings.TrimSpace(text),
				Rule: rule,
			})
		}
	}
	return scanner.Err()
}
