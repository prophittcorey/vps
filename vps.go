package vps

import (
	"net/http"
	"time"
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
	return false, nil
}
