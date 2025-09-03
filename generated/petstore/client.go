package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// APIClient represents the API client
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

// NewAPIClient creates a new APIClient instance
func NewAPIClient() *APIClient {
	return &APIClient{
		BaseURL: "https://petstore.example.com/v1",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		UserAgent: "go-api-gen/1.0.0",
	}
}

// WithBaseURL sets a custom base URL for the client
func (c *APIClient) WithBaseURL(baseURL string) *APIClient {
	c.BaseURL = baseURL
	return c
}

// WithHTTPClient sets a custom HTTP client
func (c *APIClient) WithHTTPClient(client *http.Client) *APIClient {
	c.HTTPClient = client
	return c
}

// doRequest performs an HTTP request
func (c *APIClient) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// buildURL builds a URL with path parameters and query parameters
func (c *APIClient) buildURL(path string, pathParams map[string]string, queryParams map[string]interface{}) string {
	// Replace path parameters
	for key, value := range pathParams {
		path = strings.ReplaceAll(path, "{"+key+"}", url.QueryEscape(value))
	}

	// Add query parameters
	if len(queryParams) > 0 {
		u, _ := url.Parse(c.BaseURL + path)
		q := u.Query()
		for key, value := range queryParams {
			if value != nil {
				q.Add(key, fmt.Sprintf("%v", value))
			}
		}
		u.RawQuery = q.Encode()
		return u.String()[len(c.BaseURL):]
	}

	return path
}


// Listpets List all pets
func (c *APIClient) Listpets(ctx context.Context, limit *int) error {
	pathParams := make(map[string]string)
	queryParams := make(map[string]interface{})
	
	
	
	
	if limit != nil {
		queryParams["limit"] = *limit
	}
	
	
	

	path := c.buildURL("/pets", pathParams, queryParams)
	
	return c.doRequest(ctx, "GET", path, nil, nil)
}

// Createpet Create a pet
func (c *APIClient) Createpet(ctx context.Context) error {
	pathParams := make(map[string]string)
	queryParams := make(map[string]interface{})
	
	

	path := c.buildURL("/pets", pathParams, queryParams)
	
	return c.doRequest(ctx, "POST", path, nil, nil)
}

// Getpetbyid Info for a specific pet
func (c *APIClient) Getpetbyid(ctx context.Context, petId string) error {
	pathParams := make(map[string]string)
	queryParams := make(map[string]interface{})
	
	
	
	pathParams["petId"] = fmt.Sprintf("%v", petId)
	
	

	path := c.buildURL("/pets/{petId}", pathParams, queryParams)
	
	return c.doRequest(ctx, "GET", path, nil, nil)
}

