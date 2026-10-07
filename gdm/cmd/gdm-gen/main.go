// gdm-gen generates the markdown manual from the clause data; it is a maintainer and CI tool, not part of the published gdm-cli.
package main

import (
	"fmt"
	"os"

	"github.com/ForeverSRC/golang-dev-manual/gdm/cmd/gdm-gen/wire"
)

func main() {
	container, err := wire.InitializeGenerator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize gdm-gen: %v\n", err)
		os.Exit(1)
	}

	if err := container.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
