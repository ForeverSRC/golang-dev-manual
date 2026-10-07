package main_test

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ForeverSRC/golang-dev-manual/gdm/cmd/gdm-cli/wire"
)

var (
	// clauseLinePattern matches one line of list output "id: summary".
	clauseLinePattern = regexp.MustCompile(`^([A-Z]+-\d{3}): (.+)$`)
	// explain prints its field labels in the requested language, so each language has its own patterns (the category may contain spaces).
	zhCategoryPattern = regexp.MustCompile(`归属: \S`)
	zhSinceGoPattern  = regexp.MustCompile(`起始版本: \S`)
	enCategoryPattern = regexp.MustCompile(`Category: \S`)
	enSinceGoPattern  = regexp.MustCompile(`Since: \S`)
)

// CLIITSuite is gdm's command-line integration test suite: it takes the repository's real clause data as input,
// runs the full command tree in cobra's official execution pattern, and asserts end-to-end output and artifacts.
// Every execution builds a fresh command tree via wire, so flags don't leak between cases.
type CLIITSuite struct {
	suite.Suite
}

func TestCLIITSuite(t *testing.T) {
	suite.Run(t, new(CLIITSuite))
}

func (s *CLIITSuite) TestList() {
	s.Run("should list unique clause ids from real data", func() {
		ids := s.listClauseIDs()
		s.NotEmpty(ids)

		seen := make(map[string]bool, len(ids))
		for _, id := range ids {
			s.False(seen[id], "条款编号重复: %s", id)
			seen[id] = true
		}
	})

	s.Run("should cover every clause when all legal levels given", func() {
		all := s.listClauseIDs()
		byLevel := s.listClauseIDs("--level", "MUST,SHOULD,MAY")
		s.Equal(all, byLevel, "三个合法级别应覆盖全量，漏出的条款分级非法")
	})

	s.Run("should print the summary in the requested language", func() {
		zh := s.listOutput("--level", "MUST", "--lang", "zh")
		en := s.listOutput("--level", "MUST", "--lang", "en")

		s.NotEqual(zh, en, "两种语言的条款摘要应不同")
		s.Len(strings.Split(en, "\n"), len(strings.Split(zh, "\n")), "两种语言的行数应一致")
	})

	s.Run("should filter by chapter id and by section id", func() {
		chapter := s.listClauseIDs("--category", "programming-conventions")
		section := s.listClauseIDs("--category", "programming-conventions/naming")

		s.NotEmpty(section)
		s.Subset(chapter, section, "章节过滤的结果应包含其小节过滤的结果")
	})

	s.Run("should list available category ids when category does not exist", func() {
		out, _, err := s.executeCommand("list", "--category", "nope")

		s.Require().Error(err)
		s.Contains(err.Error(), "programming-conventions/naming", "报错应列出可选归属")
		s.Empty(out)
	})
}

func (s *CLIITSuite) TestExplain() {
	s.Run("should render every clause with the Chinese labels by default", func() {
		ids := s.listClauseIDs()
		s.Require().NotEmpty(ids)

		for _, id := range ids {
			out, errOut, err := s.executeCommand("explain", id)
			s.Require().NoError(err, "条款 %s 应可展开，stderr: %s", id, errOut)
			s.Regexp(`【(MUST|SHOULD|MAY)】`+regexp.QuoteMeta(id)+` .+`, out, "条款 %s 缺少合法分级或摘要", id)
			s.Regexp(zhCategoryPattern, out, "条款 %s 缺少归属", id)
			s.Regexp(zhSinceGoPattern, out, "条款 %s 缺少起始版本", id)
			s.Contains(out, "说明:", "条款 %s 缺少说明", id)
			s.Contains(out, "为什么:", "条款 %s 缺少依据", id)
		}
	})

	s.Run("should render the English labels when lang is en", func() {
		out, errOut, err := s.executeCommand("explain", "NAMING-001", "--lang", "en")
		s.Require().NoError(err, "stderr: %s", errOut)

		s.Regexp(enCategoryPattern, out, "缺少归属")
		s.Regexp(enSinceGoPattern, out, "缺少起始版本")
		s.Contains(out, "Details:")
		s.Contains(out, "Why:")
		s.NotContains(out, "说明:")
	})
}

func (s *CLIITSuite) TestErrors() {
	s.Run("should fail when clause id does not exist", func() {
		_, _, err := s.executeCommand("explain", "NOPE-999")
		s.Require().Error(err)
	})
}

// executeCommand builds a fresh command tree, runs the subcommand, and returns the captured stdout, stderr, and execution error.
func (s *CLIITSuite) executeCommand(args ...string) (string, string, error) {
	s.T().Helper()

	container, err := wire.InitializeCLI()
	s.Require().NoError(err)

	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)
	root := container.RootCmd
	root.SetOut(out)
	root.SetErr(errOut)
	root.SetArgs(args)

	_, err = root.ExecuteC()
	return out.String(), errOut.String(), err
}

// listOutput runs list and returns the trimmed stdout.
func (s *CLIITSuite) listOutput(args ...string) string {
	s.T().Helper()

	out, errOut, err := s.executeCommand(append([]string{"list"}, args...)...)
	s.Require().NoError(err, "stderr: %s", errOut)
	return strings.TrimSpace(out)
}

// listClauseIDs runs list and parses all ids, verifying that every line is "id: summary".
func (s *CLIITSuite) listClauseIDs(args ...string) []string {
	s.T().Helper()

	var ids []string
	for line := range strings.SplitSeq(s.listOutput(args...), "\n") {
		if line == "" {
			continue
		}
		m := clauseLinePattern.FindStringSubmatch(line)
		s.Require().NotNil(m, "list 输出行应为「编号: 一句话」: %q", line)
		ids = append(ids, m[1])
	}
	return ids
}
