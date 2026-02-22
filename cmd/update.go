package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dbz to the latest version",
	Long: `Update dbz by rebuilding from source and reinstalling.

This command will:
  1. Build the dbz binary from the current source code
  2. Install it to /usr/local/bin (or equivalent on your system)

You must run this command from the dbz source directory.`,
	Example: `  dbz update
  dbz update --from /path/to/dbz`,
	RunE: runUpdate,
}

var updateFromPath string

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&updateFromPath, "from", "f", "", "Path to dbz source directory (default: current directory)")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// Determine the source directory
	sourceDir := updateFromPath
	if sourceDir == "" {
		// Try to find the source directory from the current executable
		exe, err := os.Executable()
		if err == nil {
			// If the executable is in /usr/local/bin, check if there's a source directory nearby
			if filepath.Dir(exe) == "/usr/local/bin" {
				// Check common locations
				possibleDirs := []string{
					filepath.Join(os.Getenv("HOME"), "Development", "dbz"),
					filepath.Join(os.Getenv("HOME"), "Projects", "dbz"),
					filepath.Join(os.Getenv("HOME"), "Code", "dbz"),
					"/opt/dbz",
				}
				for _, dir := range possibleDirs {
					if _, err := os.Stat(filepath.Join(dir, "main.go")); err == nil {
						sourceDir = dir
						break
					}
				}
			}
		}
		if sourceDir == "" {
			// Use current working directory
			var err error
			sourceDir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
		}
	}

	// Verify this is a dbz source directory
	mainGoPath := filepath.Join(sourceDir, "main.go")
	if _, err := os.Stat(mainGoPath); os.IsNotExist(err) {
		return fmt.Errorf("dbz source directory not found at %s (missing main.go)", sourceDir)
	}

	fmt.Printf("Updating dbz from %s...\n", sourceDir)

	// Determine install path
	installPath := "/usr/local/bin"
	if runtime.GOOS == "windows" {
		installPath = filepath.Join(os.Getenv("ProgramFiles"), "dbz")
	}

	// Build the binary
	fmt.Println("Building dbz...")
	buildCmd := exec.Command("go", "build", "-o", "dbz", "main.go")
	buildCmd.Dir = sourceDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build dbz: %w", err)
	}

	// Install the binary
	fmt.Printf("Installing dbz to %s...\n", installPath)
	binaryPath := filepath.Join(sourceDir, "dbz")
	targetPath := filepath.Join(installPath, "dbz")

	// On Unix systems, we need sudo for /usr/local/bin
	if runtime.GOOS != "windows" && installPath == "/usr/local/bin" {
		cpCmd := exec.Command("sudo", "cp", binaryPath, targetPath)
		cpCmd.Stdout = os.Stdout
		cpCmd.Stderr = os.Stderr
		if err := cpCmd.Run(); err != nil {
			return fmt.Errorf("failed to install dbz: %w", err)
		}

		// Make executable
		chmodCmd := exec.Command("sudo", "chmod", "+x", targetPath)
		chmodCmd.Stdout = os.Stdout
		chmodCmd.Stderr = os.Stderr
		if err := chmodCmd.Run(); err != nil {
			return fmt.Errorf("failed to set permissions: %w", err)
		}
	} else {
		// Direct copy for other cases
		data, err := os.ReadFile(binaryPath)
		if err != nil {
			return fmt.Errorf("failed to read built binary: %w", err)
		}
		if err := os.WriteFile(targetPath, data, 0755); err != nil {
			return fmt.Errorf("failed to install dbz: %w", err)
		}
	}

	fmt.Println("✅ dbz updated successfully!")
	fmt.Printf("Run 'dbz --version' to verify the update.\n")
	return nil
}
