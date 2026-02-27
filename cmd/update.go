package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const githubRepo = "georgeglessner/dbz"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dbz to the latest version",
	Long: `Update dbz to the latest version.

Modes:
  Local (default): Rebuild from source code
    dbz update
    dbz update --from /path/to/source
  
  Remote: Download latest release from GitHub
    dbz update --remote`,
	Example: `  dbz update
  dbz update --from /path/to/dbz
  dbz update --remote`,
	RunE: runUpdate,
}

var (
	updateFromPath string
	remoteUpdate   bool
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&updateFromPath, "from", "f", "", "Path to dbz source directory (for local mode)")
	updateCmd.Flags().BoolVar(&remoteUpdate, "remote", false, "Download and install from GitHub releases")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	if remoteUpdate {
		return runRemoteUpdate()
	}
	return runLocalUpdate()
}

// runLocalUpdate builds dbz from local source code
func runLocalUpdate() error {
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
	binaryPath := filepath.Join(sourceDir, "dbz")
	if err := installBinary(binaryPath, installPath); err != nil {
		return err
	}

	fmt.Println("✅ dbz updated successfully!")
	fmt.Printf("Run 'dbz --version' to verify the update.\n")
	return nil
}

// runRemoteUpdate downloads and installs the latest release from GitHub
func runRemoteUpdate() error {
	fmt.Println("Updating dbz from GitHub releases...")

	// Detect platform
	platform := detectPlatform()
	fmt.Printf("Detected platform: %s\n", platform)

	// Fetch latest version
	version, err := getLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to fetch latest version: %w", err)
	}
	fmt.Printf("Latest version: %s\n", version)

	// Determine install path
	installPath := "/usr/local/bin"
	if runtime.GOOS == "windows" {
		installPath = filepath.Join(os.Getenv("ProgramFiles"), "dbz")
	}

	// Check current version
	currentVersion := "v" + rootCmd.Version
	if currentVersion == version {
		fmt.Println("You already have the latest version installed.")
		return nil
	}

	// Confirm before updating
	fmt.Printf("\nThis will update dbz from %s to %s\n", currentVersion, version)
	fmt.Printf("Install location: %s\n", installPath)
	fmt.Print("Continue? [y/N]: ")

	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		// If we can't read input (e.g., EOF), default to cancelling
		fmt.Println("\nUpdate cancelled.")
		return nil
	}
	response = strings.ToLower(strings.TrimSpace(response))

	if response != "y" && response != "yes" {
		fmt.Println("Update cancelled.")
		return nil
	}

	// Download binary
	fmt.Println("\nDownloading...")
	binaryPath, err := downloadBinary(version, platform)
	if err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	defer func() { _ = os.Remove(binaryPath) }()

	// Install the binary
	if err := installBinary(binaryPath, installPath); err != nil {
		return err
	}

	fmt.Println("✅ dbz updated successfully!")
	fmt.Printf("Run 'dbz --version' to verify the update.\n")
	return nil
}

// detectPlatform returns the platform identifier for the current system
func detectPlatform() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	// Normalize architecture names
	switch arch {
	case "amd64":
		arch = "amd64"
	case "arm64":
		arch = "arm64"
	}

	return fmt.Sprintf("%s-%s", os, arch)
}

// getLatestVersion fetches the latest release version from GitHub
func getLatestVersion() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", githubRepo)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

// downloadBinary downloads the binary for the specified version and platform
func downloadBinary(version, platform string) (string, error) {
	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/dbz-%s", githubRepo, version, platform)

	tempFile, err := os.CreateTemp("", "dbz-update-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		_ = tempFile.Close()
	}()

	resp, err := http.Get(downloadURL)
	if err != nil {
		_ = os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		_ = os.Remove(tempFile.Name())
		return "", fmt.Errorf("download failed with status %d (binary may not exist for %s)", resp.StatusCode, platform)
	}

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		_ = os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to save binary: %w", err)
	}

	// Make executable
	if err := os.Chmod(tempFile.Name(), 0755); err != nil {
		_ = os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to make binary executable: %w", err)
	}

	return tempFile.Name(), nil
}

// installBinary installs the binary to the specified path
func installBinary(sourcePath, installPath string) error {
	targetPath := filepath.Join(installPath, "dbz")

	fmt.Printf("Installing dbz to %s...\n", installPath)

	// On Unix systems, we need sudo for /usr/local/bin
	if runtime.GOOS != "windows" && installPath == "/usr/local/bin" {
		cpCmd := exec.Command("sudo", "cp", sourcePath, targetPath)
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
		data, err := os.ReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("failed to read binary: %w", err)
		}
		if err := os.WriteFile(targetPath, data, 0755); err != nil {
			return fmt.Errorf("failed to install dbz: %w", err)
		}
	}

	fmt.Printf("✅ Successfully installed dbz to %s\n", installPath)
	return nil
}
