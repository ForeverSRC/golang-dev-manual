package service

import (
	"context"
	"path/filepath"
)

type generatorService struct {
	repo   ClauseRepository
	i18n   I18nRepository
	writer ManualWriter
}

// NewGenerator creates the manual generator service.
func NewGenerator(repo ClauseRepository, i18n I18nRepository, writer ManualWriter) Generator {
	return &generatorService{repo: repo, i18n: i18n, writer: writer}
}

func (s *generatorService) Generate(ctx context.Context, outDir string, langs []string) ([]string, error) {
	manual, err := s.repo.Load(ctx)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, lang := range langs {
		overlay, err := s.i18n.LoadI18n(ctx, lang)
		if err != nil {
			return nil, err
		}
		files, err := renderManual(manual, overlay, lang)
		if err != nil {
			return nil, err
		}
		written, err := s.writer.Write(ctx, filepath.Join(outDir, lang), files)
		if err != nil {
			return nil, err
		}
		paths = append(paths, written...)
	}
	return paths, nil
}
