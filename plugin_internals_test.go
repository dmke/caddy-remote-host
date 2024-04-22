package caddy_remote_host

// These tests target unexported fields and methods. For tests of the
// public API see package caddy_remote_host_test.

import (
	"context"
	"errors"
	"net"
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
	assert.Nil(t, subject.resolver)
	assert.NotNil(t, hostRegex)
}

type lookupResult struct {
	ips []net.IPAddr
	err error
}

func resolvesTo(ips ...string) lookupResult {
	r := lookupResult{ips: make([]net.IPAddr, len(ips))}
	for i, ip := range ips {
		r.ips[i] = net.IPAddr{IP: net.ParseIP(ip)}
	}
	return r
}

type mockResolver struct {
	addrs map[string]lookupResult
}

func (m *mockResolver) LookupIPAddr(_ context.Context, h string) ([]net.IPAddr, error) {
	if result, ok := m.addrs[h]; ok {
		return result.ips, result.err
	}
	return nil, errors.New("no suitable address found")
}

func TestMatchRemoteHost_resolveIPs(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	mock := &mockResolver{
		addrs: map[string]lookupResult{
			"example.com":         resolvesTo("127.0.0.1", "127.0.0.2"),
			"example.org":         resolvesTo("::1", "fe80::1"),
			"nil.records.example": {},
			"no.records.example":  {ips: make([]net.IPAddr, 0)},
		},
	}

	subject := MatchRemoteHost{resolver: mock}
	for host := range mock.addrs {
		subject.Hosts = append(subject.Hosts, host)
	}

	require.NoError(subject.Provision(caddy.Context{}))
	ips, err := subject.resolveIPs()
	require.NoError(err)

	var haveIPs []string
	for _, result := range mock.addrs {
		for _, ip := range result.ips {
			haveIPs = append(haveIPs, ip.String())
		}
	}

	toStrings := func(ips []net.IP) []string {
		s := make([]string, len(ips))
		for i, ip := range ips {
			s[i] = ip.String()
		}
		return s
	}

	assert.ElementsMatch(haveIPs, toStrings(ips))

	cachedIPs, ok := subject.cache.Get(cacheKey)
	require.True(ok)
	require.NotNil(cachedIPs)
	assert.ElementsMatch(haveIPs, toStrings(cachedIPs.([]net.IP)))
}

func TestMatchRemoteHost_resolveIPs_failure(t *testing.T) {
	subject := MatchRemoteHost{
		Hosts: []string{"example.com"},
		resolver: &mockResolver{map[string]lookupResult{
			"example.com": {err: &net.DNSError{Err: "no suitable host found"}},
		}},
	}

	require.NoError(t, subject.Provision(caddy.Context{}))

	ips, err := subject.resolveIPs()
	// XXX: the expected message is constructed within package net
	// and might change with future Go versions
	assert.EqualError(t, err, "lookup : no suitable host found")
	assert.Empty(t, ips)
}
