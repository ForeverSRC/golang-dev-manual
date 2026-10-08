package clihandler

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Search scores the clauses against a query and prints the hits in relevance order.
func (cli *CommandLineHandler) Search(cmd *cobra.Command, args []string) error {
	f, err := readFilters(cmd)
	if err != nil {
		return err
	}

	results, err := cli.manualSvc.Search(cmd.Context(), f.lang, strings.Join(args, " "), f.levels, f.category, f.tags)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if f.jsonOut {
		return writeJSON(out, searchJSONs(results))
	}
	for _, r := range results {
		if _, err := fmt.Fprintf(out, "%s: %s\n", r.Clause.ID, r.Clause.Summary); err != nil {
			return err
		}
	}
	return nil
}
