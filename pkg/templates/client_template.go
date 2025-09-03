package templates

const clientTemplate = `package {{.PackageName}}

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

// {{.ClientName}} represents the API client
type {{.ClientName}} struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

// New{{.ClientName}} creates a new {{.ClientName}} instance
func New{{.ClientName}}() *{{.ClientName}} {
	return &{{.ClientName}}{
		BaseURL: "{{.BaseURL}}",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		UserAgent: "go-api-gen/1.0.0",
	}
}

// WithBaseURL sets a custom base URL for the client
func (c *{{.ClientName}}) WithBaseURL(baseURL string) *{{.ClientName}} {
	c.BaseURL = baseURL
	return c
}

// WithHTTPClient sets a custom HTTP client
func (c *{{.ClientName}}) WithHTTPClient(client *http.Client) *{{.ClientName}} {
	c.HTTPClient = client
	return c
}

// doRequest performs an HTTP request
func (c *{{.ClientName}}) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
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
func (c *{{.ClientName}}) buildURL(path string, pathParams map[string]string, queryParams map[string]interface{}) string {
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

{{range .Operations}}
// {{.Name}} {{.Summary}}
func (c *{{$.ClientName}}) {{.Name}}(ctx context.Context{{range .Parameters}}{{if .Required}}, {{.Name}} {{.Type}}{{end}}{{end}}{{range .Parameters}}{{if not .Required}}, {{.Name}} *{{.Type}}{{end}}{{end}}) error {
	pathParams := make(map[string]string)
	queryParams := make(map[string]interface{})
	
	{{range .Parameters}}
	{{if eq .In "path"}}
	pathParams["{{.Name}}"] = fmt.Sprintf("%v", {{.Name}})
	{{else if eq .In "query"}}
	{{if .Required}}
	queryParams["{{.Name}}"] = {{.Name}}
	{{else}}
	if {{.Name}} != nil {
		queryParams["{{.Name}}"] = *{{.Name}}
	}
	{{end}}
	{{end}}
	{{end}}

	path := c.buildURL("{{.Path}}", pathParams, queryParams)
	
	return c.doRequest(ctx, "{{.Method}}", path, nil, nil)
}
{{end}}
`