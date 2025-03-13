package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// AuthConfig holds authentication details
type AuthConfig struct {
	Username string
	Password string
	Token    string
}

// Client represents the HTTP client with configuration
type Client struct {
	BaseURL    string
	Headers    map[string]string
	Auth       AuthConfig
	Timeout    time.Duration
	RetryCount int
	ProxyURL   string
}

// NewClient creates a new HTTP client with the given base URL
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Headers:    make(map[string]string),
		Timeout:    30 * time.Second,
		RetryCount: 0,
	}
}

// SetHeader sets a custom header for the client
func (c *Client) SetHeader(key, value string) {
	c.Headers[key] = value
}

// SetAuth sets basic authentication credentials
func (c *Client) SetAuth(username, password string) {
	c.Auth.Username = username
	c.Auth.Password = password
}

// SetToken sets a token for bearer authentication
func (c *Client) SetToken(token string) {
	c.Auth.Token = token
}

// Get performs a GET request to the specified path
func (c *Client) Get(path string) (*http.Response, error) {
	return c.doRequest("GET", path, nil)
}

// Post performs a POST request with the given body
func (c *Client) Post(path string, body interface{}) (*http.Response, error) {
	return c.doRequest("POST", path, body)
}

// Put performs a PUT request with the given body
func (c *Client) Put(path string, body interface{}) (*http.Response, error) {
	return c.doRequest("PUT", path, body)
}

// Delete performs a DELETE request to the specified path
func (c *Client) Delete(path string) (*http.Response, error) {
	return c.doRequest("DELETE", path, nil)
}

// UploadFile uploads a file to the specified path
func (c *Client) UploadFile(path, filePath string) (*http.Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", c.BaseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: c.Timeout}
	return client.Do(req)
}

// DownloadFile downloads a file from the specified path and saves it locally
func (c *Client) DownloadFile(path, savePath string) error {
	resp, err := c.Get(path)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}
	return nil
}

// doRequest is a helper function to perform HTTP requests
func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.Auth.Username != "" && c.Auth.Password != "" {
		req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	} else if c.Auth.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Auth.Token)
	}

	client := &http.Client{Timeout: c.Timeout}
	if c.ProxyURL != "" {
		proxyURL, err := url.Parse(c.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %v", err)
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	var resp *http.Response
	for i := 0; i <= c.RetryCount; i++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	if err != nil {
		return nil, fmt.Errorf("request failed after retries: %v", err)
	}
	return resp, nil
}
