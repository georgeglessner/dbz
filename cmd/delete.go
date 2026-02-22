package cmd

import (
	"fmt"

	"github.com/dbz/dbz/pkg/containers"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [database]",
	Short: "Delete a database container",
	Long: `Delete a database container by name or ID.

Examples:
  dbz delete postgres
  dbz delete my-mysql-container
`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
	containerName := args[0]

	err := containers.DeleteContainer(containerName)
	if err != nil {
		return fmt.Errorf("failed to delete container: %w", err)
	}

	fmt.Printf("✅ Deleted container: %s\n", containerName)
	return nil
}
