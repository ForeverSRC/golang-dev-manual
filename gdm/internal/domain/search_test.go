package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		{
			name: "should split an ASCII run into whole words",
			text: "error-handling / context.Value",
			want: []string{"error", "handling", "context", "value"},
		},
		{
			name: "should lowercase the ASCII words",
			text: "Slices.Clone",
			want: []string{"slices", "clone"},
		},
		{
			name: "should split a Han run into adjacent character pairs",
			text: "切片容量",
			want: []string{"切片", "片容", "容量"},
		},
		{
			name: "should keep a lone Han character as-is",
			text: "锁",
			want: []string{"锁"},
		},
		{
			name: "should keep the terms around punctuation and separate scripts",
			text: "用 tests 验证",
			want: []string{"用", "tests", "验证"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, domain.Tokenize(tt.text))
		})
	}
}

func TestMatchScore(t *testing.T) {
	clause := domain.Clause{
		Summary:   "元素数量已知时用 make 指定切片容量。",
		Details:   "用 make([]T, 0, n) 初始化准备 append 的切片。",
		Rationale: "容量不足会触发多次扩容。",
		Tags:      []string{"performance"},
	}

	tests := []struct {
		name    string
		query   string
		score   int
		matches []domain.Match
	}{
		{
			name:    "should sum the weights of the fields the term hits",
			query:   "make",
			score:   6,
			matches: []domain.Match{{Field: "summary", Term: "make"}, {Field: "details", Term: "make"}},
		},
		{
			name:    "should score a hit in the tags",
			query:   "performance",
			score:   3,
			matches: []domain.Match{{Field: "tags", Term: "performance"}},
		},
		{
			name:    "should score a hit in the rationale",
			query:   "扩容",
			score:   1,
			matches: []domain.Match{{Field: "rationale", Term: "扩容"}},
		},
		{
			name:  "should score nothing when no term hits",
			query: "sql-injection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, matches := clause.MatchScore(domain.Tokenize(tt.query))
			assert.Equal(t, tt.score, score)
			assert.Equal(t, tt.matches, matches)
		})
	}
}
