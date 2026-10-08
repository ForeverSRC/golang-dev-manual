package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
)

var chineseNumbers = []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "十"}

// renderManual renders the manual data into a set of markdown files per the language overlay, returning "filename -> content".
// A non-source language must provide the full table of contents and clause translations, or a missing list is returned.
func renderManual(m *domain.Manual, i18n *domain.I18n, lang string) (map[string]string, error) {
	if err := checkI18n(m, i18n, lang); err != nil {
		return nil, err
	}

	style := i18n.Render.NumberStyle
	byCategory := make(map[string][]domain.Clause, len(m.Clauses))
	for _, c := range m.Clauses {
		byCategory[c.CategoryID()] = append(byCategory[c.CategoryID()], c)
	}
	for k := range byCategory {
		sort.Slice(byCategory[k], func(i, j int) bool {
			left, right := byCategory[k][i], byCategory[k][j]
			if left.Level.Rank() != right.Level.Rank() {
				return left.Level.Rank() < right.Level.Rank()
			}
			return left.ID < right.ID
		})
	}

	files := make(map[string]string, len(m.ToC)+2)
	var index strings.Builder
	fmt.Fprintf(&index, "# %s\n\n", i18n.Manual.Title)
	fmt.Fprintf(&index, "%s\n\n", i18n.Manual.Intro)
	fmt.Fprintf(&index, "- %s\n", expand(i18n.Manual.BaselineLine, m.GoBaseline))
	fmt.Fprintf(&index, "- %s\n", expand(i18n.Manual.CountLine, strconv.Itoa(len(m.Clauses))))
	fmt.Fprintf(&index, "\n## %s\n\n", i18n.Manual.ToCHeading)

	for i, ch := range m.ToC {
		filename := fmt.Sprintf("%02d-%s.md", i+1, ch.ID)
		fmt.Fprintf(&index, "%d. [%s%s](%s)\n", i+1, indexLabel(i+1, style), chapterName(m, i18n, i), filename)
		for j := range ch.Sections {
			fmt.Fprintf(&index, "    - %s\n", sectionName(m, i18n, i, j))
		}
		files[filename] = renderChapter(renderContext{manual: m, i18n: i18n, byCategory: byCategory}, i, headingNumber(i+1, style))
	}

	appendixIndex := len(m.ToC) + 1
	appendixFilename := fmt.Sprintf("%02d-%s.md", appendixIndex, i18n.Appendix.Slug)
	fmt.Fprintf(&index, "%d. [%s%s](%s)\n", appendixIndex, indexLabel(appendixIndex, style), i18n.Appendix.Title, appendixFilename)
	files[appendixFilename] = renderAppendix(m, i18n, headingNumber(appendixIndex, style))
	files["index.md"] = index.String()
	return files, nil
}

// renderContext groups the manual, overlay, and clause grouping shared while rendering chapters.
type renderContext struct {
	manual     *domain.Manual
	i18n       *domain.I18n
	byCategory map[string][]domain.Clause
}

// renderChapter renders one chapter. number is the chapter's numbering prefix.
func renderChapter(rc renderContext, ci int, number string) string {
	ch := rc.manual.ToC[ci]
	var b strings.Builder
	fmt.Fprintf(&b, "# %s%s\n\n", number, chapterName(rc.manual, rc.i18n, ci))

	if len(ch.Sections) == 0 {
		writeSection(&b, rc.manual, rc.i18n, rc.byCategory[ch.ID])
		return b.String()
	}
	for j, s := range ch.Sections {
		fmt.Fprintf(&b, "## %s%s\n\n", sectionNumber(j+1, rc.i18n.Render.NumberStyle), sectionName(rc.manual, rc.i18n, ci, j))
		writeSection(&b, rc.manual, rc.i18n, rc.byCategory[ch.ID+"/"+s.ID])
	}
	return b.String()
}

func writeSection(b *strings.Builder, m *domain.Manual, i18n *domain.I18n, clauses []domain.Clause) {
	if len(clauses) == 0 {
		fmt.Fprintf(b, "%s\n\n", i18n.Render.EmptySection)
		return
	}
	for _, c := range clauses {
		writeClause(b, m, i18n, c)
	}
}

func writeClause(b *strings.Builder, m *domain.Manual, i18n *domain.I18n, c domain.Clause) {
	t := resolveClause(i18n, c)
	fmt.Fprintf(b, "### 【%s】%s %s\n\n", c.Level, c.ID, t.summary)
	fmt.Fprintf(b, "- %s\n", expand(i18n.Render.CategoryLine, categoryName(m, i18n, c)))
	fmt.Fprintf(b, "- %s\n", expand(i18n.Render.SinceGoLine, c.SinceGo))
	if len(c.Tags) > 0 {
		fmt.Fprintf(b, "- %s\n", expand(i18n.Render.TagsLine, strings.Join(c.Tags, ", ")))
	}
	b.WriteString("\n")
	if t.details != "" {
		fmt.Fprintf(b, "%s\n\n", t.details)
	}
	if c.Quote.Text != "" || t.rationale != "" {
		fmt.Fprintf(b, "**%s**\n\n", i18n.Render.Why)
	}
	if c.Quote.Text != "" {
		b.WriteString(quoteBlock(c.Quote, i18n.Render.QuoteMark))
	}
	if t.rationale != "" {
		fmt.Fprintf(b, "%s\n\n", t.rationale)
	}
	if c.Examples.Good != "" {
		fmt.Fprintf(b, "**%s**\n\n%s\n\n", i18n.Render.Good, codeFence(c.Examples.Good, c.Examples.Language()))
	}
	if c.Examples.Bad != "" {
		fmt.Fprintf(b, "**%s**\n\n%s\n\n", i18n.Render.Bad, codeFence(c.Examples.Bad, c.Examples.Language()))
	}
	if len(c.Sources) > 0 {
		fmt.Fprintf(b, "**%s**\n\n", i18n.Render.Sources)
		for _, s := range c.Sources {
			fmt.Fprintf(b, "- %s\n", s)
		}
		b.WriteString("\n")
	}
	fmt.Fprintf(b, "%s\n\n", expand(i18n.Render.DetectLine, detectText(i18n.Render, c.Detect, t.note)))
}

func renderAppendix(m *domain.Manual, i18n *domain.I18n, number string) string {
	clauses := append([]domain.Clause(nil), m.Clauses...)
	sort.Slice(clauses, func(i, j int) bool { return clauses[i].ID < clauses[j].ID })

	var b strings.Builder
	fmt.Fprintf(&b, "# %s%s\n\n", number, i18n.Appendix.Title)

	fmt.Fprintf(&b, "## 1 %s\n\n", i18n.Appendix.IndexHeading)
	b.WriteString(tableHeader(i18n.Appendix.IndexColumns))
	for _, c := range clauses {
		summary := resolveClause(i18n, c).summary
		fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", c.ID, c.Level, categoryName(m, i18n, c), summary)
	}

	fmt.Fprintf(&b, "\n## 2 %s\n\n", i18n.Appendix.SourcesHeading)
	sep := i18n.Render.Separator
	for _, c := range clauses {
		if len(c.Sources) == 0 {
			continue
		}
		fmt.Fprintf(&b, "- %s%s%s\n", c.ID, sep, strings.Join(c.Sources, " "))
	}

	fmt.Fprintf(&b, "\n## 3 %s\n\n", i18n.Appendix.DetectHeading)
	b.WriteString(tableHeader(i18n.Appendix.DetectColumns))
	for _, c := range clauses {
		note := resolveClause(i18n, c).note
		fmt.Fprintf(&b, "| %s | %s |\n", c.ID, detectText(i18n.Render, c.Detect, note))
	}
	return b.String()
}

// clauseText is a clause's display text in one language.
type clauseText struct {
	summary   string
	details   string
	rationale string
	note      string
}

// resolveClause gets a clause's display text: use the translation if present, otherwise fall back to the source language.
func resolveClause(i18n *domain.I18n, c domain.Clause) clauseText {
	if t, ok := i18n.Clauses[c.ID]; ok {
		return clauseText{summary: t.Summary, details: t.Details, rationale: t.Rationale, note: t.DetectNote}
	}
	var note string
	if c.Detect != nil {
		note = c.Detect.Note
	}
	return clauseText{summary: c.Summary, details: c.Details, rationale: c.Rationale, note: note}
}

// resolveClauseText returns the clause with its text fields in the requested language.
func resolveClauseText(i18n *domain.I18n, c domain.Clause) domain.Clause {
	t := resolveClause(i18n, c)
	c.Summary, c.Details, c.Rationale = t.summary, t.details, t.rationale
	return c
}

// categoryName returns a clause's display category in the requested language.
func categoryName(m *domain.Manual, i18n *domain.I18n, c domain.Clause) string {
	ci := chapterIndex(m, c.Chapter)
	if ci < 0 {
		return c.CategoryID()
	}
	name := chapterName(m, i18n, ci)
	if c.Section == "" {
		return name
	}
	for si, s := range m.ToC[ci].Sections {
		if s.ID == c.Section {
			return name + "/" + sectionName(m, i18n, ci, si)
		}
	}
	return name
}

// chapterName returns a chapter's display name in the requested language.
func chapterName(m *domain.Manual, i18n *domain.I18n, ci int) string {
	if t, ok := i18n.ToC[m.ToC[ci].ID]; ok && t.Name != "" {
		return t.Name
	}
	return m.ToC[ci].Name
}

// tagDescription returns a tag's description in the requested language.
func tagDescription(i18n *domain.I18n, t domain.Tag) string {
	if d := i18n.Tags[t.ID]; d != "" {
		return d
	}
	return t.Description
}

// sectionName returns a section's display name in the requested language.
func sectionName(m *domain.Manual, i18n *domain.I18n, ci, si int) string {
	s := m.ToC[ci].Sections[si]
	if t, ok := i18n.ToC[m.ToC[ci].ID]; ok {
		if name := t.Sections[s.ID]; name != "" {
			return name
		}
	}
	return s.Name
}

// chapterIndex finds a chapter's position in the table of contents by its id.
func chapterIndex(m *domain.Manual, id string) int {
	for i, ch := range m.ToC {
		if ch.ID == id {
			return i
		}
	}
	return -1
}

// checkI18n verifies that a non-source language's overlay is complete and lists the missing entries.
func checkI18n(m *domain.Manual, i18n *domain.I18n, lang string) error {
	if lang == domain.BaseLanguage {
		return nil
	}
	var missing []string
	for _, ch := range m.ToC {
		t, ok := i18n.ToC[ch.ID]
		if !ok || t.Name == "" {
			missing = append(missing, ch.ID)
			continue
		}
		for _, s := range ch.Sections {
			if t.Sections[s.ID] == "" {
				missing = append(missing, ch.ID+"/"+s.ID)
			}
		}
	}
	for _, t := range m.Tags {
		if i18n.Tags[t.ID] == "" {
			missing = append(missing, t.ID)
		}
	}
	for _, c := range m.Clauses {
		t, ok := i18n.Clauses[c.ID]
		if !ok || t.Summary == "" || t.Details == "" || t.Rationale == "" {
			missing = append(missing, c.ID)
			continue
		}
		if c.Detect != nil && c.Detect.Note != "" && t.DetectNote == "" {
			missing = append(missing, c.ID+"(detect.note)")
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing translations for language %s: %s", lang, strings.Join(missing, ", "))
	}
	return nil
}

// detectText renders a clause's check method using the overlay labels.
func detectText(r domain.I18nRender, d *domain.Detect, note string) string {
	if d == nil {
		return "—"
	}
	var base string
	switch d.Tool {
	case "grep":
		base = expand(r.DetectGrep, d.Pattern)
	case "golangci-lint":
		base = expand(r.DetectLint, d.Rule)
	default:
		base = r.DetectManual
	}
	if note != "" {
		base += expand(r.DetectNote, note)
	}
	return base
}

func codeFence(code, lang string) string {
	return "```" + lang + "\n" + strings.TrimRight(code, "\n") + "\n```"
}

// quoteBlock renders the quote as a markdown blockquote, with the source appended after the quote.
func quoteBlock(q domain.Quote, mark string) string {
	var b strings.Builder
	for line := range strings.SplitSeq(strings.TrimRight(q.Text, "\n"), "\n") {
		fmt.Fprintf(&b, "> %s\n", line)
	}
	if q.Source != "" {
		fmt.Fprintf(&b, ">\n> %s %s\n", mark, q.Source)
	}
	b.WriteString("\n")
	return b.String()
}

// tableHeader generates the markdown header row and separator row from the column names.
func tableHeader(columns []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "| %s |\n", strings.Join(columns, " | "))
	b.WriteString("|")
	for range columns {
		b.WriteString(" --- |")
	}
	b.WriteString("\n")
	return b.String()
}

// expand fills a value into the overlay template through the {v} placeholder.
func expand(tmpl, value string) string {
	return strings.ReplaceAll(tmpl, "{v}", value)
}

func headingNumber(n int, style string) string {
	if style == domain.NumberStyleArabic {
		return strconv.Itoa(n) + ". "
	}
	return chineseNumber(n) + "、"
}

func sectionNumber(n int, style string) string {
	if style == domain.NumberStyleArabic {
		return "(" + strconv.Itoa(n) + ") "
	}
	return "（" + chineseNumber(n) + "）"
}

// indexLabel is the number placed before a chapter name in an index entry: the Chinese style writes it in Chinese numerals followed by the ideographic comma, while the Arabic style relies on the ordered-list number.
func indexLabel(n int, style string) string {
	if style == domain.NumberStyleArabic {
		return ""
	}
	return chineseNumber(n) + "、"
}

func chineseNumber(n int) string {
	if n >= 1 && n <= len(chineseNumbers) {
		return chineseNumbers[n-1]
	}
	return strconv.Itoa(n)
}
