package cmd

import (
	"admin-cli/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var portList []int
var ip string

var portReconCmd = &cobra.Command{
	Use:   "portrecon",
	Short: "Scan a range of ports on a given IP address",
	Run: func(cmd *cobra.Command, args []string) {
		// Scan the specified ports and get the list of open ports
		openPorts := internal.ScanPorts(ip, portList)

		// Display the results
		if len(openPorts) > 0 {
			fmt.Printf("Open ports on %s: %v\n", ip, openPorts)
		} else {
			fmt.Printf("No open ports found on %s\n", ip)
		}
	},
}

func init() {
	// Define the flags for ip and ports
	portReconCmd.Flags().StringVarP(&ip, "ip", "i", "127.0.0.1", "IP address to scan")
	portReconCmd.Flags().IntSliceVarP(&portList, "ports", "p", []int{80}, "Comma-separated list of ports to scan")
	rootCmd.AddCommand(portReconCmd)
}
