// Package clihandler implements the CLI inbound adapter shared by gdm-cli and gdm-gen: it parses command arguments, calls the application service, and renders output.
package clihandler

import (
	"strings"

	"github.com/spf13/cobra"

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

// clauseFilters holds the filters list and search share, with the comma-separated arguments already split.
type clauseFilters struct {
	lang     string
	levels   []string
	category string
	tags     []string
	jsonOut  bool
}

// readFilters collects the flags shared by list and search.
func readFilters(cmd *cobra.Command) (clauseFilters, error) {
	lang, err := cmd.Flags().GetString("lang")
	if err != nil {
		return clauseFilters{}, err
	}
	level, err := cmd.Flags().GetString("level")
	if err != nil {
		return clauseFilters{}, err
	}
	category, err := cmd.Flags().GetString("category")
	if err != nil {
		return clauseFilters{}, err
	}
	tag, err := cmd.Flags().GetString("tag")
	if err != nil {
		return clauseFilters{}, err
	}
	jsonOut, err := cmd.Flags().GetBool("json")
	if err != nil {
		return clauseFilters{}, err
	}
	return clauseFilters{
		lang:     lang,
		levels:   splitList(level),
		category: category,
		tags:     splitList(tag),
		jsonOut:  jsonOut,
	}, nil
}
