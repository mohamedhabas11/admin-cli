package cmd

import (
	"admin-cli/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	input          string
	output         string
	compLevel      int
	followSymlinks bool
	passphrase     string
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Archive and compress a path",
	Run: func(cmd *cobra.Command, args []string) {
		if err := internal.Backup(input, output, compLevel, followSymlinks, passphrase); err != nil {
			fmt.Printf("Backup failed: %v\n", err)
		} else {
			fmt.Println("Backup completed successfully")
		}
	},
}

func init() {
	backupCmd.Flags().StringVarP(&input, "input", "i", ".", "Backup input path")
	backupCmd.Flags().StringVarP(&output, "output", "o", "./backup.tar.zst", "Backup output path")
	backupCmd.Flags().IntVarP(&compLevel, "compression-level", "c", 3, "Compression level (higher means better compression, slower speed)")
	backupCmd.Flags().BoolVarP(&followSymlinks, "follow-symlinks", "f", false, "Follow symbolic links when archiving")
	backupCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Age recipient passphrase")
	rootCmd.AddCommand(backupCmd)
}
