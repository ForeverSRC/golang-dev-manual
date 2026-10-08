package domain

import "strings"

// searchFields lists the clause fields a query is matched against, ordered by descending weight.
// The weight says how decisive a hit in that field is: a hit in the one-line summary means more than one in the rationale.
var searchFields = []struct {
	name   string
	weight int
	text   func(Clause) string
}{
	{"summary", 4, func(c Clause) string { return c.Summary }},
	{"tags", 3, func(c Clause) string { return strings.Join(c.Tags, " ") }},
	{"details", 2, func(c Clause) string { return c.Details }},
	{"rationale", 1, func(c Clause) string { return c.Rationale }},
}

// Match is one query term found in a clause field.
type Match struct {
	Field string `json:"field"`
	Term  string `json:"term"`
}

// MatchScore scores a clause against the query terms: a field contributes its weight once per matched term,
// and the matches record which term hit which field.
func (c Clause) MatchScore(terms []string) (int, []Match) {
	score := 0
	var matches []Match
	for _, f := range searchFields {
		hits := matchTerms(terms, f.text(c))
		if len(hits) == 0 {
			continue
		}
		score += f.weight * len(hits)
		for _, term := range hits {
			matches = append(matches, Match{Field: f.name, Term: term})
		}
	}
	return score, matches
}

// Tokenize splits text into the terms a query can hit: an ASCII run becomes a whole word split on
// non-alphanumeric characters, and a Han run is split into its adjacent character pairs,
// keeping a lone character as-is.
func Tokenize(text string) []string {
	runes := []rune(strings.ToLower(text))
	var terms []string
	for i := 0; i < len(runes); {
		switch {
		case isASCIIAlnum(runes[i]):
			j := i
			for j < len(runes) && isASCIIAlnum(runes[j]) {
				j++
			}
			terms = append(terms, string(runes[i:j]))
			i = j
		case isHan(runes[i]):
			j := i
			for j < len(runes) && isHan(runes[j]) {
				j++
			}
			terms = append(terms, hanBigrams(runes[i:j])...)
			i = j
		default:
			i++
		}
	}
	return terms
}

// matchTerms returns the query terms present in text, deduplicated and in query order.
func matchTerms(terms []string, text string) []string {
	tokens := make(map[string]bool)
	for _, t := range Tokenize(text) {
		tokens[t] = true
	}
	seen := make(map[string]bool, len(terms))
	var hits []string
	for _, t := range terms {
		if tokens[t] && !seen[t] {
			seen[t] = true
			hits = append(hits, t)
		}
	}
	return hits
}

// hanBigrams splits a Han run into its adjacent character pairs, keeping a lone character as-is.
func hanBigrams(runes []rune) []string {
	if len(runes) < 2 {
		return []string{string(runes)}
	}
	pairs := make([]string, 0, len(runes)-1)
	for i := 0; i+1 < len(runes); i++ {
		pairs = append(pairs, string(runes[i:i+2]))
	}
	return pairs
}

func isASCIIAlnum(r rune) bool {
	return r >= '0' && r <= '9' || r >= 'a' && r <= 'z'
}

func isHan(r rune) bool {
	return r >= 0x4E00 && r <= 0x9FFF
}
