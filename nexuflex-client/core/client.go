// client.go
/**
* Nexuflex Client - Client Implementation
*
* This file contains the main implementation of the nexuflex client,
* which handles communication with the Application Server.
*
* @author msto63
* @version 1.0.0
* @date 2025-03-12
 */

package core

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/msto63/nexuflex/shared/proto"
	"github.com/nexuflex/nexuflex-client/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// LogFunc defines the type for the logging function
type LogFunc func(format string, v ...interface{})

// Client represents the nexuflex client
type Client struct {
	// Configuration
	config *config.Config

	// Logger
	logger LogFunc

	// gRPC connection and client
	conn   *grpc.ClientConn
	client proto.NexuflexServiceClient

	// Session and status
	sessionToken    string
	serverInfo      *proto.ServerInfo
	lastServiceUsed string

	// Callbacks
	onStatusChanged  func(statusInfo *proto.StatusInfo)
	onServerList     func(servers []*proto.ServerInfo) (int, error)
	onOutputReceived func(output string)
}

// NewClient creates a new Client instance
func NewClient(cfg *config.Config, logger LogFunc) *Client {
	return &Client{
		config:          cfg,
		logger:          logger,
		sessionToken:    "",
		lastServiceUsed: "",
	}
}

// SetCallbacks sets the callback functions for UI updates
func (c *Client) SetCallbacks(
	onStatusChanged func(statusInfo *proto.StatusInfo),
	onServerList func(servers []*proto.ServerInfo) (int, error),
	onOutputReceived func(output string),
) {
	c.onStatusChanged = onStatusChanged
	c.onServerList = onServerList
	c.onOutputReceived = onOutputReceived
}

// DiscoverServer performs server discovery
func (c *Client) DiscoverServer(timeout time.Duration) error {
	c.logger("Starting server discovery...")

	// If already connected, close connection
	if c.conn != nil {
		c.Close()
	}

	// Perform server discovery (simulated for now)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// In a full implementation, this would send a UDP multicast
	// For this example, we simulate discovery with known servers
	knownServers := []*proto.ServerInfo{
		{
			Hostname:    "localhost",
			Address:     "localhost",
			Port:        50051,
			ShortName:   "Local Dev Server",
			Description: "Local development server",
			TlsEnabled:  false,
			Version:     "1.0.0",
		},
		{
			Hostname:    "remote-example",
			Address:     "remote-example.com",
			Port:        50051,
			ShortName:   "Remote Example",
			Description: "Example of a remote server",
			TlsEnabled:  true,
			Version:     "1.0.0",
		},
	}

	// Show server list to user, if callback is set
	if c.onServerList != nil {
		selectedIndex, err := c.onServerList(knownServers)
		if err != nil {
			return err
		}

		// Connect to selected server
		if selectedIndex >= 0 && selectedIndex < len(knownServers) {
			selectedServer := knownServers[selectedIndex]
			return c.Connect(selectedServer.Address, int(selectedServer.Port), selectedServer.TlsEnabled)
		}

		return fmt.Errorf("no server selection made")
	}

	// If no callback is set, connect to the first server
	if len(knownServers) > 0 {
		return c.Connect(knownServers[0].Address, int(knownServers[0].Port), knownServers[0].TlsEnabled)
	}

	return fmt.Errorf("no servers found")
}

// Connect establishes a connection to the server
func (c *Client) Connect(address string, port int, useTLS bool) error {
	c.logger("Connecting to %s:%d (TLS: %v)...", address, port, useTLS)

	// Close existing connection, if any
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
		c.client = nil
		c.sessionToken = ""
		c.serverInfo = nil
	}

	// Configure connection options
	var opts []grpc.DialOption
	if useTLS {
		// In a real implementation, TLS certificates would be configured here
		// For this example, we use standard TLS without certificate verification
		creds := credentials.NewClientTLSFromCert(nil, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Establish connection
	serverAddr := fmt.Sprintf("%s:%d", address, port)
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		c.logger("Connection error: %v", err)

		// Update status information
		if c.onStatusChanged != nil {
			c.onStatusChanged(&proto.StatusInfo{
				ConnectionStatus: proto.StatusInfo_CONNECTION_ERROR,
				SessionStatus:    proto.StatusInfo_NOT_LOGGED_IN,
			})
		}

		return fmt.Errorf("failed to connect to server: %v", err)
	}

	c.conn = conn
	c.client = proto.NewNexuflexServiceClient(conn)

	// Send Connect request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.Connect(ctx, &proto.ConnectRequest{
		Address: address,
		Port:    int32(port),
		UseTls:  useTLS,
	})
	if err != nil {
		c.conn.Close()
		c.conn = nil
		c.client = nil

		c.logger("Connect request failed: %v", err)

		// Update status information
		if c.onStatusChanged != nil {
			c.onStatusChanged(&proto.StatusInfo{
				ConnectionStatus: proto.StatusInfo_CONNECTION_ERROR,
				SessionStatus:    proto.StatusInfo_NOT_LOGGED_IN,
			})
		}

		return fmt.Errorf("connect request failed: %v", err)
	}

	if !resp.Success {
		c.conn.Close()
		c.conn = nil
		c.client = nil

		c.logger("Connect failed: %s", resp.ErrorMessage)

		// Update status information
		if c.onStatusChanged != nil {
			c.onStatusChanged(&proto.StatusInfo{
				ConnectionStatus: proto.StatusInfo_CONNECTION_ERROR,
				SessionStatus:    proto.StatusInfo_NOT_LOGGED_IN,
			})
		}

		return fmt.Errorf("connect failed: %s", resp.ErrorMessage)
	}

	// Store server information
	c.serverInfo = &proto.ServerInfo{
		Address:    address,
		Port:       int32(port),
		ShortName:  resp.ServerName,
		Version:    resp.Version,
		TlsEnabled: useTLS,
	}

	c.logger("Connected to server %s (Version %s)", resp.ServerName, resp.Version)

	// Report status
	if c.onStatusChanged != nil {
		c.onStatusChanged(&proto.StatusInfo{
			ConnectionStatus: proto.StatusInfo_CONNECTED,
			SessionStatus:    proto.StatusInfo_NOT_LOGGED_IN,
			ServerName:       resp.ServerName,
		})
	}

	return nil
}

// Login performs user authentication
func (c *Client) Login(username, password string) error {
	if c.client == nil {
		return fmt.Errorf("not connected to server")
	}

	c.logger("Login for user %s...", username)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.Login(ctx, &proto.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		c.logger("Login request failed: %v", err)
		return fmt.Errorf("login request failed: %v", err)
	}

	if !resp.Success {
		c.logger("Login failed: %s", resp.ErrorMessage)
		return fmt.Errorf("login failed: %s", resp.ErrorMessage)
	}

	// Store session token and user information
	c.sessionToken = resp.SessionToken
	c.logger("Login successful for %s", resp.UserInfo.DisplayName)

	// Report status
	if c.onStatusChanged != nil {
		c.onStatusChanged(&proto.StatusInfo{
			ConnectionStatus: proto.StatusInfo_CONNECTED,
			SessionStatus:    proto.StatusInfo_AUTHENTICATED,
			ServerName:       c.serverInfo.ShortName,
			Username:         username,
		})
	}

	// Output welcome message
	if c.onOutputReceived != nil {
		c.onOutputReceived(fmt.Sprintf("Welcome, %s! You are now logged in.", resp.UserInfo.DisplayName))
	}

	return nil
}

// Logout logs out the user
func (c *Client) Logout() error {
	if c.client == nil {
		return fmt.Errorf("not connected to server")
	}

	if c.sessionToken == "" {
		return fmt.Errorf("not logged in")
	}

	c.logger("Logout...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.Logout(ctx, &proto.LogoutRequest{
		SessionToken: c.sessionToken,
	})
	if err != nil {
		c.logger("Logout request failed: %v", err)
		return fmt.Errorf("logout request failed: %v", err)
	}

	if !resp.Success {
		c.logger("Logout failed: %s", resp.ErrorMessage)
		return fmt.Errorf("logout failed: %s", resp.ErrorMessage)
	}

	// Reset session token
	c.sessionToken = ""
	c.logger("Logout successful")

	// Report status
	if c.onStatusChanged != nil {
		c.onStatusChanged(&proto.StatusInfo{
			ConnectionStatus: proto.StatusInfo_CONNECTED,
			SessionStatus:    proto.StatusInfo_NOT_LOGGED_IN,
			ServerName:       c.serverInfo.ShortName,
		})
	}

	return nil
}

// ExecuteCommand executes a command on the server
func (c *Client) ExecuteCommand(command string) error {
	if c.client == nil {
		return fmt.Errorf("not connected to server")
	}

	c.logger("Executing command: %s", command)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.client.ExecuteCommand(ctx, &proto.CommandRequest{
		SessionToken: c.sessionToken,
		CommandLine:  command,
		LastContext:  c.lastServiceUsed,
	})
	if err != nil {
		c.logger("Command execution failed: %v", err)
		return fmt.Errorf("command execution failed: %v", err)
	}

	// Process output
	if !resp.Success {
		c.logger("Command failed: %s", resp.ErrorMessage)
		if c.onOutputReceived != nil {
			c.onOutputReceived(fmt.Sprintf("Error: %s", resp.ErrorMessage))
		}
	} else {
		if c.onOutputReceived != nil {
			c.onOutputReceived(resp.Output)
		}

		// Remember last used service
		if resp.NewContext != "" {
			c.lastServiceUsed = resp.NewContext
			c.logger("New service context: %s", c.lastServiceUsed)
		}
	}

	// Display status message
	if c.onStatusChanged != nil {
		c.onStatusChanged(resp.StatusInfo)
	}

	return nil
}

// ExecuteStreamingCommand executes a command that produces continuous output
func (c *Client) ExecuteStreamingCommand(command string) error {
	if c.client == nil {
		return fmt.Errorf("not connected to server")
	}

	c.logger("Executing streaming command: %s", command)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	stream, err := c.client.ExecuteStreamingCommand(ctx, &proto.CommandRequest{
		SessionToken: c.sessionToken,
		CommandLine:  command,
		LastContext:  c.lastServiceUsed,
	})
	if err != nil {
		c.logger("Streaming command execution failed: %v", err)
		return fmt.Errorf("streaming command execution failed: %v", err)
	}

	// Process stream
	for {
		output, err := stream.Recv()
		if err == io.EOF {
			// Stream ended
			c.logger("Streaming command completed")
			break
		}
		if err != nil {
			c.logger("Error receiving streaming data: %v", err)
			return fmt.Errorf("error receiving streaming data: %v", err)
		}

		// Process output by type
		switch output.Type {
		case proto.CommandOutput_TEXT:
			if c.onOutputReceived != nil {
				c.onOutputReceived(output.Content)
			}
		case proto.CommandOutput_STATUS_UPDATE:
			// Process status update (e.g., progress indicator)
			c.logger("Status update: %s (%d%%)", output.Content, output.ProgressPercent)
		case proto.CommandOutput_ERROR:
			c.logger("Streaming error: %s", output.Content)
			if c.onOutputReceived != nil {
				c.onOutputReceived(fmt.Sprintf("Error: %s", output.Content))
			}
		case proto.CommandOutput_COMPLETION:
			c.logger("Streaming command complete: %s", output.Content)
			if c.onOutputReceived != nil {
				c.onOutputReceived(fmt.Sprintf("Completed: %s", output.Content))
			}
		}
	}

	return nil
}

// AutoComplete provides command completion suggestions
func (c *Client) AutoComplete(partialInput string, cursorPos int) ([]string, string, error) {
	if c.client == nil {
		return nil, "", fmt.Errorf("not connected to server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resp, err := c.client.AutoComplete(ctx, &proto.AutoCompleteRequest{
		SessionToken:   c.sessionToken,
		PartialInput:   partialInput,
		CurrentContext: c.lastServiceUsed,
		CursorPosition: int32(cursorPos),
	})
	if err != nil {
		c.logger("Auto-completion failed: %v", err)
		return nil, "", fmt.Errorf("auto-completion failed: %v", err)
	}

	return resp.Suggestions, resp.CommonPrefix, nil
}

// GetAliases retrieves the available command aliases
func (c *Client) GetAliases() ([]*proto.AliasInfo, error) {
	if c.client == nil {
		return nil, fmt.Errorf("not connected to server")
	}

	if c.sessionToken == "" {
		return nil, fmt.Errorf("not logged in")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.GetAliases(ctx, &proto.GetAliasesRequest{
		SessionToken: c.sessionToken,
	})
	if err != nil {
		c.logger("Error retrieving aliases: %v", err)
		return nil, fmt.Errorf("error retrieving aliases: %v", err)
	}

	return resp.Aliases, nil
}

// CreateAlias creates a new command alias
func (c *Client) CreateAlias(alias, expandedCommand string) error {
	if c.client == nil {
		return fmt.Errorf("not connected to server")
	}

	if c.sessionToken == "" {
		return fmt.Errorf("not logged in")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.CreateAlias(ctx, &proto.CreateAliasRequest{
		SessionToken:    c.sessionToken,
		Alias:           alias,
		ExpandedCommand: expandedCommand,
	})
	if err != nil {
		c.logger("Error creating alias: %v", err)
		return fmt.Errorf("error creating alias: %v", err)
	}

	if !resp.Success {
		c.logger("Alias creation failed: %s", resp.ErrorMessage)
		return fmt.Errorf("alias creation failed: %s", resp.ErrorMessage)
	}

	c.logger("Alias '%s' created for '%s'", alias, expandedCommand)
	return nil
}

// DeleteAlias deletes a command alias
func (c *Client) DeleteAlias(alias string) error {
	if c.client == nil {
		return fmt.Errorf("not connected to server")
	}

	if c.sessionToken == "" {
		return fmt.Errorf("not logged in")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.DeleteAlias(ctx, &proto.DeleteAliasRequest{
		SessionToken: c.sessionToken,
		Alias:        alias,
	})
	if err != nil {
		c.logger("Error deleting alias: %v", err)
		return fmt.Errorf("error deleting alias: %v", err)
	}

	if !resp.Success {
		c.logger("Alias deletion failed: %s", resp.ErrorMessage)
		return fmt.Errorf("alias deletion failed: %s", resp.ErrorMessage)
	}

	c.logger("Alias '%s' deleted", alias)
	return nil
}

// GetAvailableServices retrieves the available services
func (c *Client) GetAvailableServices() ([]*proto.ServiceInfo, error) {
	if c.client == nil {
		return nil, fmt.Errorf("not connected to server")
	}

	if c.sessionToken == "" {
		return nil, fmt.Errorf("not logged in")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.GetAvailableServices(ctx, &proto.ServicesRequest{
		SessionToken: c.sessionToken,
	})
	if err != nil {
		c.logger("Error retrieving services: %v", err)
		return nil, fmt.Errorf("error retrieving services: %v", err)
	}

	return resp.Services, nil
}

// GetServiceCommands retrieves the available commands for a service
func (c *Client) GetServiceCommands(serviceName string) ([]*proto.CommandInfo, error) {
	if c.client == nil {
		return nil, fmt.Errorf("not connected to server")
	}

	if c.sessionToken == "" {
		return nil, fmt.Errorf("not logged in")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.GetServiceCommands(ctx, &proto.ServiceCommandsRequest{
		SessionToken: c.sessionToken,
		ServiceName:  serviceName,
	})
	if err != nil {
		c.logger("Error retrieving commands: %v", err)
		return nil, fmt.Errorf("error retrieving commands: %v", err)
	}

	return resp.Commands, nil
}

// GetCommandHelp retrieves help for a specific command
func (c *Client) GetCommandHelp(service, action, subaction string) (string, *proto.CommandInfo, error) {
	if c.client == nil {
		return "", nil, fmt.Errorf("not connected to server")
	}

	if c.sessionToken == "" {
		return "", nil, fmt.Errorf("not logged in")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.GetCommandHelp(ctx, &proto.CommandHelpRequest{
		SessionToken: c.sessionToken,
		Service:      service,
		Action:       action,
		Subaction:    subaction,
	})
	if err != nil {
		c.logger("Error retrieving help: %v", err)
		return "", nil, fmt.Errorf("error retrieving help: %v", err)
	}

	return resp.HelpText, resp.CommandInfo, nil
}

// IsConnected returns whether the client is connected to a server
func (c *Client) IsConnected() bool {
	return c.conn != nil && c.client != nil
}

// IsLoggedIn returns whether the client is logged in
func (c *Client) IsLoggedIn() bool {
	return c.sessionToken != ""
}

// GetServerInfo returns information about the connected server
func (c *Client) GetServerInfo() *proto.ServerInfo {
	return c.serverInfo
}

// GetLastServiceUsed returns the last used service
func (c *Client) GetLastServiceUsed() string {
	return c.lastServiceUsed
}

// SetLastServiceUsed sets the last used service
func (c *Client) SetLastServiceUsed(service string) {
	c.lastServiceUsed = service
}

// StartKeepAlive starts a background process for session keep-alive
func (c *Client) StartKeepAlive(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if c.client != nil && c.sessionToken != "" {
					ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
					resp, err := c.client.KeepAlive(ctx, &proto.KeepAliveRequest{
						SessionToken: c.sessionToken,
					})
					cancel()

					if err != nil {
						c.logger("KeepAlive error: %v", err)
					} else if !resp.SessionValid {
						c.logger("Session expired")
						c.sessionToken = ""

						// Report status
						if c.onStatusChanged != nil {
							c.onStatusChanged(&proto.StatusInfo{
								ConnectionStatus: proto.StatusInfo_CONNECTED,
								SessionStatus:    proto.StatusInfo_SESSION_EXPIRED,
								ServerName:       c.serverInfo.ShortName,
							})
						}

						// End KeepAlive since session has expired
						return
					}
				} else {
					// End KeepAlive if not connected or not logged in
					return
				}
			}
		}
	}()
}

// Close closes the connection to the server
func (c *Client) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		c.client = nil
		c.sessionToken = ""
		c.serverInfo = nil

		return err
	}
	return nil
}
