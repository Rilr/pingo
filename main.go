package main

import (
	"fmt"
	"log"
	"os"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"golang.org/x/crypto/ssh"
)

var devAddr = "10.100.10.1"   // This is where you'll SSH into
var tunAddr = "10.100.10.220" // This is the tunnel we're monitoring
var wanAddr = "1.1.1.1"       // This is the WAN address we're using to check connectivity

func SshIntoHost(addr, user, pass, cmd string) error {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.KeyboardInteractive(
				func(user, instruction string, questions []string, echos []bool) ([]string, error) {
					answers := make([]string, len(questions))
					for i := range questions {
						answers[i] = pass
					}
					return answers, nil
				},
			),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", addr+":22", config)
	if err != nil {
		AddtoLog(fmt.Sprintf("SSH connection failed: %v", err))
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		AddtoLog(fmt.Sprintf("Failed to create SSH session: %v", err))
		return err
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		AddtoLog(fmt.Sprintf("Failed to run command: %v", err))
		return err
	}
	AddtoLog(fmt.Sprintf("SSH command output: %s", string(output)))
	return nil
}

func initTtyToHost() {
	if !TestAddress(devAddr, 2, 1*time.Second, 10*time.Second) {
		AddtoLog(fmt.Sprintf("Device address %s is unresponsive", devAddr))
		os.Exit(3)
	} else {
		AddtoLog(fmt.Sprintf("Attempting to Tunnel into: %s", devAddr))
		if err := SshIntoHost(devAddr, "root", "pass", "ipsec restart"); err != nil {
			os.Exit(4)
		} else {
			AddtoLog(fmt.Sprintf("Command ran successfully on device address %s", devAddr))
			os.Exit(0)
		}
	}
}

func TestAddress(addr string, count int, interval time.Duration, timeout time.Duration) bool {
	var testPassed bool
	pinger, err := probing.NewPinger(addr)
	if err != nil {
		panic(err)
	}

	pinger.SetPrivileged(true)
	pinger.Count = count
	pinger.Interval = interval
	pinger.Timeout = timeout
	err = pinger.Run()
	if err != nil {
		panic(err)
	}

	stats := pinger.Statistics() // get send/receive/rtt stats
	AddtoLog(fmt.Sprintf("Packets sent: %d, Packets received: %d, RTT min/avg/max: %v/%v/%v",
		stats.PacketsSent, stats.PacketsRecv, stats.MinRtt, stats.AvgRtt, stats.MaxRtt))
	switch {

	case stats.PacketsRecv == 0 || stats.MaxRtt == 0:
		fmt.Printf("No packets received from %s, address is unreachable.\n", addr)
		testPassed = false

	case stats.PacketLoss > 0:
		fmt.Printf("Some packets were lost when pinging %s, packet loss: %.2f%%.\n", addr, stats.PacketLoss)
		testPassed = true

	default:
		fmt.Printf("Successfully pinged %s, average RTT: %v.\n", addr, stats.AvgRtt)
		testPassed = true
	}
	return testPassed
}

func AddtoLog(s string) {
	f, err := os.OpenFile("pingo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	fmt.Println(s) // Remove in production, testing purposes only
	logger.Println(s)
}

func main() {
	i := 1 * time.Second  // Interval is the wait time between each packet send. Default is 1s.
	t := 20 * time.Second // Timeout specifies a timeout before ping exits, regardless of how many packets have been received.
	c := 5                // Count tells pinger to stop after sending (and receiving) 'c' echo packets. If this option is not specified, pinger will operate until interrupted.

	for range 1 {
		if !TestAddress(tunAddr, c, i, t) {
			AddtoLog(fmt.Sprintf("Tunnel address %s is unreachable. Testing %s", tunAddr, wanAddr))

			if !TestAddress(wanAddr, c, i, t) {
				AddtoLog(fmt.Sprintf("WAN address %s is also unreachable. Testing %s", wanAddr, devAddr))

				if !TestAddress(devAddr, c, i, t) {
					AddtoLog(fmt.Sprintf("Device address %s is unreachable. Host is most likely disconnected from the network.", devAddr))
					os.Exit(1)
				} else {
					AddtoLog(fmt.Sprintf("Device address %s is reachable. Host is connected to network with no WAN connection.", devAddr))
					os.Exit(2)
				}

			} else {
				AddtoLog(fmt.Sprintf("WAN address %s is reachable. Restarting tunnel...", wanAddr))
				initTtyToHost()
			}
		} else {
			AddtoLog(fmt.Sprintf("Tunnel address %s is reachable. No action needed.", tunAddr))
		}
		time.Sleep(1 * time.Second) // Wait before the next iteration
	}
}
