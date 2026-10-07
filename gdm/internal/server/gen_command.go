package server

import (
	"github.com/spf13/cobra"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/api/clihandler"
)

// GeneratorCommand is gdm-gen's root command.
type GeneratorCommand struct {
	*cobra.Command
}

// ProvideGeneratorCommand builds gdm-gen's command, which generates the markdown manual from the embedded clause data.
func ProvideGeneratorCommand(handler *clihandler.CommandLineHandler) GeneratorCommand {
	cmd := &cobra.Command{
		Use:           "gdm-gen",
		Short:         "Generate the markdown manual from the clause data",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE:          handler.Gen,
	}
	cmd.Flags().String("out", "manual", "manual output directory")
	cmd.Flags().String("lang", "zh,en", "languages to generate, comma-separated; output goes to <--out>/<language>/")
	return GeneratorCommand{Command: cmd}
}
