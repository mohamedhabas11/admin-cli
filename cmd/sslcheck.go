package cmd

import (
	"admin-cli/internal"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var sslHost string
var sslPort int

var sslCheckCmd = &cobra.Command{
	Use:   "sslcheck",
	Short: "Check SSL certificate expiration",
	Run: func(cmd *cobra.Command, args []string) {
		expiry, err := internal.CheckSSLExpiry(sslHost, sslPort)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		daysRemaining := int(time.Until(expiry).Hours() / 24)
		fmt.Printf("Certificate for %s:%d expires on: %s (%d days remaining)\n", sslHost, sslPort, expiry.Format("2006-01-02"), daysRemaining)
	},
}

func init() {
	sslCheckCmd.Flags().StringVar(&sslHost, "host", "", "Host to check")
	sslCheckCmd.Flags().IntVar(&sslPort, "port", 443, "Port to check")
	sslCheckCmd.MarkFlagRequired("host")
	rootCmd.AddCommand(sslCheckCmd)
}
