package clihandler

import (
	"encoding/json"
	"io"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/service"
)

// clauseJSON is one clause as list prints it: the identifying fields, the display category, and the one-line summary.
type clauseJSON struct {
	ID       string   `json:"id"`
	Level    string   `json:"level"`
	Chapter  string   `json:"chapter"`
	Section  string   `json:"section,omitempty"`
	Category string   `json:"category"`
	Summary  string   `json:"summary"`
	Tags     []string `json:"tags,omitempty"`
}

// searchJSON adds the relevance score and the matched terms to the clause fields.
type searchJSON struct {
	clauseJSON
	Score   int            `json:"score"`
	Matched []domain.Match `json:"matched"`
}

// checkJSON is the machine-readable check report. Uncovered is always present, so consuming it does not depend on --verbose.
type checkJSON struct {
	Hits      []hitJSON       `json:"hits"`
	Uncovered []uncoveredJSON `json:"uncovered"`
}

// hitJSON is one grep hit.
type hitJSON struct {
	Path     string `json:"path"`
	Line     int    `json:"line"`
	Text     string `json:"text"`
	ClauseID string `json:"clause_id"`
	Level    string `json:"level"`
	Summary  string `json:"summary"`
}

// uncoveredJSON is one clause without an automated check, with its check method already rendered.
type uncoveredJSON struct {
	ID     string `json:"id"`
	Detect string `json:"detect"`
}

// clauseJSONOf builds the clause fields list and search output share.
func clauseJSONOf(c domain.Clause, category string) clauseJSON {
	return clauseJSON{
		ID:       c.ID,
		Level:    string(c.Level),
		Chapter:  c.Chapter,
		Section:  c.Section,
		Category: category,
		Summary:  c.Summary,
		Tags:     c.Tags,
	}
}

// clauseJSONs converts list results.
func clauseJSONs(views []service.ClauseView) []clauseJSON {
	out := make([]clauseJSON, 0, len(views))
	for _, v := range views {
		out = append(out, clauseJSONOf(v.Clause, v.Category))
	}
	return out
}

// searchJSONs converts search results.
func searchJSONs(results []service.SearchResult) []searchJSON {
	out := make([]searchJSON, 0, len(results))
	for _, r := range results {
		out = append(out, searchJSON{
			clauseJSON: clauseJSONOf(r.Clause, r.Category),
			Score:      r.Score,
			Matched:    r.Matched,
		})
	}
	return out
}

// checkReportJSON converts a check report.
func checkReportJSON(report *service.CheckReport) checkJSON {
	out := checkJSON{
		Hits:      make([]hitJSON, 0, len(report.Hits)),
		Uncovered: make([]uncoveredJSON, 0, len(report.Uncovered)),
	}
	for _, h := range report.Hits {
		out.Hits = append(out.Hits, hitJSON{
			Path:     h.Path,
			Line:     h.Line,
			Text:     h.Text,
			ClauseID: h.Rule.ClauseID,
			Level:    string(h.Rule.Level),
			Summary:  h.Rule.Summary,
		})
	}
	for _, v := range report.Uncovered {
		out.Uncovered = append(out.Uncovered, uncoveredJSON{ID: v.Clause.ID, Detect: v.Detect})
	}
	return out
}

// writeJSON encodes v into out as one JSON document.
func writeJSON(out io.Writer, v any) error {
	return json.NewEncoder(out).Encode(v)
}
