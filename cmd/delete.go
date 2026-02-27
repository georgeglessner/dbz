package cmd

import (
	"fmt"

	"github.com/dbz/dbz/pkg/containers"
	"github.com/spf13/cobra"
)

var deleteRemoveVolumes bool

var deleteCmd = &cobra.Command{
	Use:   "delete [database]",
	Short: "Delete a database container",
	Long: `Delete a database container by name or ID.

Examples:
  dbz delete postgres
  dbz delete my-mysql-container
  dbz delete postgres --volume    # Delete container and associated volumes`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&deleteRemoveVolumes, "volume", "v", false, "Also remove volumes associated with the container")
}

func runDelete(cmd *cobra.Command, args []string) error {
	containerName := args[0]

	err := containers.DeleteContainer(containerName, deleteRemoveVolumes)
	if err != nil {
		return fmt.Errorf("failed to delete container: %w", err)
	}

	if deleteRemoveVolumes {
		fmt.Printf("✅ Deleted container and volumes: %s\n", containerName)
	} else {
		fmt.Printf("✅ Deleted container: %s\n", containerName)
	}
	return nil
}
