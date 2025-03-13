// discovery.go
/**
* Nexuflex Client - Server Discovery Implementation
*
* This file contains the implementation of the discovery mechanism
* for automatic detection of nexuflex servers on the network.
*
* @author msto63
* @version 1.0.0
* @date 2025-03-12
 */

package core

import (
	"fmt"
	"log"
	"net"
	"time"
)

// UDP multicast addresses for server discovery
const (
	DefaultMulticastAddress = "239.0.0.1:5000"
	DiscoveryTimeout        = 5 * time.Second
	DiscoveryPacketSize     = 1024
)

// DiscoveryPacket represents a discovery packet
type DiscoveryPacket struct {
	Type    string `json:"type"`    // "request" or "response"
	Token   string `json:"token"`   // Security token
	Address string `json:"address"` // Server address (only for "response")
	Port    int    `json:"port"`    // Server port (only for "response")
	Name    string `json:"name"`    // Server name (only for "response")
	Version string `json:"version"` // Server version (only for "response")
}

// PerformMulticastDiscovery performs a multicast discovery
// In a complete implementation, this function would be used
// to discover servers on the network. For this example, we simulate it.
func PerformMulticastDiscovery(multicastAddr, discoveryToken string, timeout time.Duration) error {
	// Parse multicast address
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		return fmt.Errorf("invalid multicast address: %v", err)
	}

	// Create UDP socket
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return fmt.Errorf("error creating UDP socket: %v", err)
	}
	defer conn.Close()

	// Set timeout
	err = conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return fmt.Errorf("error setting timeout: %v", err)
	}

	// Create discovery packet
	// In a real implementation, a JSON-encoded packet would be sent
	discoveryMessage := fmt.Sprintf("NEXUFLEX_DISCOVERY:%s", discoveryToken)

	// Send discovery packet
	_, err = conn.WriteToUDP([]byte(discoveryMessage), addr)
	if err != nil {
		return fmt.Errorf("error sending discovery packet: %v", err)
	}

	// Wait for responses
	buffer := make([]byte, DiscoveryPacketSize)
	servers := make(map[string]string) // Address -> Name

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			// If timeout reached, exit
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}
			log.Printf("Error receiving discovery response: %v", err)
			continue
		}

		// Process response
		// In a real implementation, a JSON-encoded packet would be received
		response := string(buffer[:n])
		log.Printf("Response from %s: %s", remoteAddr, response)

		// Add server to list
		servers[remoteAddr.String()] = response
	}

	// Output results
	log.Printf("Servers found: %d", len(servers))
	for addr, name := range servers {
		log.Printf("  %s: %s", addr, name)
	}

	return nil
}
