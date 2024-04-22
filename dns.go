package caddy_remote_host

import (
	"context"
	"net"
)

// A resolver is able to lookup IP addresses from a given host name.
// The net.Resolver type (found at net.DefaultResolver) matches this
// interface.
//
// This is intended for testing.
type resolver interface {
	LookupIPAddr(ctx context.Context, host string) (addrs []net.IPAddr, err error)
}

// lookupIP does the same thing as net.LookupIP, except it doesn't use
// net.DefaultResolver when given an alternative resolver implementation
// as first argument.
// (Setting resolv to nil will fallback to net.DefaultResolver.)
func lookupIP(resolv resolver, host string) ([]net.IP, error) {
	r := resolv
	if r == nil {
		r = net.DefaultResolver
	}

	addrs, err := r.LookupIPAddr(context.Background(), host)
	if err != nil {
		return nil, err
	}
	ips := make([]net.IP, len(addrs))
	for i, ia := range addrs {
		ips[i] = ia.IP
	}
	return ips, nil
}
