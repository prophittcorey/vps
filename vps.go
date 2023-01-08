package vps

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	ErrInvalidIP = fmt.Errorf("vps: invalid ip address")
	ErrNotFound  = fmt.Errorf("vps: ip address isn't associated with a known vps")
)

var (
	// A map of vps IP range sources. All sources will be fetched concurrently
	// and merged together.
	Sources = map[string]map[string][]byte{
		"aws": {
			"https://raw.githubusercontent.com/lord-alfred/ipranges/main/amazon/ipv4.txt": []byte{},
			"https://raw.githubusercontent.com/lord-alfred/ipranges/main/amazon/ipv6.txt": []byte{},
		},
		"linode": {
			"https://raw.githubusercontent.com/lord-alfred/ipranges/main/linode/ipv4.txt": []byte{},
			"https://raw.githubusercontent.com/lord-alfred/ipranges/main/linode/ipv6.txt": []byte{},
		},
		"digitalocean": {
			"https://raw.githubusercontent.com/lord-alfred/ipranges/main/digitalocean/ipv4.txt": []byte{},
			"https://raw.githubusercontent.com/lord-alfred/ipranges/main/digitalocean/ipv6.txt": []byte{},
		},
		"azure": {
			"https://raw.githubusercontent.com/lord-alfred/ipranges/main/microsoft/ipv4.txt": []byte{},
			"https://raw.githubusercontent.com/lord-alfred/ipranges/main/microsoft/ipv6.txt": []byte{},
		},
		"googlecloud": {
			"https://raw.githubusercontent.com/lord-alfred/ipranges/main/google/ipv4.txt": []byte{},
			"https://raw.githubusercontent.com/lord-alfred/ipranges/main/google/ipv6.txt": []byte{},
		},
		"oracle": {
			"https://raw.githubusercontent.com/lord-alfred/ipranges/main/oracle/ipv4.txt": []byte{},
		},
		"hetzner": {
			"https://raw.githubusercontent.com/Pymmdrza/Datacenter_List_DataBase_IP/mainx/Hetzner/CIDR.txt": []byte{},
		},

		// TODO: Add additional sources (and/or backup sources if one or more fail).
	}

	// HTTPClient is used to perform all HTTP requests. You can specify your own
	// to set a custom timeout, proxy, etc.
	HTTPClient = http.Client{
		Timeout: 3 * time.Second,
	}

	// CachePeriod specifies the amount of time an internal cache of vps subnets are used
	// before refreshing the subnets.
	CachePeriod = 45 * time.Minute

	// UserAgent will be used in each request's user agent header field.
	UserAgent = "github.com/prophittcorey/vps"
)

// TODO: Is there a more efficient data structure for IP -> Subnet look ups? Currently it's
// O(n), but there is likely a O(1) solution out there.

type vpses struct {
	sync.RWMutex
	subnets map[string][]*net.IPNet /* origin -> [*IPNet], ex: 'aws' => [...] */
}

var (
	networks    = vpses{}
	lastFetched = time.Now()
)

func refresh() error {
	if len(networks.subnets) == 0 || time.Now().After(lastFetched.Add(CachePeriod)) {
		wg := sync.WaitGroup{}

		for origin, sources := range Sources {
			origin := origin
			sources := sources

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
			origin := origin
			sources := sources

			for _, bs := range sources {
				for _, cidr := range bytes.Fields(bs) {
					if _, subnet, err := net.ParseCIDR(string(cidr)); err == nil {
						if subnet.String() == "45.79.99.0/24" {
							log.Println(origin)
						}
						subnets[origin] = append(subnets[origin], subnet)
					}
				}
			}
		}

		/* clear byte cache of sources */

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
	}

	return nil
}

// Check returns the VPS associated with the IP, or an error.
func Check(ipstr string) (string, error) {
	ip := net.ParseIP(ipstr)

	if ip == nil {
		return "", ErrInvalidIP
	}

	refresh()

	for origin, subnets := range networks.subnets {
		origin := origin
		subnets := subnets

		for _, subnet := range subnets {
			if subnet.Contains(ip) {
				return origin, nil
			}
		}
	}

	return "", ErrNotFound
}

// Subnets will return all known VPS subnets.
func Subnets() []*net.IPNet {
	refresh()

	ss := []*net.IPNet{}

	for _, subnets := range networks.subnets {
		for _, subnet := range subnets {
			ss = append(ss, subnet)
		}
	}

	return ss
}
