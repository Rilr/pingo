package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"golang.org/x/crypto/ssh"
)

var devAddr = "10.100.10.1"  // This is where you'll SSH into
var tunAddr = "10.100.10.12" // This is the tunnel we're monitoring
var wanAddr = "1.1.1.1"      // This is the WAN address we're using to check connectivity

func postNewTicket() int {
	auth := ManageAuth()
	baseURL := "http://na.myconnectwise.net/v4_6_release/apis/3.0/service/tickets"
	jsonData := PostTicketPayload()
	u, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating the webrequest", err)
	}
	req.Header.Add("clientId", "3e53e6c4-d9ca-4916-8651-bc1e33e1c132")
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error doing the webrequest", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
	defer res.Body.Close()
	var ticket Ticket
	err = json.Unmarshal(body, &ticket)
	if err != nil {
		fmt.Println("Error decoding the data", err)
	}
	return ticket.ID
}

func PutTicketNote(ticketID int, note string) {
	fmt.Printf("Adding note to ticket %d: %s\n", ticketID, note)
	// Placeholder logic for adding a note to a ticket
}

func checkLogForTicket() (int, bool) {
	// Parse pingo.log for the latest ticket number
	f, err := os.Open("pingo.log")
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return 0, false
	}
	defer f.Close()

	buf := make([]byte, 4096)
	var content []byte
	for {
		n, err := f.Read(buf)
		if n > 0 {
			content = append(content, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading log file:", err)
			return 0, false
		}
	}

	lines := bytes.SplitSeq(content, []byte{'\n'})
	for line := range lines {
		// Truncate the first 20 characters (date portion) if line is long enough
		var trimmedLine []byte
		if len(line) > 20 {
			trimmedLine = line[20:]
		} else {
			trimmedLine = line
		}
		// Look for lines like: "Ticket created with ID: <number>"
		var id int
		n, _ := fmt.Sscanf(string(trimmedLine), "Ticket created with ID: %d", &id)
		if n == 1 {
			return id, true
		}
	}
	return 0, false
}

func checkManageForTicket(ticketID int) bool {
	// Placeholder logic for checking if a ticket exists in Manage
	// In a real implementation, this would query the Manage API
	fmt.Printf("Checking if ticket %d exists in Manage...\n", ticketID)
	return true // Assume the ticket exists for this example
}

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
	// AddtoLog(fmt.Sprintf("Packets sent: %d, Packets received: %d, RTT min/avg/max: %v/%v/%v", // Commented out for cleaner logging
	// 	stats.PacketsSent, stats.PacketsRecv, stats.MinRtt, stats.AvgRtt, stats.MaxRtt))
	switch {

	case stats.PacketsRecv == 0 || stats.MaxRtt == 0:
		testPassed = false

	case stats.PacketsSent > stats.PacketsRecv || stats.PacketLoss > 0:
		AddtoLog(fmt.Sprintf("Ping to address %s reveals packet loss at: %f%%", addr, stats.PacketLoss))
		testPassed = true

	default:
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

	for range 1 { // Wrapping in a for range loop to allow for termination or extension in the future
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
				AddtoLog(fmt.Sprintf("WAN address %s is reachable. Creating a ticket and restarting the tunnels...", wanAddr))
				i, b := checkLogForTicket()

				if b {
					AddtoLog(fmt.Sprintf("Ticket %d Present in Log. Checking it's status...", i))

					if checkManageForTicket(i) {
						AddtoLog(fmt.Sprintf("Ticket %d is already present in Manage. Adding a note and exiting.", i))
						PutTicketNote(i, "Tunnel is down. Host is attempting to restart the tunnel.")
						initTtyToHost()
					} else {
						AddtoLog(fmt.Sprintf("Ticket %d is not active in Manage. Creating a new ticket.", i))
						t := postNewTicket()
						AddtoLog(fmt.Sprintf("Ticket created with ID: %d", t))
						PutTicketNote(t, "Tunnel is down. Host is attempting to restart the tunnel.")
						initTtyToHost()
					}
				} else {
					t := postNewTicket()
					AddtoLog(fmt.Sprintf("Ticket created with ID: %d", t))
				}
				initTtyToHost()
			}
		} else {
			AddtoLog(fmt.Sprintf("Tunnel address %s is reachable. No action needed.", tunAddr))
		}
		time.Sleep(1 * time.Second) // Wait before the next iteration
	}
}
