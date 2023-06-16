# go-gpsd

[GPSD](https://gpsd.gitlab.io/gpsd/index.html) client for Golang without CGO and additional dependencies.

It provides interface similar to the official C, C++ and Python libraries.

Tested with GPSD 3.16, proto version 3.11 and GPSD 3.22, proto version 3.14

## Usage

The following example illustrates basic concepts and how to use this library:

```go
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/amenzhinsky/go-gpsd"
)

var addrFlag string

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [FLAG...]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&addrFlag, "addr", gpsd.DefaultAddress, "address to connect to")
	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(2)
	}
	if err := mon(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func mon() error {
	g, err := gpsd.Dial(addrFlag)
	if err != nil {
		return err
	}
	defer g.Close()

	if err := g.Stream(gpsd.WATCH_ENABLE|gpsd.WATCH_JSON, ""); err != nil {
		return err
	}
	defer g.Stream(gpsd.WATCH_DISABLE, "")

	for v := range g.C() {
		switch t := v.(type) {
		case *gpsd.VERSION:
			fmt.Printf("GPSD Version: %s, Proto: %.0f.%.0f\n", t.Release, t.ProtoMajor, t.ProtoMinor)
		case *gpsd.DEVICES:
			if len(t.Devices) == 0 {
				return errors.New("no devices available")
			}
			fmt.Println("Available devices:")
			for _, d := range t.Devices {
				fmt.Printf("\t%s\n", d.Path)
			}
		case *gpsd.TPV:
			if t.Lat != 0 && t.Lon != 0 {
				fmt.Printf("%.9f %.9f\n", t.Lat, t.Lon)
			} else {
				fmt.Println("n/a")
			}
		}
	}
	return g.Err()
}
```

## Testing

GPSD address `TEST_ADDR` should be set to enable integrated testing, e.g. `TEST_ADDR=localhost:2947 go test`.

## Debugging

A positive value of `DEBUG` or `DEBUG_GPSD` enables debug output, unless you want to use a custom logger that can be set with `WithLogger(logger)` option.

## Contributing

Any contributions are welcome, just fork the repo and open a pull request or simply file an issue.
