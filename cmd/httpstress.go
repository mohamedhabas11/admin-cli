package cmd

import (
	"admin-cli/internal"

	"github.com/spf13/cobra"
)

var (
	url         string
	numRequests int
	concurrency int
)

var httpStressCmd = &cobra.Command{
	Use:   "httpstress",
	Short: "Perform HTTP stress testing on a server",
	Run: func(cmd *cobra.Command, args []string) {
		// Call the function from the internal package
		internal.PerformStressTest(url, numRequests, concurrency)
	},
}

func init() {
	// Flags for HTTP stress testing
	httpStressCmd.Flags().StringVarP(&url, "url", "u", "", "URL to stress test (required)")
	httpStressCmd.Flags().IntVarP(&numRequests, "requests", "r", 100, "Number of HTTP requests to send")
	httpStressCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 10, "Number of concurrent requests")
	httpStressCmd.MarkFlagRequired("url") // URL is required

	// Add the command to rootCmd
	rootCmd.AddCommand(httpStressCmd)
}
