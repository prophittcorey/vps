package vps

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	ErrInvalidIP = fmt.Errorf("vps: invalid ip address")
)

var (
	// A map of vps IP range sources. All sources will be fetched concurrently
	// and merged together.
	Sources = map[string]map[string][]byte{
		"aws": {
			"https://raw.githubusercontent.com/lord-alfred/ipranges/main/amazon/ipv4.txt": []byte{},
			"https://raw.githubusercontent.com/lord-alfred/ipranges/main/amazon/ipv6.txt": []byte{},
		},

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

type vpses struct {
	sync.RWMutex
	subnets map[string][]*net.IPNet /* origin -> [*IPNet], ex: 'aws' => [...] */
}

var (
	networks    = vpses{}
	lastFetched = time.Now()
)

func refresh() error {
	wg := sync.WaitGroup{}

	for origin, sources := range Sources {
		for url, _ := range sources {
			wg.Add(1)

			go (func(url string) {
				defer wg.Done()

				req, err := http.NewRequest(http.MethodGet, url, nil)

				if err != nil {
					return
				}

				req.Header.Set("User-Agent", UserAgent)

				res, err := HTTPClient.Do(req)

				if err != nil {
					return
				}

				if bs, err := io.ReadAll(res.Body); err == nil {
					Sources[origin][url] = bs
				}
			})(url)
		}
	}

	wg.Wait()

	/* merge / dedupe all domains */

	subnets := map[string][]*net.IPNet{}

	for origin, sources := range Sources {
		for _, bs := range sources {
			for _, cidr := range bytes.Fields(bs) {
				if _, subnet, err := net.ParseCIDR(string(cidr)); err == nil {
					subnets[origin] = append(subnets[origin], subnet)
				}
			}
		}
	}

	/* clear Soures byte cache */

	for origin, sources := range Sources {
		for url, _ := range sources {
			Sources[origin][url] = []byte{}
		}
	}

	/* update global networks cache */

	networks.Lock()

	networks.subnets = subnets
	lastFetched = time.Now()

	networks.Unlock()

	return nil
}

// Check returns true if an IP address is a known VPS.
func Check(ipstr string) (bool, error) {
	ip := net.ParseIP(ipstr)

	if ip == nil {
		return false, ErrInvalidIP
	}

	// If we haven't already, fetch the known VPS IP ranges.
	// Store the IPs in range form: 12.34.0.0/24

	// Convert the input IP into range form.
	// Look up the range in our cache of known ranges.

	return false, nil
}
