package cmd

import (
	"fmt"

	"github.com/dbz/dbz/pkg/containers"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [database]",
	Short: "Start a stopped database container",
	Long: `Start a previously stopped database container by name or ID.

Examples:
  dbz start postgres
  dbz start my-mysql-container
`,
	Args: cobra.ExactArgs(1),
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) error {
	containerName := args[0]

	err := containers.StartContainer(containerName)
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	fmt.Printf("✅ Started container: %s\n", containerName)
	return nil
}
