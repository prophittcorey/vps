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
		if val, err := vps.Check(ip); err == nil {
			fmt.Println(val)
		}

		return
	}

	flag.Usage()
}
