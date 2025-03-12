package cmd

import (
	"admin-cli/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	restoreInput      string
	restoreOutput     string
	restorePassphrase string
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore backup from a tar.zst archive",
	Run: func(cmd *cobra.Command, args []string) {
		if err := internal.Restore(restoreInput, restoreOutput, restorePassphrase); err != nil {
			fmt.Printf("Restore failed: %v\n", err)
		} else {
			fmt.Println("Restore completed successfully")
		}
	},
}

func init() {
	restoreCmd.Flags().StringVarP(&restoreInput, "input", "i", "./backup.tar.zst", "Path to backup tar.zst file")
	restoreCmd.Flags().StringVarP(&restoreOutput, "output", "o", ".", "Destination path to restore the backup")
	restoreCmd.Flags().StringVarP(&restorePassphrase, "passphrase", "p", "", "Age recipient passphrase")
	rootCmd.AddCommand(restoreCmd)
}
