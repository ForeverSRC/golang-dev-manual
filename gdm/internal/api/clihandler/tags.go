package clihandler

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Tags prints the cross-cutting tag vocabulary with each tag's description in the requested language.
func (cli *CommandLineHandler) Tags(cmd *cobra.Command, _ []string) error {
	lang, err := cmd.Flags().GetString("lang")
	if err != nil {
		return err
	}
	jsonOut, err := cmd.Flags().GetBool("json")
	if err != nil {
		return err
	}

	tags, err := cli.manualSvc.Tags(cmd.Context(), lang)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if jsonOut {
		return writeJSON(out, tags)
	}
	for _, t := range tags {
		if _, err := fmt.Fprintf(out, "%s: %s\n", t.ID, t.Description); err != nil {
			return err
		}
	}
	return nil
}
