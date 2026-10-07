package main_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ForeverSRC/golang-dev-manual/gdm/cmd/gdm-gen/wire"
	"github.com/ForeverSRC/golang-dev-manual/gdm/data"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/repository/jsonfile"
)

// GenITSuite is gdm-gen's command-line integration test suite: it takes the embedded clause data as input,
// runs the command tree in cobra's official execution pattern, and asserts the generated artifacts.
type GenITSuite struct {
	suite.Suite
}

func TestGenITSuite(t *testing.T) {
	suite.Run(t, new(GenITSuite))
}

func (s *GenITSuite) TestGenerate() {
	s.Run("should place every clause into a chapter file", func() {
		manual, err := jsonfile.New(data.FS).Load(s.T().Context())
		s.Require().NoError(err)
		s.Require().NotEmpty(manual.Clauses)

		outDir := s.T().TempDir()
		out, errOut, err := s.execute("--out", outDir)
		s.Require().NoError(err, "stderr: %s", errOut)
		s.Contains(out, "generated ")

		chapters := s.chapterContent(filepath.Join(outDir, "zh"))
		for _, c := range manual.Clauses {
			s.Contains(chapters, c.ID, "条款 %s 未落入任何章节文件", c.ID)
		}
	})
}

// execute builds a fresh command tree, runs the generator, and returns the captured stdout, stderr, and execution error.
func (s *GenITSuite) execute(args ...string) (string, string, error) {
	s.T().Helper()

	container, err := wire.InitializeGenerator()
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

// chapterContent concatenates the chapter file contents of one language under the output directory, excluding the appendix and index.
func (s *GenITSuite) chapterContent(dir string) string {
	s.T().Helper()

	entries, err := os.ReadDir(dir)
	s.Require().NoError(err)

	var b strings.Builder
	for _, e := range entries {
		name := e.Name()
		if name == "index.md" || strings.Contains(name, "appendix") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, name)) //nolint:gosec // the path comes from the test temp directory
		s.Require().NoError(err)
		b.Write(raw)
	}
	return b.String()
}
