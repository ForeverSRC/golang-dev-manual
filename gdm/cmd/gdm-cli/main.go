// gdm-cli is a CLI for searching and checking the Golang development manual.
package main

import (
	"fmt"
	"os"

	"github.com/ForeverSRC/golang-dev-manual/gdm/cmd/gdm-cli/wire"
)

func main() {
	container, err := wire.InitializeCLI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize gdm-cli: %v\n", err)
		os.Exit(1)
	}

	if err := container.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
