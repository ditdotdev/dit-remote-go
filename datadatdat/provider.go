// Package datadatdat provides a remote provider for connecting d3 CLI to datadatdat-remote-server.
// This provider implements the remote.Remote interface and communicates with the datadatdat-remote-server
// HTTP APIs to store and retrieve commits.
package datadatdat

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/datadatdat/remote-sdk-go/remote"
)

// Provider implements the remote.Remote interface for datadatdat-remote-server.
type Provider struct {
}

// NewProvider creates a new datadatdat remote provider instance.
func NewProvider() *Provider {
	return &Provider{}
}

// Type returns the remote type identifier for the datadatdat remote.
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

	// Extract org/repo from path
	path := strings.Trim(parsedURL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		return nil, fmt.Errorf("invalid path: expected /org/repo format")
	}

	org := pathParts[0]
	repo := pathParts[1]

	// Build API base URL (scheme://host:port)
	apiBaseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	properties := map[string]interface{}{
		"api_base_url": apiBaseURL,
		"org":          org,
		"repo":         repo,
		"scheme":       parsedURL.Scheme,
		"host":         parsedURL.Host,
	}

	// Add any additional properties (e.g., API tokens)
	for k, v := range additionalProperties {
		properties[k] = v
	}

	return properties, nil
}

// ToURL converts remote properties back to a URL and additional properties.
func (p *Provider) ToURL(properties map[string]interface{}) (string, map[string]string, error) {
	apiBaseURL, ok := properties["api_base_url"].(string)
	if !ok {
		return "", nil, fmt.Errorf("missing api_base_url property")
	}

	org, ok := properties["org"].(string)
	if !ok {
		return "", nil, fmt.Errorf("missing org property")
	}

	repo, ok := properties["repo"].(string)
	if !ok {
		return "", nil, fmt.Errorf("missing repo property")
	}

	// Construct the full URL
	fullURL := fmt.Sprintf("%s/%s/%s", apiBaseURL, org, repo)

	// Extract additional properties (excluding system properties)
	additionalProps := make(map[string]string)
	systemProps := map[string]bool{
		"api_base_url": true,
		"org":          true,
		"repo":         true,
		"scheme":       true,
		"host":         true,
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
func (p *Provider) GetParameters(remoteProperties map[string]interface{}) (map[string]interface{}, error) {
	// For now, pass through all properties as parameters
	// In the future, this could prompt for API tokens if not present
	return remoteProperties, nil
}

// ValidateRemote validates the remote connection properties.
func (p *Provider) ValidateRemote(properties map[string]interface{}) error {
	requiredProps := []string{"api_base_url", "org", "repo"}
	for _, prop := range requiredProps {
		if _, ok := properties[prop]; !ok {
			return fmt.Errorf("missing required property: %s", prop)
		}
	}
	return nil
}

// ValidateParameters validates the operation parameters.
func (p *Provider) ValidateParameters(parameters map[string]interface{}) error {
	// For MVP, no additional validation needed
	// In the future, validate API tokens here
	return nil
}

// ListCommits returns a list of available commits from the remote server, optionally filtered by tags.
func (p *Provider) ListCommits(properties map[string]interface{}, parameters map[string]interface{}, tags []remote.Tag) ([]remote.Commit, error) {
	// TODO: Implement HTTP API call to list commits
	// GET /api/v1/repos/{org}/{repo}/commits?tag=key:value
	return nil, fmt.Errorf("ListCommits not yet implemented")
}

// GetCommit retrieves a specific commit by its identifier from the remote server.
func (p *Provider) GetCommit(properties map[string]interface{}, parameters map[string]interface{}, commitID string) (*remote.Commit, error) {
	// TODO: Implement HTTP API call to get commit metadata
	// GET /api/v1/repos/{org}/{repo}/commits/{commitID}
	return nil, fmt.Errorf("GetCommit not yet implemented")
}
