package main

import (
	"fmt"
	"os"
	"time"

	p "github.com/prometheus-community/pro-bing"
)

func main() {
	address := "8.8.8.8"
	interval := time.Minute

	for {
		success := false
		for i := 0; i < 5; i++ {
			pinger, err := p.NewPinger(address)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create pinger: %v\n", err)
				continue
			}
			pinger.SetPrivileged(true)
			pinger.Count = 1
			pinger.Timeout = time.Second * 2
			err = pinger.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Ping error: %v\n", err)
				continue
			}
			stats := pinger.Statistics()
			if stats.PacketsRecv > 0 {
				success = true
				fmt.Printf("Ping %d: reply from %s\n", i+1, address)
			} else {
				fmt.Printf("Ping %d: no reply from %s\n", i+1, address)
			}
			time.Sleep(2 * time.Second)
		}
		if success {
			fmt.Printf("Result: At least one reply from %s in this minute.\n\n", address)
		} else {
			fmt.Printf("Result: No replies from %s in this minute.\n\n", address)
		}
		time.Sleep(interval - 10*time.Second) // 5 pings x 2s = 10s, so wait the rest of the minute
	}
}
