/*
 * Copyright Datadatdat.
 */
package datadatdat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/datadatdat/remote-sdk-go/remote"
	"github.com/stretchr/testify/assert"
)

func TestRegistered(t *testing.T) {
	r := remote.Get("datadatdat")

	ret, err := r.Type()
	if assert.NoError(t, err) {
		assert.Equal(t, "datadatdat", ret)
	}
}

// TestFromURL tests parsing complete HTTP URLs
func TestFromURL(t *testing.T) {
	r := remote.Get("datadatdat")

	props, err := r.FromURL("http://data.datadatdat.io:8080/myorg/myrepo", map[string]string{})
	if assert.NoError(t, err) {
		assert.Equal(t, "http://data.datadatdat.io:8080", props["api_base_url"])
		assert.Equal(t, "myorg", props["org"])
		assert.Equal(t, "myrepo", props["repo"])
		assert.Equal(t, "http", props["scheme"])
		assert.Equal(t, "data.datadatdat.io", props["host"])
		assert.Equal(t, 8080, props["port"])
	}
}

// TestFromURLHTTPS tests parsing HTTPS URLs
func TestFromURLHTTPS(t *testing.T) {
	r := remote.Get("datadatdat")

	props, err := r.FromURL("https://data.datadatdat.io/myorg/myrepo", map[string]string{})
	if assert.NoError(t, err) {
		assert.Equal(t, "https://data.datadatdat.io", props["api_base_url"])
		assert.Equal(t, "myorg", props["org"])
		assert.Equal(t, "myrepo", props["repo"])
		assert.Equal(t, "https", props["scheme"])
		assert.Equal(t, "data.datadatdat.io", props["host"])
		assert.Nil(t, props["port"])
	}
}

// TestFromURLSimple tests parsing URLs without port
func TestFromURLSimple(t *testing.T) {
	r := remote.Get("datadatdat")

	props, err := r.FromURL("http://localhost/org/repo", map[string]string{})
	if assert.NoError(t, err) {
		assert.Equal(t, "http://localhost", props["api_base_url"])
		assert.Equal(t, "org", props["org"])
		assert.Equal(t, "repo", props["repo"])
		assert.Equal(t, "http", props["scheme"])
		assert.Equal(t, "localhost", props["host"])
		assert.Nil(t, props["port"])
	}
}

// TestFromURLWithAPIToken tests passing API token as parameter
func TestFromURLWithAPIToken(t *testing.T) {
	r := remote.Get("datadatdat")

	props, err := r.FromURL("http://localhost/org/repo", map[string]string{"api_token": "secret123"})
	if assert.NoError(t, err) {
		assert.Equal(t, "secret123", props["api_token"])
	}
}

// TestBadURLMalformed tests rejection of malformed URLs
func TestBadURLMalformed(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.FromURL("http://host\nname", map[string]string{})
	assert.Error(t, err)
}

// TestBadScheme tests rejection of wrong scheme
func TestBadScheme(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.FromURL("ssh://host/org/repo", map[string]string{})
	assert.Error(t, err)
}

// TestBadSchemeFTP tests rejection of FTP scheme
func TestBadSchemeFTP(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.FromURL("ftp://host/org/repo", map[string]string{})
	assert.Error(t, err)
}

// TestBadProperty tests rejection of unknown properties
func TestBadProperty(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.FromURL("http://host/org/repo", map[string]string{"foo": "bar"})
	assert.Error(t, err)
}

// TestBadMissingHost tests rejection of URLs without host
func TestBadMissingHost(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.FromURL("http:///org/repo", map[string]string{})
	assert.Error(t, err)
}

// TestBadSchemeOnly tests rejection of scheme-only strings
func TestBadSchemeOnly(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.FromURL("http", map[string]string{})
	assert.Error(t, err)
}

// TestBadMissingOrg tests rejection of URLs without org
func TestBadMissingOrg(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.FromURL("http://host/", map[string]string{})
	assert.Error(t, err)
}

// TestBadMissingRepo tests rejection of URLs without repo
func TestBadMissingRepo(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.FromURL("http://host/org", map[string]string{})
	assert.Error(t, err)
}

// TestBadMissingRepoTrailingSlash tests rejection with trailing slash but no repo
func TestBadMissingRepoTrailingSlash(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.FromURL("http://host/org/", map[string]string{})
	assert.Error(t, err)
}

// TestBadPort tests rejection of invalid port numbers
func TestBadPort(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.FromURL("http://host:99999999999999999999/org/repo", map[string]string{})
	assert.Error(t, err)
}

// TestBadNegativePort tests rejection of negative ports
func TestBadNegativePort(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.FromURL("http://host:-1/org/repo", map[string]string{})
	assert.Error(t, err)
}

// TestBadExtraPathSegments tests rejection of URLs with extra path segments
func TestBadExtraPathSegments(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.FromURL("http://host/org/repo/extra", map[string]string{})
	assert.Error(t, err)
}

// TestToURL tests reconstructing URLs from properties
func TestToURL(t *testing.T) {
	r := remote.Get("datadatdat")

	u, props, err := r.ToURL(map[string]interface{}{
		"api_base_url": "http://localhost:8080",
		"org":          "myorg",
		"repo":         "myrepo",
	})
	if assert.NoError(t, err) {
		assert.Equal(t, "http://localhost:8080/myorg/myrepo", u)
		assert.Empty(t, props)
	}
}

// TestToURLHTTPS tests reconstructing HTTPS URLs
func TestToURLHTTPS(t *testing.T) {
	r := remote.Get("datadatdat")

	u, props, err := r.ToURL(map[string]interface{}{
		"api_base_url": "https://data.datadatdat.io",
		"org":          "myorg",
		"repo":         "myrepo",
	})
	if assert.NoError(t, err) {
		assert.Equal(t, "https://data.datadatdat.io/myorg/myrepo", u)
		assert.Empty(t, props)
	}
}

// TestToURLWithToken tests that API tokens are returned as properties, not in URL
func TestToURLWithToken(t *testing.T) {
	r := remote.Get("datadatdat")

	u, props, err := r.ToURL(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          "org",
		"repo":         "repo",
		"api_token":    "secret123",
	})
	if assert.NoError(t, err) {
		assert.Equal(t, "http://localhost/org/repo", u)
		assert.Len(t, props, 1)
		assert.Equal(t, "secret123", props["api_token"])
	}
}

// TestToURLMissingAPIBaseURL tests error when api_base_url is missing
func TestToURLMissingAPIBaseURL(t *testing.T) {
	r := remote.Get("datadatdat")
	_, _, err := r.ToURL(map[string]interface{}{
		"org":  "org",
		"repo": "repo",
	})
	assert.Error(t, err)
}

// TestToURLMissingOrg tests error when org is missing
func TestToURLMissingOrg(t *testing.T) {
	r := remote.Get("datadatdat")
	_, _, err := r.ToURL(map[string]interface{}{
		"api_base_url": "http://localhost",
		"repo":         "repo",
	})
	assert.Error(t, err)
}

// TestToURLMissingRepo tests error when repo is missing
func TestToURLMissingRepo(t *testing.T) {
	r := remote.Get("datadatdat")
	_, _, err := r.ToURL(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          "org",
	})
	assert.Error(t, err)
}

// TestToURLBadAPIBaseURLType tests error when api_base_url is wrong type
func TestToURLBadAPIBaseURLType(t *testing.T) {
	r := remote.Get("datadatdat")
	_, _, err := r.ToURL(map[string]interface{}{
		"api_base_url": 123,
		"org":          "org",
		"repo":         "repo",
	})
	assert.Error(t, err)
}

// TestToURLBadOrgType tests error when org is wrong type
func TestToURLBadOrgType(t *testing.T) {
	r := remote.Get("datadatdat")
	_, _, err := r.ToURL(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          123,
		"repo":         "repo",
	})
	assert.Error(t, err)
}

// TestToURLBadRepoType tests error when repo is wrong type
func TestToURLBadRepoType(t *testing.T) {
	r := remote.Get("datadatdat")
	_, _, err := r.ToURL(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          "org",
		"repo":         123,
	})
	assert.Error(t, err)
}

// TestGetParameters tests parameter pass-through
func TestGetParameters(t *testing.T) {
	r := remote.Get("datadatdat")

	props, err := r.GetParameters(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          "org",
		"repo":         "repo",
	})
	if assert.NoError(t, err) {
		assert.Equal(t, "http://localhost", props["api_base_url"])
		assert.Equal(t, "org", props["org"])
		assert.Equal(t, "repo", props["repo"])
	}
}

// TestGetParametersWithToken tests API token pass-through
func TestGetParametersWithToken(t *testing.T) {
	r := remote.Get("datadatdat")

	props, err := r.GetParameters(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          "org",
		"repo":         "repo",
		"api_token":    "secret123",
	})
	if assert.NoError(t, err) {
		assert.Equal(t, "secret123", props["api_token"])
	}
}

// TestValidateRemoteRequiredOnly tests validation with only required properties
func TestValidateRemoteRequiredOnly(t *testing.T) {
	r := remote.Get("datadatdat")
	err := r.ValidateRemote(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          "org",
		"repo":         "repo",
	})
	assert.NoError(t, err)
}

// TestValidateRemoteAllProperties tests validation with all properties
func TestValidateRemoteAllProperties(t *testing.T) {
	r := remote.Get("datadatdat")
	err := r.ValidateRemote(map[string]interface{}{
		"api_base_url": "http://localhost:8080",
		"org":          "org",
		"repo":         "repo",
		"scheme":       "http",
		"host":         "localhost",
		"port":         8080,
		"api_token":    "secret123",
	})
	assert.NoError(t, err)
}

// TestValidateRemoteMissingAPIBaseURL tests error when api_base_url is missing
func TestValidateRemoteMissingAPIBaseURL(t *testing.T) {
	r := remote.Get("datadatdat")
	err := r.ValidateRemote(map[string]interface{}{
		"org":  "org",
		"repo": "repo",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "api_base_url")
}

// TestValidateRemoteMissingOrg tests error when org is missing
func TestValidateRemoteMissingOrg(t *testing.T) {
	r := remote.Get("datadatdat")
	err := r.ValidateRemote(map[string]interface{}{
		"api_base_url": "http://localhost",
		"repo":         "repo",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "org")
}

// TestValidateRemoteMissingRepo tests error when repo is missing
func TestValidateRemoteMissingRepo(t *testing.T) {
	r := remote.Get("datadatdat")
	err := r.ValidateRemote(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          "org",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repo")
}

// TestValidateRemoteExtraProperty tests error with unknown property
func TestValidateRemoteExtraProperty(t *testing.T) {
	r := remote.Get("datadatdat")
	err := r.ValidateRemote(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          "org",
		"repo":         "repo",
		"foo":          "bar",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "foo")
}

// TestValidateRemoteBadPortType tests error when port is wrong type
func TestValidateRemoteBadPortType(t *testing.T) {
	r := remote.Get("datadatdat")
	err := r.ValidateRemote(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          "org",
		"repo":         "repo",
		"port":         "not_a_number",
	})
	assert.Error(t, err)
}

// TestValidateRemoteBadPortNegative tests error when port is negative
func TestValidateRemoteBadPortNegative(t *testing.T) {
	r := remote.Get("datadatdat")
	err := r.ValidateRemote(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          "org",
		"repo":         "repo",
		"port":         -1,
	})
	assert.Error(t, err)
}

// TestValidateRemotePortFloat tests that float ports are accepted and converted
func TestValidateRemotePortFloat(t *testing.T) {
	r := remote.Get("datadatdat")
	err := r.ValidateRemote(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          "org",
		"repo":         "repo",
		"port":         8080.0,
	})
	assert.NoError(t, err)
}

// TestValidateRemotePortFloat32 tests that float32 ports work
func TestValidateRemotePortFloat32(t *testing.T) {
	r := remote.Get("datadatdat")

	var p float32 = 8080.0

	err := r.ValidateRemote(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          "org",
		"repo":         "repo",
		"port":         p,
	})
	assert.NoError(t, err)
}

// TestValidateParametersEmpty tests validation with no parameters
func TestValidateParametersEmpty(t *testing.T) {
	r := remote.Get("datadatdat")
	err := r.ValidateParameters(map[string]interface{}{})
	assert.NoError(t, err)
}

// TestValidateParametersWithToken tests validation with API token
func TestValidateParametersWithToken(t *testing.T) {
	r := remote.Get("datadatdat")
	err := r.ValidateParameters(map[string]interface{}{
		"api_token": "secret123",
	})
	assert.NoError(t, err)
}

// TestValidateParametersUnknown tests error with unknown parameter
func TestValidateParametersUnknown(t *testing.T) {
	r := remote.Get("datadatdat")
	err := r.ValidateParameters(map[string]interface{}{
		"foo": "bar",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "foo")
}

// HTTP Tests with mock server

// TestListCommitsSuccess tests successful listing of commits
func TestListCommitsSuccess(t *testing.T) {
	// Create mock HTTP server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "/api/v1/repos/testorg/testrepo/commits", r.URL.Path)
		assert.Equal(t, "list-commits", r.URL.Query().Get("action"))
		assert.Equal(t, "GET", r.Method)

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(listCommitsResponse{
			Repo: "testrepo",
			Commits: []commitResponse{
				{
					CommitID:  "commit1",
					Repo:      "testrepo",
					Timestamp: time.Now(),
					Size:      1024,
					Author:    "user1",
					Message:   "First commit",
				},
				{
					CommitID:  "commit2",
					Repo:      "testrepo",
					Timestamp: time.Now().Add(-1 * time.Hour),
					Size:      2048,
					Author:    "user2",
					Message:   "Second commit",
					Tags:      []string{"env:prod", "version:1.0"},
				},
			},
		})
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Test ListCommits
	r := remote.Get("datadatdat")
	commits, err := r.ListCommits(
		map[string]interface{}{
			"api_base_url": server.URL,
			"org":          "testorg",
			"repo":         "testrepo",
		},
		map[string]interface{}{},
		[]remote.Tag{},
	)

	assert.NoError(t, err)
	assert.Len(t, commits, 2)
	assert.Equal(t, "commit1", commits[0].ID)
	assert.Equal(t, "commit2", commits[1].ID)
	assert.Equal(t, "user1", commits[0].Properties["author"])
	assert.Equal(t, "Second commit", commits[1].Properties["message"])
}

// TestListCommitsWithAuthentication tests that Bearer token is sent
func TestListCommitsWithAuthentication(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Bearer token
		assert.Equal(t, "Bearer secret123", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(listCommitsResponse{
			Repo:    "testrepo",
			Commits: []commitResponse{},
		})
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	r := remote.Get("datadatdat")
	_, err := r.ListCommits(
		map[string]interface{}{
			"api_base_url": server.URL,
			"org":          "testorg",
			"repo":         "testrepo",
		},
		map[string]interface{}{
			"api_token": "secret123",
		},
		[]remote.Tag{},
	)

	assert.NoError(t, err)
}

// TestListCommitsWithTagFiltering tests client-side tag filtering
func TestListCommitsWithTagFiltering(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(listCommitsResponse{
			Repo: "testrepo",
			Commits: []commitResponse{
				{
					CommitID: "commit1",
					Tags:     []string{"env:prod", "version:1.0"},
				},
				{
					CommitID: "commit2",
					Tags:     []string{"env:dev", "version:1.0"},
				},
				{
					CommitID: "commit3",
					Tags:     []string{"env:prod", "version:2.0"},
				},
			},
		})
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	r := remote.Get("datadatdat")

	// Filter for env:prod
	envProd := "prod"
	commits, err := r.ListCommits(
		map[string]interface{}{
			"api_base_url": server.URL,
			"org":          "testorg",
			"repo":         "testrepo",
		},
		map[string]interface{}{},
		[]remote.Tag{{Key: "env", Value: &envProd}},
	)

	assert.NoError(t, err)
	assert.Len(t, commits, 2)
	assert.Equal(t, "commit1", commits[0].ID)
	assert.Equal(t, "commit3", commits[1].ID)
}

// TestListCommitsServerError tests handling of server errors
func TestListCommitsServerError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	r := remote.Get("datadatdat")
	_, err := r.ListCommits(
		map[string]interface{}{
			"api_base_url": server.URL,
			"org":          "testorg",
			"repo":         "testrepo",
		},
		map[string]interface{}{},
		[]remote.Tag{},
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

// TestListCommitsInvalidResponse tests handling of invalid JSON
func TestListCommitsInvalidResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	r := remote.Get("datadatdat")
	_, err := r.ListCommits(
		map[string]interface{}{
			"api_base_url": server.URL,
			"org":          "testorg",
			"repo":         "testrepo",
		},
		map[string]interface{}{},
		[]remote.Tag{},
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decode")
}

// TestGetCommitSuccess tests successful retrieval of a commit
func TestGetCommitSuccess(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/repos/testorg/testrepo/commits/commit123", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(commitResponse{
			CommitID:  "commit123",
			Repo:      "testrepo",
			Timestamp: time.Now(),
			Size:      1024,
			Author:    "testuser",
			Message:   "Test commit",
			ParentID:  "parent123",
			Tags:      []string{"env:test"},
		})
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	r := remote.Get("datadatdat")
	commit, err := r.GetCommit(
		map[string]interface{}{
			"api_base_url": server.URL,
			"org":          "testorg",
			"repo":         "testrepo",
		},
		map[string]interface{}{},
		"commit123",
	)

	assert.NoError(t, err)
	assert.NotNil(t, commit)
	assert.Equal(t, "commit123", commit.ID)
	assert.Equal(t, "testuser", commit.Properties["author"])
	assert.Equal(t, "Test commit", commit.Properties["message"])
	assert.Equal(t, "parent123", commit.Properties["parentId"])
}

// TestGetCommitNotFound tests 404 response
func TestGetCommitNotFound(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Commit not found"))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	r := remote.Get("datadatdat")
	commit, err := r.GetCommit(
		map[string]interface{}{
			"api_base_url": server.URL,
			"org":          "testorg",
			"repo":         "testrepo",
		},
		map[string]interface{}{},
		"nonexistent",
	)

	assert.NoError(t, err)
	assert.Nil(t, commit)
}

// TestGetCommitWithAuthentication tests that Bearer token is sent
func TestGetCommitWithAuthentication(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer mytoken", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(commitResponse{
			CommitID: "commit123",
		})
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	r := remote.Get("datadatdat")
	_, err := r.GetCommit(
		map[string]interface{}{
			"api_base_url": server.URL,
			"org":          "testorg",
			"repo":         "testrepo",
		},
		map[string]interface{}{
			"api_token": "mytoken",
		},
		"commit123",
	)

	assert.NoError(t, err)
}

// TestGetCommitServerError tests handling of server errors
func TestGetCommitServerError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database error"))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	r := remote.Get("datadatdat")
	_, err := r.GetCommit(
		map[string]interface{}{
			"api_base_url": server.URL,
			"org":          "testorg",
			"repo":         "testrepo",
		},
		map[string]interface{}{},
		"commit123",
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

// TestGetCommitInvalidResponse tests handling of invalid JSON
func TestGetCommitInvalidResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	r := remote.Get("datadatdat")
	_, err := r.GetCommit(
		map[string]interface{}{
			"api_base_url": server.URL,
			"org":          "testorg",
			"repo":         "testrepo",
		},
		map[string]interface{}{},
		"commit123",
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decode")
}

// TestBuildAPIURLMissingProperties tests error handling for missing properties
func TestBuildAPIURLMissingProperties(t *testing.T) {
	p := NewProvider()

	// Missing api_base_url
	_, err := p.buildAPIURL(map[string]interface{}{
		"org":  "testorg",
		"repo": "testrepo",
	}, "/commits")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "api_base_url")

	// Missing org
	_, err = p.buildAPIURL(map[string]interface{}{
		"api_base_url": "http://localhost",
		"repo":         "testrepo",
	}, "/commits")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "org")

	// Missing repo
	_, err = p.buildAPIURL(map[string]interface{}{
		"api_base_url": "http://localhost",
		"org":          "testorg",
	}, "/commits")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repo")
}

// TestParseCommitResponseMinimal tests parsing a minimal commit response
func TestParseCommitResponseMinimal(t *testing.T) {
	cr := commitResponse{
		CommitID:  "minimal123",
		Timestamp: time.Now(),
		Size:      100,
	}

	commit := parseCommitResponse(cr)
	assert.Equal(t, "minimal123", commit.ID)
	assert.NotNil(t, commit.Properties["timestamp"])
	assert.Equal(t, int64(100), commit.Properties["size"])
	assert.Nil(t, commit.Properties["author"])
	assert.Nil(t, commit.Properties["message"])
}

// TestParseCommitResponseWithMetadata tests parsing commit with custom metadata
func TestParseCommitResponseWithMetadata(t *testing.T) {
	cr := commitResponse{
		CommitID:  "meta123",
		Timestamp: time.Now(),
		Size:      200,
		Metadata: map[string]interface{}{
			"custom_field": "custom_value",
			"number":       42,
		},
	}

	commit := parseCommitResponse(cr)
	assert.Equal(t, "meta123", commit.ID)
	assert.Equal(t, "custom_value", commit.Properties["custom_field"])
	assert.Equal(t, 42, commit.Properties["number"])
}

// TestParseCommitResponseWithSingleTag tests parsing a commit with single tag
func TestParseCommitResponseWithSingleTag(t *testing.T) {
	cr := commitResponse{
		CommitID:  "tag123",
		Timestamp: time.Now(),
		Size:      300,
		Tags:      []string{"single"},
	}

	commit := parseCommitResponse(cr)
	tagsMap, ok := commit.Properties["tags"].(map[string]string)
	assert.True(t, ok)
	assert.Equal(t, "", tagsMap["single"])
}

// TestMatchesTagsEmpty tests matching with no tags specified
func TestMatchesTagsEmpty(t *testing.T) {
	commit := remote.Commit{
		ID:         "commit1",
		Properties: map[string]interface{}{},
	}

	// Should match when no tags specified
	assert.True(t, matchesTags(commit, []remote.Tag{}))
}

// TestMatchesTagsNoCommitTags tests matching when commit has no tags
func TestMatchesTagsNoCommitTags(t *testing.T) {
	commit := remote.Commit{
		ID:         "commit1",
		Properties: map[string]interface{}{},
	}

	envVal := "prod"
	// Should not match when filter tags specified but commit has none
	assert.False(t, matchesTags(commit, []remote.Tag{{Key: "env", Value: &envVal}}))
}

// TestMatchesTagsKeyOnly tests matching with key-only tag filter
func TestMatchesTagsKeyOnly(t *testing.T) {
	commit := remote.Commit{
		ID: "commit1",
		Properties: map[string]interface{}{
			"tags": map[string]string{
				"env":     "prod",
				"version": "1.0",
			},
		},
	}

	// Should match when key exists (value not specified in filter)
	assert.True(t, matchesTags(commit, []remote.Tag{{Key: "env", Value: nil}}))
}

// TestListCommitsEmpty tests listing commits with empty response
func TestListCommitsEmpty(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(listCommitsResponse{
			Repo:    "testrepo",
			Commits: []commitResponse{},
		})
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	r := remote.Get("datadatdat")
	commits, err := r.ListCommits(
		map[string]interface{}{
			"api_base_url": server.URL,
			"org":          "testorg",
			"repo":         "testrepo",
		},
		map[string]interface{}{},
		[]remote.Tag{},
	)

	assert.NoError(t, err)
	assert.Len(t, commits, 0)
}

// TestListCommitsBadProperties tests error when properties are invalid
func TestListCommitsBadProperties(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.ListCommits(
		map[string]interface{}{
			"org":  "testorg",
			"repo": "testrepo",
			// Missing api_base_url
		},
		map[string]interface{}{},
		[]remote.Tag{},
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "api_base_url")
}

// TestGetCommitBadProperties tests error when properties are invalid
func TestGetCommitBadProperties(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.GetCommit(
		map[string]interface{}{
			"org":  "testorg",
			"repo": "testrepo",
			// Missing api_base_url
		},
		map[string]interface{}{},
		"commit123",
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "api_base_url")
}
