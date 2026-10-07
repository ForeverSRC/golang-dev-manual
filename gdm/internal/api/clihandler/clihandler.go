// Package clihandler implements the CLI inbound adapter shared by gdm-cli and gdm-gen: it parses command arguments, calls the application service, and renders output.
package clihandler

import (
	"strings"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/service"
)

// CommandLineHandler holds the application services, called by each cobra command's RunE.
type CommandLineHandler struct {
	manualSvc service.ManualService
	genSvc    service.Generator
}

// New creates the CLI inbound adapter.
func New(manualSvc service.ManualService, genSvc service.Generator) *CommandLineHandler {
	return &CommandLineHandler{manualSvc: manualSvc, genSvc: genSvc}
}

// splitList splits a comma-separated argument into a list of non-empty items.
func splitList(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
