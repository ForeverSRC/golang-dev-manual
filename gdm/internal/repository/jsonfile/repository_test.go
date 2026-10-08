package jsonfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/repository/jsonfile"
)

type RepositorySuite struct {
	suite.Suite

	underTest *jsonfile.Repository
}

func (s *RepositorySuite) SetupSuite() {
	s.underTest = jsonfile.New(os.DirFS("testdata"))
}

func (s *RepositorySuite) TestLoad() {
	got, err := s.underTest.Load(s.T().Context())
	s.Require().NoError(err)

	want := &domain.Manual{
		Version:    "1.0",
		GoBaseline: "1.27.1",
		ToC: []domain.Chapter{
			{ID: "programming-conventions", Name: "编程规约", Sections: []domain.Section{{ID: "naming", Name: "命名规约"}}},
			{ID: "performance", Name: "性能规约"},
		},
		Clauses: []domain.Clause{
			{
				ID:        "NAMING-001",
				Level:     domain.LevelMust,
				Chapter:   "programming-conventions",
				Section:   "naming",
				SinceGo:   "1.0",
				Summary:   "包名使用小写单词连写。",
				Details:   "包名全小写、无分隔符，多个单词直接连写。",
				Rationale: "包名出现在每个外部引用点，命名不一致会持续抬高阅读成本。",
				Examples:  domain.Examples{Good: "package orderbook", Bad: "package order_book"},
				Sources:   []string{"https://go.dev/wiki/CodeReviewComments#package-names"},
				Detect: &domain.Detect{
					Tool:    "golangci-lint",
					Rule:    "stylecheck(ST1003)",
					Pattern: "", // zero value written explicitly to show the assertion is intentional
					Note:    "检查包名标识符",
				},
			},
			{
				ID:        "PERF-001",
				Level:     domain.LevelMay,
				Chapter:   "performance",
				SinceGo:   "1.0",
				Summary:   "按已知容量预分配切片。",
				Details:   "能预知元素数量时用 make 指定容量。",
				Rationale: "容量不足会触发多次扩容与整体复制。",
				Examples:  domain.Examples{Good: "ids := make([]string, 0, len(clauses))", Bad: "var ids []string"},
				Sources:   []string{"https://go.dev/wiki/CodeReviewComments"},
				Detect:    nil, // zero value written explicitly to show the assertion is intentional
			},
		},
	}

	s.Equal(want, got)
}

func (s *RepositorySuite) TestLoadError() {
	tests := []struct {
		name    string
		prepare func() string
	}{
		{
			name: "should return error when json is invalid",
			prepare: func() string {
				dir := s.T().TempDir()
				s.Require().NoError(os.WriteFile(filepath.Join(dir, "manual.json"), []byte("{"), 0o600))
				return dir
			},
		},
		{
			name:    "should return error when file is missing",
			prepare: func() string { return s.T().TempDir() },
		},
		{
			name: "should return error when a chapter file is missing",
			prepare: func() string {
				dir := s.T().TempDir()
				meta := `{"toc":[{"id":"missing","name":"x"}]}`
				s.Require().NoError(os.WriteFile(filepath.Join(dir, "manual.json"), []byte(meta), 0o600))
				return dir
			},
		},
		{
			name: "should return error when a chapter id repeats",
			prepare: func() string {
				dir := s.T().TempDir()
				writeFile(s, filepath.Join(dir, "manual.json"), `{"toc":[{"id":"x","name":"x"},{"id":"x","name":"x"}]}`)
				writeFile(s, filepath.Join(dir, "clauses", "x.json"), `[]`)
				return dir
			},
		},
		{
			name: "should return error when a clause id repeats",
			prepare: func() string {
				dir := s.T().TempDir()
				writeFile(s, filepath.Join(dir, "manual.json"), `{"toc":[{"id":"x","name":"x"}]}`)
				writeFile(s, filepath.Join(dir, "clauses", "x.json"), `[{"id":"A-001","chapter":"x"},{"id":"A-001","chapter":"x"}]`)
				return dir
			},
		},
		{
			name: "should return error when a clause references an unknown section",
			prepare: func() string {
				dir := s.T().TempDir()
				writeFile(s, filepath.Join(dir, "manual.json"), `{"toc":[{"id":"x","name":"x","sections":[{"id":"y","name":"y"}]}]}`)
				writeFile(s, filepath.Join(dir, "clauses", "x.json"), `[{"id":"A-001","chapter":"x","section":"z"}]`)
				return dir
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			_, err := jsonfile.New(os.DirFS(tt.prepare())).Load(s.T().Context())
			s.Require().Error(err)
		})
	}
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

// writeFile writes content under a path, creating the parent directories.
func writeFile(s *RepositorySuite, path, content string) {
	s.T().Helper()
	s.Require().NoError(os.MkdirAll(filepath.Dir(path), 0o750))
	s.Require().NoError(os.WriteFile(path, []byte(content), 0o600))
}
