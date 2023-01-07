package main

import (
	"flag"
	"fmt"

	"github.com/prophittcorey/vps"
)

func main() {
	var ip string

	flag.StringVar(&ip, "check", "", "an ip address to analyze (returns 'true' if it's a known disposable address, 'false' otherwise")

	flag.Parse()

	if len(ip) > 0 {
		vpsname, err := vps.Check(ip)

		if err == nil {
			fmt.Printf("Looks like a '%s' address.\n", vpsname)
		} else {
			fmt.Printf("Does not look like a vps address.\n")
		}

		return
	}

	flag.Usage()
}
