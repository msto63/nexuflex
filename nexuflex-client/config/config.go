// config.go
/**
 * Nexuflex Client - Configuration Management
 *
 * This file contains the data structures and functions for managing
 * the client configuration.
 *
 * @author msto63
 * @version 1.0.0
 * @date 2025-03-12
 */

package config

import (
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// Config represents the overall configuration of the client
type Config struct {
	Server   ServerConfig   `ini:"server"`
	UI       UIConfig       `ini:"ui"`
	Commands CommandsConfig `ini:"commands"`
}

// ServerConfig contains the configuration for the server connection
type ServerConfig struct {
	Address                string `ini:"address"`
	Port                   int    `ini:"port"`
	UseTLS                 bool   `ini:"use_tls"`
	DiscoveryToken         string `ini:"discovery_token"`
	AutoDiscover           bool   `ini:"auto_discover"`
	DiscoverTimeoutSeconds int    `ini:"discover_timeout_seconds"`
}

// UIConfig contains configuration options for the user interface
type UIConfig struct {
	ColorScheme           string `ini:"color_scheme"`
	HeaderText            string `ini:"header_text"`
	ShowTimestamps        bool   `ini:"show_timestamps"`
	EnableSounds          bool   `ini:"enable_sounds"`
	MaxOutputLines        int    `ini:"max_output_lines"`
	MaxHistoryEntries     int    `ini:"max_history_entries"`
	AutoCompleteEnabled   bool   `ini:"auto_complete_enabled"`
	AutoFillServicePrefix bool   `ini:"auto_fill_service_prefix"`
	Language              string `ini:"language"`
}

// CommandsConfig contains configuration options for command processing
type CommandsConfig struct {
	SaveHistory           bool `ini:"save_history"`
	UseLocalAliases       bool `ini:"use_local_aliases"`
	MaxLocalAliases       int  `ini:"max_local_aliases"`
	EnableMultilineInput  bool `ini:"enable_multiline_input"`
	SaveHistoryOnShutdown bool `ini:"save_history_on_shutdown"`
}

// LoadConfig loads the configuration from a file
func LoadConfig(configPath string) (Config, error) {
	// Default configuration as base
	config := GetDefaultConfig()

	// If no path is specified, try standard paths
	if configPath == "" {
		// Determine user's directory
		userConfigDir, err := os.UserConfigDir()
		if err == nil {
			// First try the configuration file in the user directory
			configPath = filepath.Join(userConfigDir, "nexuflex", "client.ini")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				// Try alternative: configuration file in current directory
				configPath = "client.ini"
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					// No configuration file found, use default configuration
					return config, nil
				}
			}
		} else {
			// Error determining user directory, use default configuration
			return config, nil
		}
	}

	// Load configuration file
	cfg, err := ini.Load(configPath)
	if err != nil {
		// If the file cannot be loaded, use default configuration
		return config, err
	}

	// Map configuration to structure
	err = cfg.MapTo(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config Config, configPath string) error {
	// If no path is specified, use default path
	if configPath == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		// Ensure directory exists
		configDir := filepath.Join(userConfigDir, "nexuflex")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
		configPath = filepath.Join(configDir, "client.ini")
	}

	// Create new .ini file
	cfg := ini.Empty()

	// Write configuration to .ini file
	err := ini.ReflectFrom(cfg, &config)
	if err != nil {
		return err
	}

	// Save file
	return cfg.SaveTo(configPath)
}
