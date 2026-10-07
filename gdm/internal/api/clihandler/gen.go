package clihandler

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Gen generates the markdown manual from the clause data.
func (cli *CommandLineHandler) Gen(cmd *cobra.Command, _ []string) error {
	outDir, err := cmd.Flags().GetString("out")
	if err != nil {
		return err
	}
	lang, err := cmd.Flags().GetString("lang")
	if err != nil {
		return err
	}

	paths, err := cli.genSvc.Generate(cmd.Context(), outDir, splitList(lang))
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	for _, p := range paths {
		if _, err := fmt.Fprintf(out, "generated %s\n", p); err != nil {
			return err
		}
	}
	return nil
}
