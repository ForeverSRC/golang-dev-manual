package service

import (
	"context"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
)

// ClauseRepository is the clause data reading port, implemented by an adapter in the repository layer.
type ClauseRepository interface {
	Load(ctx context.Context) (*domain.Manual, error)
}

// I18nRepository is the language overlay reading port, implemented by an adapter in the repository layer.
type I18nRepository interface {
	LoadI18n(ctx context.Context, lang string) (*domain.I18n, error)
}

// GrepScanner is the source scanning port, implemented by an adapter in the adapter layer.
type GrepScanner interface {
	Scan(ctx context.Context, root string, rules []GrepRule) ([]GrepHit, error)
}

// ManualWriter is the manual file writing port, implemented by an adapter in the adapter layer.
type ManualWriter interface {
	Write(ctx context.Context, dir string, files map[string]string) ([]string, error)
}
