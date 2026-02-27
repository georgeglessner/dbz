package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/dbz/dbz/pkg/containers"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all database containers",
	Long:    `List all running database containers created by dbz.`,
	Aliases: []string{"ls"},
	RunE:    runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	dbContainers, err := containers.ListContainers()
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	if len(dbContainers) == 0 {
		fmt.Println("No database containers found.")
		return nil
	}

	// Create a tabwriter for formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "NAME\tTYPE\tVERSION\tSTATUS\tPORT\tUSER\tDATABASE")
	_, _ = fmt.Fprintln(w, "----\t----\t-------\t------\t----\t----\t--------")

	for _, container := range dbContainers {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\t%s\n",
			container.Name,
			container.Type,
			container.Version,
			container.Status,
			container.Port,
			container.User,
			container.Database,
		)
	}

	_ = w.Flush()
	return nil
}
