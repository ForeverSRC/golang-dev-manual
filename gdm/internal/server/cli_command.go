// Package server assembles gdm's cobra command tree.
package server

import (
	"github.com/spf13/cobra"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/api/clihandler"
)

// CLICommand is gdm's root command.
type CLICommand struct {
	*cobra.Command
}

// ProvideCLICommand builds gdm's command tree.
func ProvideCLICommand(handler *clihandler.CommandLineHandler) CLICommand {
	rootCmd := &cobra.Command{
		Use:           "gdm-cli",
		Short:         "Golang development manual CLI",
		Long:          "gdm-cli searches the clauses of the Golang development manual and runs the automatable checks against Go code.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(newListCmd(handler))
	rootCmd.AddCommand(newSearchCmd(handler))
	rootCmd.AddCommand(newTagsCmd(handler))
	rootCmd.AddCommand(newExplainCmd(handler))
	rootCmd.AddCommand(newCheckCmd(handler))
	rootCmd.PersistentFlags().String("lang", "zh", "language of the printed clause text: zh or en")

	return CLICommand{Command: rootCmd}
}

func newListCmd(handler *clihandler.CommandLineHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: `Filter by level, category, or tag, printing only "id: summary"`,
		Args:  cobra.NoArgs,
		RunE:  handler.List,
	}
	addClauseFilters(cmd)
	return cmd
}

func newSearchCmd(handler *clihandler.CommandLineHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <query> [query...]",
		Short: "Score the clauses against a query and print the hits in relevance order",
		Args:  cobra.MinimumNArgs(1),
		RunE:  handler.Search,
	}
	addClauseFilters(cmd)
	return cmd
}

// addClauseFilters registers the filters shared by list and search.
func addClauseFilters(cmd *cobra.Command) {
	cmd.Flags().String("level", "", "filter by level, comma-separated, e.g. MUST,SHOULD")
	cmd.Flags().String("category", "", `filter by category id: a chapter id or "chapter/section"`)
	cmd.Flags().String("tag", "", "filter by tag, comma-separated; a clause matches any of the given tags")
	cmd.Flags().Bool("json", false, "print a JSON document instead of plain text")
}

func newTagsCmd(handler *clihandler.CommandLineHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "List the cross-cutting tags with their descriptions",
		Args:  cobra.NoArgs,
		RunE:  handler.Tags,
	}
	cmd.Flags().Bool("json", false, "print a JSON document instead of plain text")
	return cmd
}

func newExplainCmd(handler *clihandler.CommandLineHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "explain <id> [id...]",
		Short: "Show the details, examples, and references of the given clauses",
		Args:  cobra.MinimumNArgs(1),
		RunE:  handler.Explain,
	}
	return cmd
}

func newCheckCmd(handler *clihandler.CommandLineHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check <path>",
		Short: "Run the automatable checks against a code directory and print the hits",
		Args:  cobra.ExactArgs(1),
		RunE:  handler.Check,
	}
	cmd.Flags().Bool("verbose", false, "also list the clauses without an automated check")
	cmd.Flags().Bool("json", false, "print a JSON document instead of plain text")
	return cmd
}
