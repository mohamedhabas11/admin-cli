package cmd

import (
	"admin-cli/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var length int

var genPassCmd = &cobra.Command{
	Use:   "genpass",
	Short: "Generate a password",
	Run: func(cmd *cobra.Command, args []string) {
		pwd, _ := internal.GeneratePassword(length)
		fmt.Println("Generated Password:", pwd)
	},
}

func init() {
	genPassCmd.Flags().IntVarP(&length, "length", "l", 12, "Length of the password")
	rootCmd.AddCommand(genPassCmd)
}
