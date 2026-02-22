package cmd

import (
	"fmt"

	"github.com/dbz/dbz/pkg/containers"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset [database]",
	Short: "Reset a database container",
	Long: `Reset a database container by stopping it (if running) and starting it again.
This is equivalent to "dbz stop <container>" followed by "dbz start <container>".

For a complete reset including data deletion, use "dbz delete" followed by "dbz create".

Examples:
  dbz reset postgres
  dbz reset my-mysql-container
`,
	Args: cobra.ExactArgs(1),
	RunE: runReset,
}

func init() {
	rootCmd.AddCommand(resetCmd)
}

func runReset(cmd *cobra.Command, args []string) error {
	containerName := args[0]

	// First try to stop the container (ignore error if already stopped)
	err := containers.StopContainer(containerName)
	if err != nil {
		// If container doesn't exist, return error
		if err.Error() == fmt.Sprintf("container not found: %s", containerName) {
			return fmt.Errorf("container not found: %s", containerName)
		}
		// Otherwise continue (might already be stopped)
	}

	// Start the container
	err = containers.StartContainer(containerName)
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	fmt.Printf("✅ Reset container: %s\n", containerName)
	return nil
}
