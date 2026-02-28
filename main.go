package main

import (
	"fmt"
	"os"

	"github.com/dbz/dbz/cmd"
)

// Version is set at build time via ldflags
var Version = "dev"

func main() {
	cmd.SetVersion(Version)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
