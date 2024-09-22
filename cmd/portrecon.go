package cmd

import (
	"admin-cli/internal"

	"github.com/spf13/cobra"
)

var port int
var ip string

var portReconCmd = &cobra.Command{
	Use:   "portrecon",
	Short: "Scan a port",
	Run: func(cmd *cobra.Command, args []string) {
		internal.ScanPort(ip, port)
	},
}

func init() {
	portReconCmd.Flags().IntVarP(&port, "port", "p", 80, "Port to scan")
	portReconCmd.Flags().StringVarP(&ip, "ip", "i", "localhost", "IP address to scan")
	rootCmd.AddCommand(portReconCmd)
}
