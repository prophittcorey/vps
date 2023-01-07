package vps

import "testing"

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
	if result, err := Check("12.34.56.78"); err != nil || result != "fake-vps" {
		t.Fatalf("failed to identify a known vps; 12.34.56.78")
	}
}
