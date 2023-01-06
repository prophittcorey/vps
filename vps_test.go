package vps

import "testing"

func TestCheckWithInvalidIPs(t *testing.T) {
	if _, err := Check("---"); err != ErrInvalidIP {
		t.Fatalf("failed to return an invalid IP error; got %v", err)
	}
}

func Test(t *testing.T) {
	if ok, err := Check("12.34.56.78"); err != nil || !ok {
		t.Fatalf("failed to identify a known vps; 12.34.56.78")
	}
}
