package clihandler

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Check runs the automatable checks against a code directory and outputs the hits.
func (cli *CommandLineHandler) Check(cmd *cobra.Command, args []string) error {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}
	jsonOut, err := cmd.Flags().GetBool("json")
	if err != nil {
		return err
	}
	lang, err := cmd.Flags().GetString("lang")
	if err != nil {
		return err
	}

	report, err := cli.manualSvc.Check(cmd.Context(), lang, args[0])
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if jsonOut {
		return writeJSON(out, checkReportJSON(report))
	}
	labels, err := cli.manualSvc.Labels(cmd.Context(), lang)
	if err != nil {
		return err
	}
	if len(report.Hits) == 0 {
		if _, err := fmt.Fprintln(out, labels.NoHits); err != nil {
			return err
		}
	}
	for _, h := range report.Hits {
		if _, err := fmt.Fprintf(out, "%s:%d: %s  [%s %s] %s\n",
			h.Path, h.Line, h.Text, h.Rule.ClauseID, h.Rule.Level, h.Rule.Summary); err != nil {
			return err
		}
	}
	if verbose && len(report.Uncovered) > 0 {
		if _, err := fmt.Fprintf(out, "\n%s\n", labels.Uncovered); err != nil {
			return err
		}
		for _, v := range report.Uncovered {
			if _, err := fmt.Fprintf(out, "  [%s] %s\n", v.Clause.ID, v.Detect); err != nil {
				return err
			}
		}
	}
	return nil
}
