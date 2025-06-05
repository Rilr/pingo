package main

import (
	"fmt"
	"github.com/prometheus-community/pro-bing"
)

pinger, err := probing.NewPinger("www.google.com")
if err != nil {
	panic(err)
}
pinger.Count = 3
err = pinger.Run() // Blocks until finished.
if err != nil {
	panic(err)
}
stats := pinger.Statistics()