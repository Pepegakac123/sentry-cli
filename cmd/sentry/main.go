package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	interval := flag.Duration("interval", 2*time.Second, "Check interval (e.g. 1s, 500ms)")
	flag.Parse()

	fmt.Printf("The sentry-cli starts with interval of %s\n", *interval)
}
