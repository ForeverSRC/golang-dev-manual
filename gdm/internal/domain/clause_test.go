package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
)

func TestFilter(t *testing.T) {
	manual := &domain.Manual{
		Clauses: []domain.Clause{
			{ID: "A-001", Level: domain.LevelMust, Chapter: "programming-conventions", Section: "naming"},
			{ID: "A-002", Level: domain.LevelShould, Chapter: "programming-conventions", Section: "comments"},
			{ID: "B-001", Level: domain.LevelMust, Chapter: "errors-and-logging", Section: "logging"},
			{ID: "C-001", Level: domain.LevelMay, Chapter: "performance"},
		},
	}

	tests := []struct {
		name     string
		levels   []string
		category string
		want     []domain.Clause
	}{
		{
			name: "should return all clauses when no filter given",
			want: []domain.Clause{
				{ID: "A-001", Level: domain.LevelMust, Chapter: "programming-conventions", Section: "naming"},
				{ID: "A-002", Level: domain.LevelShould, Chapter: "programming-conventions", Section: "comments"},
				{ID: "B-001", Level: domain.LevelMust, Chapter: "errors-and-logging", Section: "logging"},
				{ID: "C-001", Level: domain.LevelMay, Chapter: "performance"},
			},
		},
		{
			name:   "should filter by level when level given",
			levels: []string{"MUST"},
			want: []domain.Clause{
				{ID: "A-001", Level: domain.LevelMust, Chapter: "programming-conventions", Section: "naming"},
				{ID: "B-001", Level: domain.LevelMust, Chapter: "errors-and-logging", Section: "logging"},
			},
		},
		{
			name:   "should ignore level case when filtering by level",
			levels: []string{"must", "may"},
			want: []domain.Clause{
				{ID: "A-001", Level: domain.LevelMust, Chapter: "programming-conventions", Section: "naming"},
				{ID: "B-001", Level: domain.LevelMust, Chapter: "errors-and-logging", Section: "logging"},
				{ID: "C-001", Level: domain.LevelMay, Chapter: "performance"},
			},
		},
		{
			name:     "should filter by chapter when chapter id given",
			category: "programming-conventions",
			want: []domain.Clause{
				{ID: "A-001", Level: domain.LevelMust, Chapter: "programming-conventions", Section: "naming"},
				{ID: "A-002", Level: domain.LevelShould, Chapter: "programming-conventions", Section: "comments"},
			},
		},
		{
			name:     "should filter by section when section id given",
			category: "programming-conventions/naming",
			want: []domain.Clause{
				{ID: "A-001", Level: domain.LevelMust, Chapter: "programming-conventions", Section: "naming"},
			},
		},
		{
			name:     "should filter by level and chapter when both given",
			levels:   []string{"MUST"},
			category: "programming-conventions",
			want: []domain.Clause{
				{ID: "A-001", Level: domain.LevelMust, Chapter: "programming-conventions", Section: "naming"},
			},
		},
		{
			name:   "should return nothing when nothing matches",
			levels: []string{"NOPE"},
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, manual.Filter(tt.levels, tt.category))
		})
	}
}

func TestByID(t *testing.T) {
	manual := &domain.Manual{Clauses: []domain.Clause{{ID: "A-001", Summary: "x"}}}

	got, ok := manual.ByID("A-001")
	assert.True(t, ok)
	assert.Equal(t, domain.Clause{ID: "A-001", Summary: "x"}, got)

	missing, ok := manual.ByID("A-999")
	assert.False(t, ok)
	assert.Equal(t, domain.Clause{}, missing)
}

func TestLevelRank(t *testing.T) {
	tests := []struct {
		name  string
		level domain.Level
		want  int
	}{
		{name: "should rank MUST first", level: domain.LevelMust, want: 0},
		{name: "should rank SHOULD after MUST", level: domain.LevelShould, want: 1},
		{name: "should rank MAY last", level: domain.LevelMay, want: 2},
		{name: "should rank unknown level last", level: domain.Level("NOPE"), want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.level.Rank())
		})
	}
}

func TestCategoryID(t *testing.T) {
	tests := []struct {
		name   string
		clause domain.Clause
		want   string
	}{
		{
			name:   "should join chapter and section when the clause has a section",
			clause: domain.Clause{Chapter: "programming-conventions", Section: "naming"},
			want:   "programming-conventions/naming",
		},
		{
			name:   "should report the chapter when the clause has no section",
			clause: domain.Clause{Chapter: "performance"},
			want:   "performance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.clause.CategoryID())
		})
	}
}

func TestCategoryIDs(t *testing.T) {
	manual := &domain.Manual{ToC: []domain.Chapter{
		{ID: "programming-conventions", Sections: []domain.Section{{ID: "naming"}, {ID: "comments"}}},
		{ID: "performance"},
	}}

	assert.Equal(t, []string{
		"programming-conventions",
		"programming-conventions/naming",
		"programming-conventions/comments",
		"performance",
	}, manual.CategoryIDs())
}
