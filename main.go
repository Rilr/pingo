package main

import (
	"fmt"
	"os"
	"time"

	p "github.com/prometheus-community/pro-bing"
)

func BadPing() {
	if TestISP() { // if the ISP responds, tunnel is probably down
		CreateTicket()
		TtyIntoHost()
	} else {
		AddtoLog() // if the ISP does not respond, we can't create a ticket nor get the tunnel restored until the ISP is back online
		os.Exit(11)
	}
}

func TestISP() bool {
	AddtoLog()
	return true
}

func CreateTicket() int{
	AddtoLog()
	ticketNumber := 123456 // Simulated ticket number
	return ticketNumber
}

func TtyIntoHost() {
	fmt.Println("This function is not used in the main program.")
}

func ElevateTicket(t int) error{
	if t >= 0 {
		fmt.Print(t)
		return nil
	} else {
		return fmt.Errorf("invalid ticket number: %d", t)
	}
}

func CloseTicket(t int) error{
	if t >= 0 {
		fmt.Print(t)
		return nil
	} else {
		return fmt.Errorf("invalid ticket number: %d", t)
	}
}

func AddtoLog() {
	fmt.Println("This function is not used in the main program.")
}

func main() {
	address := "8.8.8.8"
		for i := 0; i < 2; i++ {
		pinger, err := p.NewPinger(address)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create pinger: %v\n", err)
		}
		pinger.SetPrivileged(true) // Allows to run without sudo
		pinger.Count = 10
		err = pinger.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ping error: %v\n", err)
			continue
		}
		stats := pinger.Statistics()
		switch{
			case stats.PacketLoss >= 20:
				fmt.Printf("%s: Ping to %s resulted in %f%% packet loss.\n", time.Now().Format("2006-01-02_15:04:05"), address, stats.PacketLoss)

			case stats.PacketLoss < 20 && stats.PacketLoss > 0:
				fmt.Printf("Ping to %s resulted in %f%% packet loss.\n Average RTT: %s. Deviation: %s\n", address, stats.PacketLoss, stats.AvgRtt.String(), stats.StdDevRtt.String())
			default:
				fmt.Printf("Ping to %s resulted in zero packet loss.\n Average RTT: %s. Deviation: %s\n", address, stats.AvgRtt.String(), stats.StdDevRtt.String())
		}
	}
}
