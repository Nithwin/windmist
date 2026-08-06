package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Nithwin/WindMist/internal/config"
	"github.com/spf13/cobra"
)

var (
	uninstallYesFlag   bool
	uninstallPurgeFlag bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the WindMist CLI executable and optionally purge configuration",
	Long:  `Safely removes the WindMist CLI executable from your system binary path and optionally cleans up saved configuration and data (` + "`~/.config/windmist`" + ` and ` + "`~/.windmist`" + `).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to locate windmist executable: %w", err)
		}

		binaryPath, err := filepath.EvalSymlinks(execPath)
		if err != nil {
			binaryPath = execPath
		}

		reader := bufio.NewReader(os.Stdin)

		if !uninstallYesFlag {
			fmt.Printf("⚠️  Are you sure you want to uninstall WindMist from %s? [y/N]: ", binaryPath)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if !strings.EqualFold(input, "y") && !strings.EqualFold(input, "yes") {
				fmt.Println("Uninstallation cancelled.")
				return nil
			}
		}

		if !uninstallPurgeFlag && !uninstallYesFlag {
			cfgDir, _ := config.ConfigDir()
			home, _ := os.UserHomeDir()
			dataDir := ""
			if home != "" {
				dataDir = filepath.Join(home, ".windmist")
			}
			purgeTarget := cfgDir
			if dataDir != "" {
				if purgeTarget != "" {
					purgeTarget = purgeTarget + " and " + dataDir
				} else {
					purgeTarget = dataDir
				}
			}
			if purgeTarget != "" {
				cfgExists := false
				if cfgDir != "" {
					if _, err := os.Stat(cfgDir); err == nil {
						cfgExists = true
					}
				}
				dataExists := false
				if dataDir != "" {
					if _, err := os.Stat(dataDir); err == nil {
						dataExists = true
					}
				}
				if cfgExists || dataExists {
					fmt.Printf("❓ Also remove all configuration and chat history in %s? [y/N]: ", purgeTarget)
					input, _ := reader.ReadString('\n')
					input = strings.TrimSpace(input)
					if strings.EqualFold(input, "y") || strings.EqualFold(input, "yes") {
						uninstallPurgeFlag = true
					}
				}
			}
		}

		fmt.Printf("🗑️  Removing binary at %s...\n", binaryPath)
		err = os.Remove(binaryPath)
		if err != nil && os.IsPermission(err) {
			if runtime.GOOS != "windows" {
				fmt.Println("🔑 Permission denied. Requesting root permissions via sudo...")
				sudoCmd := exec.Command("sudo", "rm", "-f", binaryPath)
				sudoCmd.Stdin = os.Stdin
				sudoCmd.Stdout = os.Stdout
				sudoCmd.Stderr = os.Stderr
				if err := sudoCmd.Run(); err != nil {
					return fmt.Errorf("failed to remove binary via sudo: %w", err)
				}
			} else {
				return fmt.Errorf("permission denied. Please run this command as Administrator to uninstall")
			}
		} else if err != nil {
			return fmt.Errorf("failed to remove binary: %w", err)
		}

		if uninstallPurgeFlag {
			cfgDir, _ := config.ConfigDir()
			if cfgDir != "" {
				if _, err := os.Stat(cfgDir); err == nil {
					fmt.Printf("🗑️  Removing configuration directory at %s...\n", cfgDir)
					if err := os.RemoveAll(cfgDir); err != nil {
						fmt.Printf("⚠️  Warning: failed to remove configuration directory: %v\n", err)
					}
				}
			}
			if home, err := os.UserHomeDir(); err == nil {
				dataDir := filepath.Join(home, ".windmist")
				if _, err := os.Stat(dataDir); err == nil {
					fmt.Printf("🗑️  Removing data directory at %s...\n", dataDir)
					if err := os.RemoveAll(dataDir); err != nil {
						fmt.Printf("⚠️  Warning: failed to remove data directory: %v\n", err)
					}
				}
			}
		}

		fmt.Println("\n✨ WindMist has been successfully uninstalled. Goodbye!")
		return nil
	},
}

func init() {
	uninstallCmd.Flags().BoolVarP(&uninstallYesFlag, "yes", "y", false, "Skip confirmation prompts and uninstall immediately")
	uninstallCmd.Flags().BoolVarP(&uninstallPurgeFlag, "purge", "p", false, "Also remove user configuration and chat history directory")
	rootCmd.AddCommand(uninstallCmd)
}
