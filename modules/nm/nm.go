package nm

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/highercomve/gohtmx/modules/nm/nmmodules"
)

type NetworkManager struct{}

func Init() nmmodules.WifiManager {
	return &NetworkManager{}
}

func (nm *NetworkManager) List() (conns []nmmodules.WifiConn, err error) {
	// Run the nmcli command to list available WiFi networks, filtering for the necessary fields
	cmd := exec.Command(
		"nmcli",
		"-f",
		"ssid,mode,freq,signal,active,security",
		"device",
		"wifi",
		"list",
	)

	// Get the command output (stdout) and check for errors
	stdout, err := cmd.Output()
	if err != nil {
		// If there is an error running the nmcli command, log it and return
		fmt.Println("Error running nmcli command:", err)
		return
	}

	// Initialize a scanner to read the command output line by line
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))

	// Skip the first line (header) as it contains the column titles
	scanner.Scan()

	// Initialize a map to store the WiFi connections, keyed by SSID
	connections := map[string]nmmodules.WifiConn{}

	// Define a regex pattern to extract SSID, mode, frequency, signal strength, active status, and security
	wifiRegex := regexp.MustCompile(`(?P<SSID>.+?)\s{2,}(?P<MODE>\S+)\s+(?P<FREQ>\S+\sMHz)\s+(?P<SIGNAL>\d+)\s+(?P<ACTIVE>\S+)\s+(?P<SECURITY>.+)`)

	// Process each line of the command output
	for scanner.Scan() {
		line := scanner.Text()

		// Apply the regex to extract the WiFi details from the current line
		matches := wifiRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			// Create a new WifiConn struct and populate it with extracted data
			conn := nmmodules.WifiConn{}
			conn.ID = strings.TrimSpace(matches[1])
			conn.SSID = strings.TrimSpace(matches[1])  // SSID (WiFi network name)
			conn.Mode = matches[2]                     // Mode (e.g., Infra)
			conn.Frequency = matches[3]                // Frequency (e.g., 2442 MHz)
			conn.Active = matches[5] == "yes"          // Active status (true if 'yes')
			conn.Security = strings.Fields(matches[6]) // Security protocols (e.g., WPA1, WPA2)

			// Convert the signal strength to an integer
			conn.Strength, err = strconv.Atoi(matches[4])
			if err != nil {
				// If conversion fails, return the error
				return conns, err
			}

			// Check if this SSID is already present in the connections map
			if c, ok := connections[conn.ID]; !ok {
				// If it's not in the map, add the new connection
				connections[conn.ID] = conn
			} else {
				// If it's already in the map, check the active status
				// If the current connection is active and the one in the map is not, update it
				if !c.Active && conn.Active {
					connections[conn.ID] = conn
				}
			}
		}
	}

	// Handle any scanning errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading command output:", err)
		return conns, err
	}

	// Convert the map to a slice of WifiConn structs
	for _, v := range connections {
		conns = append(conns, v)
	}

	// Return the list of connections and any error encountered
	slices.SortFunc(conns, func(a nmmodules.WifiConn, b nmmodules.WifiConn) int {
		return b.Strength - a.Strength
	})
	return conns, err
}

func (nm *NetworkManager) Save(ssid, password string) error {
	// Use nmcli to connect to the network
	cmd := exec.Command("nmcli", "device", "wifi", "connect", ssid, "password", password)

	// Capture both the output and errors separately
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		// If there's an error, return it as a Go error, including nmcli's stderr
		return fmt.Errorf("Failed to connect to network %s: %s, nmcli error: %s", ssid, err.Error(), stderr.String())
	}

	// Optionally log the success message from stdout
	fmt.Println("nmcli output:", stdout.String())

	return nil
}
