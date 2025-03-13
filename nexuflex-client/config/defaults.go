// defaults.go
/**
 * Nexuflex Client - Default Configurations
 *
 * This file contains the default configurations for the client.
 *
 * @author msto63
 * @version 1.0.0
 * @date 2025-03-12
 */

package config

// GetDefaultConfig returns the default configuration for the client
func GetDefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Address:                "",
			Port:                   50051,
			UseTLS:                 false,
			DiscoveryToken:         "NEXUFLEX_DISCOVERY",
			AutoDiscover:           true,
			DiscoverTimeoutSeconds: 5,
		},
		UI: UIConfig{
			ColorScheme:           "default",
			HeaderText:            "nexuflex Terminal",
			ShowTimestamps:        true,
			EnableSounds:          false,
			MaxOutputLines:        1000,
			MaxHistoryEntries:     100,
			AutoCompleteEnabled:   true,
			AutoFillServicePrefix: true,
			Language:              "en",
		},
		Commands: CommandsConfig{
			SaveHistory:           true,
			UseLocalAliases:       true,
			MaxLocalAliases:       50,
			EnableMultilineInput:  true,
			SaveHistoryOnShutdown: true,
		},
	}
}
