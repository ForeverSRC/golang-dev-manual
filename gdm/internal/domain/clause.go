// Package domain defines the manual's domain models and domain methods, with no IO.
package domain

import (
	"slices"
	"strings"
)

// Level is a clause's requirement level.
type Level string

const (
	LevelMust   Level = "MUST"
	LevelShould Level = "SHOULD"
	LevelMay    Level = "MAY"
)

// Rank returns the level's sort weight within a manual section: MUST first, SHOULD next, MAY last.
func (l Level) Rank() int {
	switch l {
	case LevelMust:
		return 0
	case LevelShould:
		return 1
	case LevelMay:
		return 2
	default:
		return 2
	}
}

// Section is a chapter's second-level entry. ID is a language-independent ASCII identifier.
type Section struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Tag is one cross-cutting theme in the manual's controlled vocabulary. ID is a language-independent ASCII identifier.
type Tag struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// Chapter is one part of the manual table of contents and its sections. ID is a language-independent ASCII identifier used to derive file names.
type Chapter struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Sections []Section `json:"sections,omitempty"`
}

// Detect describes how a clause is checked.
type Detect struct {
	Tool    string `json:"tool"`              // grep / golangci-lint / manual
	Rule    string `json:"rule,omitempty"`    // golangci-lint rule name
	Pattern string `json:"pattern,omitempty"` // grep regex
	Note    string `json:"note,omitempty"`
}

// Examples holds a clause's good and bad examples.
type Examples struct {
	Good string `json:"good,omitempty"`
	Bad  string `json:"bad,omitempty"`
	Lang string `json:"lang,omitempty"`
}

// Language returns the language tag for the example code block, defaulting to Go when unspecified.
func (e Examples) Language() string {
	if e.Lang == "" {
		return "go"
	}
	return e.Lang
}

// Quote is an excerpt of the upstream text backing a clause, along with its source.
type Quote struct {
	Text   string `json:"text"`
	Source string `json:"source"`
}

// Clause is a third-level clause. Chapter and Section are ids referencing the table of contents.
type Clause struct {
	ID        string   `json:"id"`
	Level     Level    `json:"level"`
	Chapter   string   `json:"chapter"`
	Section   string   `json:"section,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	SinceGo   string   `json:"since_go"`
	Summary   string   `json:"summary"`
	Details   string   `json:"details"`
	Rationale string   `json:"rationale"`
	Quote     Quote    `json:"quote"`
	Examples  Examples `json:"examples"`
	Sources   []string `json:"sources"`
	Detect    *Detect  `json:"detect,omitempty"`
}

// Manual is the manual's single source of truth.
type Manual struct {
	Version    string    `json:"version"`
	GoBaseline string    `json:"go_baseline"`
	Tags       []Tag     `json:"tags,omitempty"`
	ToC        []Chapter `json:"toc"`
	Clauses    []Clause  `json:"clauses"`
}

// BaseLanguage is the manual's source-of-truth language; other languages hang off it as overlays.
const BaseLanguage = "zh"

// Number styles, determining how chapters and sections are numbered.
const (
	NumberStyleChinese = "chinese"
	NumberStyleArabic  = "arabic"
)

// I18nClause is a clause's translation. Chinese is the source, so only non-source languages need to provide it.
type I18nClause struct {
	Summary    string `json:"summary"`
	Details    string `json:"details"`
	Rationale  string `json:"rationale"`
	DetectNote string `json:"detect_note,omitempty"`
}

// I18nChapter is a chapter's translation: its name and its sections keyed by section id.
type I18nChapter struct {
	Name     string            `json:"name"`
	Sections map[string]string `json:"sections,omitempty"`
}

// I18nManual is the copy for the index page. Fields containing {v} are templates.
type I18nManual struct {
	Title        string `json:"title"`
	Intro        string `json:"intro"`
	BaselineLine string `json:"baseline_line"`
	CountLine    string `json:"count_line"`
	ToCHeading   string `json:"toc_heading"`
}

// I18nAppendix is the copy for the appendix.
type I18nAppendix struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	IndexHeading   string   `json:"index_heading"`
	SourcesHeading string   `json:"sources_heading"`
	DetectHeading  string   `json:"detect_heading"`
	IndexColumns   []string `json:"index_columns"`
	DetectColumns  []string `json:"detect_columns"`
}

// I18nRender holds the labels, templates, and punctuation used to render the body. Fields containing {v} are templates.
type I18nRender struct {
	NumberStyle  string `json:"number_style"`
	Separator    string `json:"separator"`
	EmptySection string `json:"empty_section"`
	Why          string `json:"why"`
	Good         string `json:"good"`
	Bad          string `json:"bad"`
	Sources      string `json:"sources"`
	CategoryLine string `json:"category_line"`
	SinceGoLine  string `json:"since_go_line"`
	TagsLine     string `json:"tags_line"`
	DetectLine   string `json:"detect_line"`
	QuoteMark    string `json:"quote_mark"`
	DetectGrep   string `json:"detect_grep"`
	DetectLint   string `json:"detect_lint"`
	DetectManual string `json:"detect_manual"`
	DetectNote   string `json:"detect_note"`
}

// I18nCli holds the labels gdm-cli prints around clause fields.
type I18nCli struct {
	Details    string `json:"details"`
	Why        string `json:"why"`
	Quote      string `json:"quote"`
	Source     string `json:"source"`
	Good       string `json:"good"`
	Bad        string `json:"bad"`
	References string `json:"references"`
	Detection  string `json:"detection"`
	Category   string `json:"category"`
	SinceGo    string `json:"since_go"`
	NoHits     string `json:"no_hits"`
	Uncovered  string `json:"uncovered"`
}

// I18n is the render overlay for one language. ToC, Clauses, and Tags are only needed for non-source languages, falling back to the source language when missing.
type I18n struct {
	ToC      map[string]I18nChapter `json:"toc,omitempty"`
	Clauses  map[string]I18nClause  `json:"clauses,omitempty"`
	Tags     map[string]string      `json:"tags,omitempty"`
	Manual   I18nManual             `json:"manual"`
	Appendix I18nAppendix           `json:"appendix"`
	Render   I18nRender             `json:"render"`
	Cli      I18nCli                `json:"cli"`
}

// CategoryID returns the clause's category id: its chapter id, or "chapter/section" when it has a section.
func (c Clause) CategoryID() string {
	if c.Section == "" {
		return c.Chapter
	}
	return c.Chapter + "/" + c.Section
}

// ByID looks up a clause by id.
func (m *Manual) ByID(id string) (Clause, bool) {
	for _, c := range m.Clauses {
		if c.ID == id {
			return c, true
		}
	}
	return Clause{}, false
}

// TagIDs lists every valid tag id in the vocabulary.
func (m *Manual) TagIDs() []string {
	ids := make([]string, 0, len(m.Tags))
	for _, t := range m.Tags {
		ids = append(ids, t.ID)
	}
	return ids
}

// CategoryIDs lists every valid category id: each chapter id and, for chapters with sections, each "chapter/section" id.
func (m *Manual) CategoryIDs() []string {
	ids := make([]string, 0, len(m.ToC))
	for _, ch := range m.ToC {
		ids = append(ids, ch.ID)
		for _, s := range ch.Sections {
			ids = append(ids, ch.ID+"/"+s.ID)
		}
	}
	return ids
}

// Filter filters clauses by level, category, and tags; an empty condition means no filtering.
// category accepts either a chapter id or a two-level "chapter/section" id.
// tags keep a clause that carries any of the given tags.
func (m *Manual) Filter(levels []string, category string, tags []string) []Clause {
	wanted := make(map[string]bool, len(levels))
	for _, l := range levels {
		wanted[strings.ToUpper(strings.TrimSpace(l))] = true
	}
	wantedTags := make(map[string]bool, len(tags))
	for _, t := range tags {
		wantedTags[t] = true
	}
	var result []Clause
	for _, c := range m.Clauses {
		if len(wanted) > 0 && !wanted[string(c.Level)] {
			continue
		}
		if category != "" && c.Chapter != category && c.CategoryID() != category {
			continue
		}
		if len(wantedTags) > 0 && !slices.ContainsFunc(c.Tags, func(t string) bool { return wantedTags[t] }) {
			continue
		}
		result = append(result, c)
	}
	return result
}
