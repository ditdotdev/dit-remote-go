/*
 * Copyright Datadatdat.
 */
package datadatdat

import (
	"testing"

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

// TestListCommitsNotImplemented tests that ListCommits returns error
func TestListCommitsNotImplemented(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.ListCommits(
		map[string]interface{}{
			"api_base_url": "http://localhost",
			"org":          "org",
			"repo":         "repo",
		},
		map[string]interface{}{},
		[]remote.Tag{},
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")
}

// TestGetCommitNotImplemented tests that GetCommit returns error
func TestGetCommitNotImplemented(t *testing.T) {
	r := remote.Get("datadatdat")
	_, err := r.GetCommit(
		map[string]interface{}{
			"api_base_url": "http://localhost",
			"org":          "org",
			"repo":         "repo",
		},
		map[string]interface{}{},
		"commit123",
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")
}
