package vps

import "testing"

func Test(t *testing.T) {
	if ok, err := Check("12.34.56.78"); err != nil || !ok {
		t.Fatalf("failed to identify a known vps; 12.34.56.78")
	}
}
