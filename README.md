# VPS

A golang package and command line tool for the identification of VPS IP
addresses.

## Package Usage

```golang
import "github.com/prophittcorey/vps"

if provider, err := vps.Check("12.34.56.78"); err == nil {
  fmt.Printf("Looks like a %s IP address.\n", provider)
}

vps.Subnets() // => [*IPNet, ...]
```

## Tool Usage

```bash
# Install the latest tool.
$ go install github.com/prophittcorey/vps/cmd/vps@latest

# Dump all known VPS subnets.
$ vps --subnets

# Check a specific ip address.
$ vps --check 12.34.56.78
```

## License

The source code for this repository is licensed under the MIT license, which you can
find in the [LICENSE](LICENSE.md) file.
