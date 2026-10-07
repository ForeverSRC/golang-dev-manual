package clihandler

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/service"
)

// Explain renders the details, examples, and sources of the given clauses.
func (cli *CommandLineHandler) Explain(cmd *cobra.Command, args []string) error {
	lang, err := cmd.Flags().GetString("lang")
	if err != nil {
		return err
	}
	labels, err := cli.manualSvc.Labels(cmd.Context(), lang)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	for i, id := range args {
		view, err := cli.manualSvc.Explain(cmd.Context(), lang, id)
		if err != nil {
			return err
		}
		if i > 0 {
			if _, err := fmt.Fprintln(out); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(out, clauseText(view, labels)); err != nil {
			return err
		}
	}
	return nil
}

// clauseText renders one clause as plain text, with the field labels taken from the language overlay.
func clauseText(v service.ClauseView, l domain.I18nCli) string {
	c := v.Clause
	var b strings.Builder
	fmt.Fprintf(&b, "【%s】%s %s\n\n", c.Level, c.ID, c.Summary)
	fmt.Fprintf(&b, "%s: %s | %s: %s\n\n", l.Category, v.Category, l.SinceGo, c.SinceGo)
	if c.Details != "" {
		fmt.Fprintf(&b, "%s:\n%s\n\n", l.Details, c.Details)
	}
	if c.Quote.Text != "" || c.Rationale != "" {
		fmt.Fprintf(&b, "%s:\n", l.Why)
		if c.Quote.Text != "" {
			fmt.Fprintf(&b, "%s: %s\n", l.Quote, c.Quote.Text)
			if c.Quote.Source != "" {
				fmt.Fprintf(&b, "%s: %s\n", l.Source, c.Quote.Source)
			}
		}
		if c.Rationale != "" {
			fmt.Fprintf(&b, "%s\n", c.Rationale)
		}
		b.WriteString("\n")
	}
	if c.Examples.Good != "" {
		fmt.Fprintf(&b, "%s:\n%s\n\n", l.Good, codeBlock(c.Examples.Good, c.Examples.Language()))
	}
	if c.Examples.Bad != "" {
		fmt.Fprintf(&b, "%s:\n%s\n\n", l.Bad, codeBlock(c.Examples.Bad, c.Examples.Language()))
	}
	if len(c.Sources) > 0 {
		fmt.Fprintf(&b, "%s:\n", l.References)
		for _, s := range c.Sources {
			fmt.Fprintf(&b, "- %s\n", s)
		}
		b.WriteString("\n")
	}
	fmt.Fprintf(&b, "%s: %s\n", l.Detection, v.Detect)
	return b.String()
}

func codeBlock(code, lang string) string {
	return "```" + lang + "\n" + strings.TrimRight(code, "\n") + "\n```"
}
