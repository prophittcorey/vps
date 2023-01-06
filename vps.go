package vps

import (
	"fmt"
	"net/http"
	"time"
)

var (
	ErrInvalidIP = fmt.Errorf("vps: invalid ip address")
)

var (
	// A map of vps IP range sources. All sources will be fetched concurrently
	// and merged together.
	Sources = map[string][]byte{
		"https://raw.githubusercontent.com/lord-alfred/ipranges/main/amazon/ipv4.txt": []byte{},
		"https://raw.githubusercontent.com/lord-alfred/ipranges/main/amazon/ipv6.txt": []byte{},

		// TODO: Add additional sources.
	}

	// HTTPClient is used to perform all HTTP requests. You can specify your own
	// to set a custom timeout, proxy, etc.
	HTTPClient = http.Client{
		Timeout: 3 * time.Second,
	}

	// CachePeriod specifies the amount of time an internal cache of disposable email domains are used
	// before refreshing the domains.
	CachePeriod = 45 * time.Minute

	// UserAgent will be used in each request's user agent header field.
	UserAgent = "github.com/prophittcorey/vps"
)

// Check returns true if an IP address is a known VPS.
func Check(ip string) (bool, error) {
	// First, check if the IP is even valid. Does it work for IPv4 and IPv6?

	// If we haven't already, fetch the known VPS IP ranges.
	// Store the IPs in range form: 12.34.0.0/24

	// Convert the input IP into range form.
	// Look up the range in our cache of known ranges.

	return false, nil
}
