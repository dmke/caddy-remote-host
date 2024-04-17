package caddy_remote_host

// These tests target unexported fields and methods. For tests of the
// public API see package caddy_remote_host_test.

import (
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchRemoteHost_Provision(t *testing.T) {
	hostRegex = nil

	subject := MatchRemoteHost{}
	err := subject.Provision(caddy.Context{})
	require.NoError(t, err)
	assert.NotNil(t, subject.logger)
	assert.NotNil(t, subject.cache)
	assert.NotNil(t, hostRegex)
}
