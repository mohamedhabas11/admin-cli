package cmd

import (
	"admin-cli/internal"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var portList []int
var ip string
var portRange string
var allFlag bool // Flag to indicate --all option

var portReconCmd = &cobra.Command{
	Use:   "portrecon",
	Short: "Scan a range of ports on a given IP address",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse the port range if provided
		if portRange != "" {
			start, end, err := parsePortRange(portRange)
			if err != nil {
				fmt.Println("Invalid port range:", err)
				return
			}
			// Append the range of ports to the portList
			for i := start; i <= end; i++ {
				portList = append(portList, i)
			}
		}

		// Remove duplicate ports
		portList = removeDuplicates(portList)

		// Scan the specified ports and get the list of open ports
		openPorts := internal.ScanPorts(ip, portList, allFlag || portRange != "")

		// Display the results
		if len(openPorts) > 0 {
			fmt.Printf("Open ports on %s: %v\n", ip, openPorts)
		} else {
			fmt.Printf("No open ports found on %s\n", ip)
		}
	},
}

func init() {
	// Define the flags for ip, ports, and port range
	portReconCmd.Flags().StringVarP(&ip, "ip", "i", "127.0.0.1", "IP address to scan")
	portReconCmd.Flags().IntSliceVarP(&portList, "ports", "p", []int{80}, "Comma-separated list of ports to scan")
	portReconCmd.Flags().StringVarP(&portRange, "range", "r", "", "Range of ports to scan (e.g., 20-80)")
	portReconCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Scan all ports (1-65535)")
	rootCmd.AddCommand(portReconCmd)
}

// parsePortRange converts a port range string (e.g., "20-80") to start and end ports
func parsePortRange(rangeStr string) (int, int, error) {
	ports := strings.Split(rangeStr, "-")
	if len(ports) != 2 {
		return 0, 0, fmt.Errorf("invalid range format")
	}
	start, err1 := strconv.Atoi(ports[0])
	end, err2 := strconv.Atoi(ports[1])
	if err1 != nil || err2 != nil || start < 1 || end < 1 || start > 65535 || end > 65535 || start > end {
		return 0, 0, fmt.Errorf("invalid port range values")
	}
	return start, end, nil
}

// removeDuplicates removes duplicate ports from the port list
func removeDuplicates(ports []int) []int {
	uniquePorts := make(map[int]struct{})
	for _, port := range ports {
		uniquePorts[port] = struct{}{}
	}

	var result []int
	for port := range uniquePorts {
		result = append(result, port)
	}
	return result
}
