package main_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ForeverSRC/golang-dev-manual/gdm/cmd/gdm-cli/wireit"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/server"
)

// CLIITSuite is gdm-cli's command-line integration test suite: wire assembles the real layers over the fixture data tree
// under internal/ittest, and each case runs the full command tree in cobra's official execution pattern and compares
// the printed output against expectations written out by hand.
type CLIITSuite struct {
	suite.Suite

	container *wireit.CLIITTestContainer
}

func TestCLIITSuite(t *testing.T) {
	suite.Run(t, new(CLIITSuite))
}

func (s *CLIITSuite) SetupSuite() {
	container, err := wireit.InitializeCLIITTestContainer()
	s.Require().NoError(err)
	s.container = container
}

func (s *CLIITSuite) TestList() {
	s.Run("should print every clause in data order", func() {
		s.Equal(`NAMING-001: 包名使用小写单词连写。
COMMENT-001: 注释只写代码表达不出的东西。
PERF-001: 元素数量已知时用 make 指定切片容量。
`, s.output("list"))
	})

	s.Run("should print the summaries in the requested language", func() {
		s.Equal(`NAMING-001: Write package names as lowercase words run together.
COMMENT-001: Comment only what the code cannot express.
PERF-001: Specify a slice capacity with make when the element count is known.
`, s.output("list", "--lang", "en"))
	})

	s.Run("should keep the clauses of the given level", func() {
		s.Equal(`NAMING-001: 包名使用小写单词连写。
`, s.output("list", "--level", "MUST"))
	})

	s.Run("should keep the clauses of several levels", func() {
		s.Equal(`NAMING-001: 包名使用小写单词连写。
PERF-001: 元素数量已知时用 make 指定切片容量。
`, s.output("list", "--level", "MUST,MAY"))
	})

	s.Run("should keep the clauses of the given chapter", func() {
		s.Equal(`NAMING-001: 包名使用小写单词连写。
COMMENT-001: 注释只写代码表达不出的东西。
`, s.output("list", "--category", "programming-conventions"))
	})

	s.Run("should keep the clauses of the given section", func() {
		s.Equal(`COMMENT-001: 注释只写代码表达不出的东西。
`, s.output("list", "--category", "programming-conventions/comments"))
	})

	s.Run("should keep the clauses carrying the given tag", func() {
		s.Equal(`PERF-001: 元素数量已知时用 make 指定切片容量。
`, s.output("list", "--tag", "performance"))
	})

	s.Run("should fail when the category does not exist", func() {
		out, err := s.executeCommand("list", "--category", "nope")

		s.Require().Error(err)
		s.Empty(out)
	})

	s.Run("should fail when the tag does not exist", func() {
		out, err := s.executeCommand("list", "--tag", "nope")

		s.Require().Error(err)
		s.Empty(out)
	})
}

func (s *CLIITSuite) TestSearch() {
	s.Run("should print the hits in relevance order", func() {
		s.Equal(`NAMING-001: 包名使用小写单词连写。
COMMENT-001: 注释只写代码表达不出的东西。
`, s.output("search", "标识符"))
	})

	s.Run("should print the clause whose text carries the query", func() {
		s.Equal(`PERF-001: 元素数量已知时用 make 指定切片容量。
`, s.output("search", "容量"))
	})

	s.Run("should accept the filters list accepts", func() {
		s.Equal(`NAMING-001: 包名使用小写单词连写。
`, s.output("search", "使用", "--level", "MUST", "--tag", "naming"))
	})

	s.Run("should print nothing when no clause matches", func() {
		s.Empty(s.output("search", "并发"))
	})
}

func (s *CLIITSuite) TestTags() {
	s.Run("should print every tag with its description", func() {
		s.Equal(`comments: 注释写法
naming: 包名与标识符命名
performance: 容量预分配与开销
`, s.output("tags"))
	})

	s.Run("should print the descriptions in the requested language", func() {
		s.Equal(`comments: How comments and doc comments are written
naming: Package, identifier, and import naming
performance: Capacity preallocation, string concatenation, and slice cloning
`, s.output("tags", "--lang", "en"))
	})
}

func (s *CLIITSuite) TestExplain() {
	s.Run("should render the clause with the Chinese labels by default", func() {
		s.Equal("【MUST】NAMING-001 包名使用小写单词连写。\n"+
			"\n"+
			"归属: 编程规约/命名规约 | 起始版本: 1.0\n"+
			"\n"+
			"说明:\n"+
			"包名全小写、无下划线与驼峰，标识符同理。\n"+
			"\n"+
			"为什么:\n"+
			"引文: Package names should be short, concise, and evocative.\n"+
			"出处: https://go.dev/blog/package-names\n"+
			"包名出现在每个外部引用点，命名不一致会持续抬高阅读成本。\n"+
			"\n"+
			"正例:\n"+
			"```go\n"+
			"package orderbook\n"+
			"```\n"+
			"\n"+
			"反例:\n"+
			"```go\n"+
			"package order_book\n"+
			"```\n"+
			"\n"+
			"依据:\n"+
			"- https://go.dev/wiki/CodeReviewComments#package-names\n"+
			"\n"+
			"检测: golangci-lint stylecheck(ST1003)（检查包名标识符）\n",
			s.output("explain", "NAMING-001"))
	})

	s.Run("should render the English labels and the translated text when lang is en", func() {
		s.Equal("【MUST】NAMING-001 Write package names as lowercase words run together.\n"+
			"\n"+
			"Category: Programming Conventions/Naming | Since: 1.0\n"+
			"\n"+
			"Details:\n"+
			"Keep package names in lowercase with no underscores or camel case; the same holds for identifiers.\n"+
			"\n"+
			"Why:\n"+
			"Quote: Package names should be short, concise, and evocative.\n"+
			"Source: https://go.dev/blog/package-names\n"+
			"A package name appears at every external reference, so inconsistent naming keeps raising the reading cost.\n"+
			"\n"+
			"Good:\n"+
			"```go\n"+
			"package orderbook\n"+
			"```\n"+
			"\n"+
			"Bad:\n"+
			"```go\n"+
			"package order_book\n"+
			"```\n"+
			"\n"+
			"References:\n"+
			"- https://go.dev/wiki/CodeReviewComments#package-names\n"+
			"\n"+
			"Detection: golangci-lint stylecheck(ST1003) (check the package name identifier)\n",
			s.output("explain", "NAMING-001", "--lang", "en"))
	})
}

func (s *CLIITSuite) TestErrors() {
	s.Run("should fail when the clause id does not exist", func() {
		out, err := s.executeCommand("explain", "NOPE-999")

		s.Require().Error(err)
		s.Empty(out)
	})
}

// executeCommand builds a fresh command tree, runs the subcommand, and returns the captured stdout and execution error.
func (s *CLIITSuite) executeCommand(args ...string) (string, error) {
	s.T().Helper()

	root := server.ProvideCLICommand(s.container.Handler)
	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)
	root.SetOut(out)
	root.SetErr(errOut)
	root.SetArgs(args)

	_, err := root.ExecuteC()
	return out.String(), err
}

// output runs the subcommand and returns its stdout, requiring the execution to succeed.
func (s *CLIITSuite) output(args ...string) string {
	s.T().Helper()

	out, err := s.executeCommand(args...)
	s.Require().NoError(err)
	return out
}
