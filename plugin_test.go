package caddy_remote_host_test

// These tests target the public API. For tests of unexported fields and
// methods see package caddy_remote_host.

import (
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	plugin "github.com/muety/caddy-remote-host"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchRemoteHost_UnmarshalCaddyfile(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	for _, tt := range []struct {
		name      string
		input     string
		hosts     []string
		forwarded bool
		nocache   bool
	}{
		{
			name:  "simple",
			input: "remote_host example.com",
			hosts: []string{"example.com"},
		}, {
			name:  "list",
			input: "remote_host example.com example.org",
			hosts: []string{"example.com", "example.org"},
		}, {
			name:      "fwd",
			input:     "remote_host forwarded example.com",
			hosts:     []string{"example.com"},
			forwarded: true,
		}, {
			name:    "noc",
			input:   "remote_host nocache example.com",
			hosts:   []string{"example.com"},
			nocache: true,
		}, {
			name:      "fwdnoc",
			input:     "remote_host forwarded nocache example.com",
			hosts:     []string{"example.com"},
			forwarded: true,
			nocache:   true,
		}, {
			name:      "nocfwd",
			input:     "remote_host nocache forwarded example.com",
			hosts:     []string{"example.com"},
			forwarded: true,
			nocache:   true,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			disp := caddyfile.NewTestDispenser(tt.input)
			subject := plugin.MatchRemoteHost{}

			err := subject.UnmarshalCaddyfile(disp)
			require.NoError(err)

			assert.EqualValues(tt.hosts, subject.Hosts)
			assert.Equal(tt.forwarded, subject.Forwarded)
			assert.Equal(tt.nocache, subject.NoCache)
		})
	}
}

func TestMatchRemoteHost_UnmarshalCaddyfile_invalid(t *testing.T) {
	assert := assert.New(t)

	for _, tt := range []struct {
		name  string
		input string
		err   string
	}{
		{
			name:  "forwarded after host",
			input: "remote_host example.com forwarded",
			err:   "if used, 'forwarded' must appear before 'hosts' argument, at Testfile:1",
		}, {
			name:  "forwarded between hosts",
			input: "remote_host example.com forwarded example.org",
			err:   "if used, 'forwarded' must appear before 'hosts' argument, at Testfile:1",
		}, {
			name:  "nocache after host",
			input: "remote_host example.com nocache",
			err:   "if used, 'nocache' must appear before 'hosts' argument, at Testfile:1",
		}, {
			name:  "nocache between hosts",
			input: "remote_host example.com nocache example.org",
			err:   "if used, 'nocache' must appear before 'hosts' argument, at Testfile:1",
		}, {
			name:  "block",
			input: "remote_host example.com {\nnot supported\n}",
			err:   "malformed remote_host matcher: blocks are not supported, at Testfile:2",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			disp := caddyfile.NewTestDispenser(tt.input)
			subject := plugin.MatchRemoteHost{}

			err := subject.UnmarshalCaddyfile(disp)
			assert.EqualError(err, tt.err)
		})
	}
}
