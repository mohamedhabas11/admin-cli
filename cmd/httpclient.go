package cmd

import (
	"admin-cli/http"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
)

var (
	baseURL    string
	headers    []string
	authUser   string
	authPass   string
	token      string
	timeout    int
	retryCount int
	proxyURL   string
)

// httpClientCmd is the root command for the HTTP client utility
var httpClientCmd = &cobra.Command{
	Use:   "httpclient",
	Short: "HTTP client utility for making requests",
}

// getCmd performs a GET request
var getCmd = &cobra.Command{
	Use:   "get [path]",
	Short: "Perform a GET request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := setupClient()
		resp, err := client.Get(args[0])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
	},
}

// uploadCmd uploads a file
var uploadCmd = &cobra.Command{
	Use:   "upload [path] [file]",
	Short: "Upload a file to the specified path",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := setupClient()
		resp, err := client.UploadFile(args[0], args[1])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Upload Response:", string(body))
	},
}

// downloadCmd downloads a file
var downloadCmd = &cobra.Command{
	Use:   "download [path] [save-path]",
	Short: "Download a file and save it locally",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := setupClient()
		err := client.DownloadFile(args[0], args[1])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("File downloaded successfully to", args[1])
	},
}

// setupClient configures the HTTP client based on flags
func setupClient() *http.Client {
	client := http.NewClient(baseURL)
	client.Timeout = time.Duration(timeout) * time.Second
	client.RetryCount = retryCount
	client.ProxyURL = proxyURL

	for _, h := range headers {
		var key, value string
		fmt.Sscanf(h, "%s:%s", &key, &value)
		client.SetHeader(key, value)
	}

	if authUser != "" && authPass != "" {
		client.SetAuth(authUser, authPass)
	} else if token != "" {
		client.SetToken(token)
	}

	return client
}

func init() {
	// Add subcommands to the root command
	httpClientCmd.AddCommand(getCmd, uploadCmd, downloadCmd)

	// Define persistent flags
	httpClientCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "Base URL for requests")
	httpClientCmd.PersistentFlags().StringSliceVar(&headers, "header", []string{}, "Custom headers (e.g., 'Key:Value')")
	httpClientCmd.PersistentFlags().StringVar(&authUser, "user", "", "Username for basic auth")
	httpClientCmd.PersistentFlags().StringVar(&authPass, "pass", "", "Password for basic auth")
	httpClientCmd.PersistentFlags().StringVar(&token, "token", "", "Token for bearer auth")
	httpClientCmd.PersistentFlags().IntVar(&timeout, "timeout", 30, "Request timeout in seconds")
	httpClientCmd.PersistentFlags().IntVar(&retryCount, "retry", 0, "Number of retries")
	httpClientCmd.PersistentFlags().StringVar(&proxyURL, "proxy", "", "Proxy URL")

	// Add the httpclient command to the root command (assumes a rootCmd exists)
	rootCmd.AddCommand(httpClientCmd)
}
