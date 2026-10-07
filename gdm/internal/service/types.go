package service

import (
	"regexp"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
)

// GrepRule is an executable grep check rule, compiled from a clause's detect field.
type GrepRule struct {
	ClauseID string
	Level    domain.Level
	Summary  string
	Pattern  *regexp.Regexp
}

// GrepHit is one hit record.
type GrepHit struct {
	Path string
	Line int
	Text string
	Rule GrepRule
}

// CheckReport is the result of one check.
type CheckReport struct {
	Hits      []GrepHit
	Uncovered []ClauseView
}
