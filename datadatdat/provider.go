// Package datadatdat provides a remote provider for connecting d3 CLI to datadatdat-remote-server.
// This provider implements the remote.Remote interface and communicates with the datadatdat-remote-server
// HTTP APIs to store and retrieve commits.
package datadatdat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/datadatdat/remote-sdk-go/remote"
)

const (
	propAPIBaseURL = "api_base_url"
	propOrg        = "org"
	propRepo       = "repo"
	propAPIToken   = "api_token"
	propPort       = "port"
)

// Provider implements the remote.Remote interface for datadatdat-remote-server.
type Provider struct {
	httpClient *http.Client
}

// getHTTPClient returns the HTTP client, initializing it if necessary
func (p *Provider) getHTTPClient() *http.Client {
	if p.httpClient == nil {
		p.httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return p.httpClient
}

// commitResponse represents the API response for a single commit
type commitResponse struct {
	Repo      string                 `json:"repo"`
	CommitID  string                 `json:"commitId"`
	Timestamp time.Time              `json:"timestamp"`
	Size      int64                  `json:"size"`
	Author    string                 `json:"author,omitempty"`
	Message   string                 `json:"message,omitempty"`
	ParentID  string                 `json:"parentId,omitempty"`
	Tags      []string               `json:"tags,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// listCommitsResponse represents the API response for listing commits
type listCommitsResponse struct {
	Repo       string           `json:"repo"`
	Commits    []commitResponse `json:"commits"`
	NextCursor string           `json:"nextCursor,omitempty"`
}

// NewProvider creates a new datadatdat remote provider instance.
func NewProvider() *Provider {
	return &Provider{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// buildAPIURL constructs the full API URL for a given path
func (p *Provider) buildAPIURL(properties map[string]interface{}, path string) (string, error) {
	apiBaseURL, ok := properties[propAPIBaseURL].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid api_base_url property")
	}

	org, ok := properties[propOrg].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid org property")
	}

	repo, ok := properties[propRepo].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid repo property")
	}

	// Build the full URL
	fullURL := fmt.Sprintf("%s/api/v1/repos/%s/%s%s", apiBaseURL, org, repo, path)
	return fullURL, nil
}

// doRequest performs an HTTP request with optional authentication
func (p *Provider) doRequest(ctx context.Context, method, url string, body io.Reader, parameters map[string]interface{}) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add Bearer token if present
	if apiToken, ok := parameters[propAPIToken].(string); ok && apiToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	}

	resp, err := p.getHTTPClient().Do(req) // #nosec G704 -- URL is constructed from user-configured remote properties; SSRF is by design for a remote provider
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return resp, nil
}

// parseCommitResponse converts API commit response to remote.Commit
func parseCommitResponse(cr commitResponse) remote.Commit {
	properties := make(map[string]interface{})

	// Add standard properties
	properties["timestamp"] = cr.Timestamp
	properties["size"] = cr.Size

	if cr.Author != "" {
		properties["author"] = cr.Author
	}
	if cr.Message != "" {
		properties["message"] = cr.Message
	}
	if cr.ParentID != "" {
		properties["parentId"] = cr.ParentID
	}
	if len(cr.Tags) > 0 {
		// Convert tags array to map for compatibility
		tagsMap := make(map[string]string)
		for _, tag := range cr.Tags {
			parts := strings.SplitN(tag, ":", 2)
			if len(parts) == 2 {
				tagsMap[parts[0]] = parts[1]
			} else {
				tagsMap[tag] = ""
			}
		}
		properties["tags"] = tagsMap
	}

	// Merge any additional metadata
	for k, v := range cr.Metadata {
		if _, exists := properties[k]; !exists {
			properties[k] = v
		}
	}

	return remote.Commit{
		ID:         cr.CommitID,
		Properties: properties,
	}
}

// Type returns the remote type identifier for the datadatdat remote.
// Returns "datadatdat" to match the Kotlin provider name
func (p *Provider) Type() (string, error) {
	return "datadatdat", nil
}

// FromURL parses a datadatdat remote URL and additional properties to create remote properties.
// Expected URL format: http://hostname:port/org/repo or https://hostname/org/repo
func (p *Provider) FromURL(rawURL string, additionalProperties map[string]string) (map[string]interface{}, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Validate scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme %s: must be http or https", parsedURL.Scheme)
	}

	// Validate host is present
	if parsedURL.Host == "" {
		return nil, fmt.Errorf("missing host in URL")
	}

	// Extract org/repo from path
	path := strings.Trim(parsedURL.Path, "/")
	if path == "" {
		return nil, fmt.Errorf("invalid path: expected /org/repo format")
	}
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		return nil, fmt.Errorf("invalid path: expected /org/repo format")
	}
	if len(pathParts) > 2 {
		return nil, fmt.Errorf("invalid path: too many segments, expected /org/repo format")
	}

	org := pathParts[0]
	repo := pathParts[1]

	if org == "" || repo == "" {
		return nil, fmt.Errorf("invalid path: org and repo cannot be empty")
	}

	// Build API base URL (scheme://host:port)
	apiBaseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	portStr := parsedURL.Port()

	properties := map[string]interface{}{
		propAPIBaseURL: apiBaseURL,
		propOrg:          org,
		propRepo:         repo,
	}

	// Add port if present
	if portStr != "" {
		port := 0
		_, err := fmt.Sscanf(portStr, "%d", &port)
		if err != nil {
			return nil, fmt.Errorf("invalid port: %s", portStr)
		}
		if port < 1 || port > 65535 {
			return nil, fmt.Errorf("invalid port: %d (must be between 1 and 65535)", port)
		}
		properties[propPort] = port
	}

	// Validate and add additional properties (only api_token is allowed)
	allowedProps := map[string]bool{
		propAPIToken: true,
	}
	for k, v := range additionalProperties {
		if !allowedProps[k] {
			return nil, fmt.Errorf("unknown additional property: %s", k)
		}
		properties[k] = v
	}

	return properties, nil
}

// ToURL converts remote properties back to a URL and additional properties.
func (p *Provider) ToURL(properties map[string]interface{}) (string, map[string]string, error) {
	apiBaseURL, ok := properties[propAPIBaseURL].(string)
	if !ok {
		return "", nil, fmt.Errorf("missing api_base_url property")
	}

	org, ok := properties[propOrg].(string)
	if !ok {
		return "", nil, fmt.Errorf("missing org property")
	}

	repo, ok := properties[propRepo].(string)
	if !ok {
		return "", nil, fmt.Errorf("missing repo property")
	}

	// Construct the full URL
	fullURL := fmt.Sprintf("%s/%s/%s", apiBaseURL, org, repo)

	// Extract additional properties (excluding system properties)
	additionalProps := make(map[string]string)
	systemProps := map[string]bool{
		propAPIBaseURL: true,
		propOrg:          true,
		propRepo:         true,
		propPort:         true,
	}

	for k, v := range properties {
		if !systemProps[k] {
			if strVal, ok := v.(string); ok {
				additionalProps[k] = strVal
			}
		}
	}

	return fullURL, additionalProps, nil
}

// GetParameters extracts operation parameters from remote properties.
// This is called before each operation to get parameters like API tokens.
// If api_token is not present in properties, it will be read from the DATADATDAT_API_KEY environment variable.
func (p *Provider) GetParameters(remoteProperties map[string]interface{}) (map[string]interface{}, error) {
	// Create a copy of the properties to avoid modifying the original
	params := make(map[string]interface{})
	for k, v := range remoteProperties {
		params[k] = v
	}

	// If api_token is not already present, try to read from environment variable
	if _, hasToken := params[propAPIToken]; !hasToken {
		if envToken := os.Getenv("DATADATDAT_API_KEY"); envToken != "" {
			params[propAPIToken] = envToken
		}
	}

	return params, nil
}

// ValidateRemote validates the remote connection properties.
func (p *Provider) ValidateRemote(properties map[string]interface{}) error {
	requiredProps := []string{propAPIBaseURL, propOrg, propRepo}
	for _, prop := range requiredProps {
		if _, ok := properties[prop]; !ok {
			return fmt.Errorf("missing required property: %s", prop)
		}
	}

	// Validate allowed properties (must match what the Kotlin server accepts)
	allowedProps := map[string]bool{
		propAPIBaseURL: true,
		propOrg:          true,
		propRepo:         true,
		propPort:         true,
		propAPIToken:    true,
	}

	for k := range properties {
		if !allowedProps[k] {
			return fmt.Errorf("unknown property: %s", k)
		}
	}

	// Validate port if present
	if port, ok := properties[propPort]; ok {
		portInt := 0
		switch v := port.(type) {
		case int:
			portInt = v
		case float64:
			portInt = int(v)
		case float32:
			portInt = int(v)
		default:
			return fmt.Errorf("invalid port type: must be integer")
		}

		if portInt < 1 || portInt > 65535 {
			return fmt.Errorf("invalid port: %d (must be between 1 and 65535)", portInt)
		}
	}

	return nil
}

// ValidateParameters validates the operation parameters.
func (p *Provider) ValidateParameters(parameters map[string]interface{}) error {
	// Validate allowed parameters (must match what the Kotlin server accepts)
	allowedParams := map[string]bool{
		propAPIToken:    true,
		propAPIBaseURL: true,
		propOrg:          true,
		propRepo:         true,
		propPort:         true,
	}

	for k := range parameters {
		if !allowedParams[k] {
			return fmt.Errorf("unknown parameter: %s", k)
		}
	}

	return nil
}

// ListCommits returns a list of available commits from the remote server, optionally filtered by tags.
func (p *Provider) ListCommits(properties map[string]interface{}, parameters map[string]interface{}, tags []remote.Tag) ([]remote.Commit, error) {
	// Build the API URL
	apiURL, err := p.buildAPIURL(properties, "/commits")
	if err != nil {
		return nil, err
	}

	// Use the correct API endpoint format (no query parameter)
	fullURL := apiURL

	// TODO: Add tag filtering support when server implements it
	// For now, we'll filter client-side

	// Make the HTTP request
	ctx := context.Background()
	resp, err := p.doRequest(ctx, http.MethodGet, fullURL, nil, parameters)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Handle error responses
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the response
	var listResp listCommitsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to remote.Commit format
	commits := make([]remote.Commit, 0, len(listResp.Commits))
	for _, cr := range listResp.Commits {
		commit := parseCommitResponse(cr)

		// Apply tag filtering if tags are specified
		if len(tags) == 0 {
			commits = append(commits, commit)
		} else if matchesTags(commit, tags) {
			commits = append(commits, commit)
		}
	}

	return commits, nil
}

// GetCommit retrieves a specific commit by its identifier from the remote server.
func (p *Provider) GetCommit(properties map[string]interface{}, parameters map[string]interface{}, commitID string) (*remote.Commit, error) {
	// Build the API URL
	apiURL, err := p.buildAPIURL(properties, fmt.Sprintf("/commits/%s", commitID))
	if err != nil {
		return nil, err
	}

	// Make the HTTP request
	ctx := context.Background()
	resp, err := p.doRequest(ctx, http.MethodGet, apiURL, nil, parameters)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Handle 404 - commit not found
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	// Handle other error responses
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the response
	var cr commitResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	commit := parseCommitResponse(cr)
	return &commit, nil
}

// matchesTags checks if a commit matches all specified tag filters
func matchesTags(commit remote.Commit, tags []remote.Tag) bool {
	commitTags, ok := commit.Properties["tags"].(map[string]string)
	if !ok {
		return len(tags) == 0
	}

	for _, tag := range tags {
		value, exists := commitTags[tag.Key]

		// If tag key doesn't exist in commit, no match
		if !exists {
			return false
		}

		// If tag value is specified and doesn't match, no match
		if tag.Value != nil && value != *tag.Value {
			return false
		}
	}

	return true
}

func init() {
	p := &Provider{}
	remote.Register(p)
	// NOTE: We don't register separate http/https providers because
	// the datadatdat provider handles both http:// and https:// URLs
	// directly in its FromURL() method. Registering separate providers
	// would cause the wrong provider name to be returned.
}
