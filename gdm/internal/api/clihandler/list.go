package clihandler

import (
	"fmt"

	"github.com/spf13/cobra"
)

// List outputs "id: summary" so an agent can get the clause list cheaply.
func (cli *CommandLineHandler) List(cmd *cobra.Command, _ []string) error {
	lang, err := cmd.Flags().GetString("lang")
	if err != nil {
		return err
	}
	level, err := cmd.Flags().GetString("level")
	if err != nil {
		return err
	}
	category, err := cmd.Flags().GetString("category")
	if err != nil {
		return err
	}

	clauses, err := cli.manualSvc.List(cmd.Context(), lang, splitList(level), category)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	for _, c := range clauses {
		if _, err := fmt.Fprintf(out, "%s: %s\n", c.ID, c.Summary); err != nil {
			return err
		}
	}
	return nil
}
