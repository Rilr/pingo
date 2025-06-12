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
	"pingo/static"
	"strconv"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"golang.org/x/crypto/ssh"
)

var devAddr = static.Addr.Dev // This is where you'll SSH into if the tunnel is down
var tunAddr = static.Addr.Tun // This is the tunnel we're monitoring
var wanAddr = static.Addr.Wan // This is the WAN address we're using to check connectivity

// checkManageForTicket checks the status of a ticket in ConnectWise Manage and returns true if the ticket is still valid (not closed).
func checkManageForTicket(ticketID int) bool {
	auth := ManageAuth()
	var ticketValid bool
	// Convert the Int to a string
	ticketNumberStr := strconv.Itoa(ticketID)
	// Create the webrequest
	baseURL := "http://na.myconnectwise.net/v4_6_release/apis/3.0/service/tickets/" + ticketNumberStr
	params := url.Values{}
	u, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
	}
	u.RawQuery = params.Encode()
	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		fmt.Println("Error creating the webrequest", err)
	}
	req.Header.Add("clientId", "3e53e6c4-d9ca-4916-8651-bc1e33e1c132")
	req.Header.Add("Authorization", "Basic "+auth)
	// Do the webrequest
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error doing the webrequest", err)
	}
	defer res.Body.Close()
	// Handle the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading the body", err)
	}
	var ticketData Ticket
	err = json.Unmarshal(body, &ticketData)
	if err != nil {
		fmt.Println("Error decoding the data", err)
	}
	switch ticketData.Status.ID {
	case 736, 612, 452, 737, 739, 778, 17, 80, 9: // >Completed(QA Review), >QA Reviewed Closed/No Response, >QA Reviewed/Closed etc...
		ticketValid = false
	default:
		ticketValid = true
	}
	AddtoLog(fmt.Sprintf("Ticket %d status: %s (ID: %d)\n", ticketID, ticketData.Status.Name, ticketData.Status.ID))
	return ticketValid // If the ticket is valid, we won't create a new one. If it's been closed (which returns false), we will create a new one.
}

// postNewTicket creates a new ticket in ConnectWise Manage and returns the ticket ID.
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

// TODO: Implement the PostTicketPayload function to return the JSON payload for creating a new ticket
func putTicketNote(ticketID int, note string) {
	// fmt.Printf("PutTicketNote Ran with the inputs of %d and %s\n", ticketID, note)
	// Placeholder logic for adding a note to a ticket
}

// checkLogForTicket checks the pingo.log file for the latest ticket number and returns the ticket ID and a boolean indicating if a ticket was found.
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

// sshIntoHost connects to a host via SSH and executes a command after InitTtyToHost is called to check if the host is reachable first.
func sshIntoHost(addr, user, pass, cmd string) error {
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
	outputStr := string(bytes.TrimSpace(output)) // Trim whitespace/newlines
	if err != nil {
		AddtoLog(fmt.Sprintf("SSH command error: %v", err))
		AddtoLog(fmt.Sprintf("SSH command output: %s", outputStr))
		return err
	}
	AddtoLog(fmt.Sprintf("SSH command output: %s", string(output)))
	return nil
}

// TestAddress pings the specified address and returns true if the address is reachable.
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

// Logging function to write messages to pingo.log
func AddtoLog(s string) {
	f, err := os.OpenFile("pingo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	// fmt.Println(s) // Remove in production, testing purposes only
	logger.Println(s)
}

// Tests the tunnel constantly
// If the tunnel is down, it will check the WAN address
// If the WAN address is down, it will check the device address (if all three are down, the device is most likely disconnected from the network)
// If the tunnel is down, but the WAN address is up, it will attempt to recover and submit a ticket
// After restarting the tunnel and submitting a ticket, it should check if the tunnel is up again
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
				AddtoLog(fmt.Sprintf("WAN address %s is reachable. Checking for an open ticket and restarting the tunnels...", wanAddr))
				i, b := checkLogForTicket()

				if b {
					AddtoLog(fmt.Sprintf("Ticket %d Present in Log. Checking it's validity via it's status ID...", i))

					if checkManageForTicket(i) {
						AddtoLog(fmt.Sprintf("Ticket %d is valid ticket. Adding a note and exiting.", i))
						putTicketNote(i, "Tunnel is down. Host is attempting to restart the tunnel.")
						InitTtyToHost()
					} else {
						AddtoLog(fmt.Sprintf("Ticket %d is not active in Manage. Creating a new ticket.", i))
						t := postNewTicket()
						AddtoLog(fmt.Sprintf("Ticket created with ID: %d", t))
						putTicketNote(t, "Tunnel is down. Host is attempting to restart the tunnel.")
						InitTtyToHost()
					}
				} else {
					t := postNewTicket()
					AddtoLog(fmt.Sprintf("Ticket created with ID: %d", t))
				}
				InitTtyToHost()
			}
		} else {
			AddtoLog(fmt.Sprintf("Tunnel address %s is reachable. No action needed.", tunAddr))
			fmt.Println(ManageAuth())

		}
		time.Sleep(1 * time.Second) // Wait before the next iteration
	}
}
