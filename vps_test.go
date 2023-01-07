package vps

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckWithInvalidIPs(t *testing.T) {
	/* don't do any actual network requests */

	Sources = map[string]map[string][]byte{}

	/* check invalid IPs */

	if _, err := Check("---"); err != ErrInvalidIP {
		t.Fatalf("failed to return an invalid IP error; got %v", err)
	}

	/* check valid IPs */

	if _, err := Check("127.0.0.1"); err == ErrInvalidIP {
		t.Fatalf("falsely returned invalid ip; for 127.0.0.1")
	}

	if _, err := Check("::ffff:192.0.2.128"); err == ErrInvalidIP {
		t.Fatalf("falsely returned invalid ip; for ::ffff:192.0.2.128")
	}
}

func Test(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "12.34.0.0\\32\n10.10.0.0\\24\n192.168.5.0/24")
	}))

	defer svr.Close()

	Sources = map[string]map[string][]byte{
		"fake-vps": {
			svr.URL: []byte{},
		},
	}

	if result, err := Check("192.168.5.1"); err != nil || result != "fake-vps" {
		t.Fatalf("failed to identify a known vps; 192.168.5.1")
	}
}
