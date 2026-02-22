package cmd

import (
	"fmt"

	"github.com/dbz/dbz/pkg/containers"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [database]",
	Short: "Stop a database container",
	Long: `Stop a running database container by name or ID.
The container can be started again with the "start" command.

Examples:
  dbz stop postgres
  dbz stop my-mysql-container
`,
	Args: cobra.ExactArgs(1),
	RunE: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func runStop(cmd *cobra.Command, args []string) error {
	containerName := args[0]

	err := containers.StopContainer(containerName)
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	fmt.Printf("✅ Stopped container: %s\n", containerName)
	return nil
}
