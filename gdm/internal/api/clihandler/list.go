package clihandler

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/service"
)

// List outputs "id: summary" so an agent can get the clause list cheaply, or the same clauses as a JSON document.
func (cli *CommandLineHandler) List(cmd *cobra.Command, _ []string) error {
	f, err := readFilters(cmd)
	if err != nil {
		return err
	}

	views, err := cli.manualSvc.List(cmd.Context(), service.ClauseQuery{
		Lang:     f.lang,
		Levels:   f.levels,
		Category: f.category,
		Tags:     f.tags,
	})
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if f.jsonOut {
		return writeJSON(out, clauseJSONs(views))
	}
	for _, v := range views {
		if _, err := fmt.Fprintf(out, "%s: %s\n", v.Clause.ID, v.Clause.Summary); err != nil {
			return err
		}
	}
	return nil
}
