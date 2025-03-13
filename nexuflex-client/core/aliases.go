// aliases.go
/**
 * Nexuflex Client - Local Alias Management
 *
 * This file contains functions for managing local command aliases
 * that are not stored on the server.
 *
 * @author msto63
 * @version 1.0.0
 * @date 2025-03-12
 */

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AliasManager manages local command aliases
type AliasManager struct {
	aliases  map[string]string
	maxCount int
}

// NewAliasManager creates a new AliasManager
func NewAliasManager(maxCount int) *AliasManager {
	return &AliasManager{
		aliases:  make(map[string]string),
		maxCount: maxCount,
	}
}

// AddAlias adds a local alias
func (am *AliasManager) AddAlias(alias, command string) error {
	// Check if there are already too many aliases
	if len(am.aliases) >= am.maxCount {
		return fmt.Errorf("maximum number of aliases (%d) reached", am.maxCount)
	}

	// Check if the alias is valid
	if strings.Contains(alias, " ") || strings.Contains(alias, ".") {
		return fmt.Errorf("alias cannot contain spaces or periods")
	}

	// Check if the alias already exists
	if _, exists := am.aliases[alias]; exists {
		return fmt.Errorf("an alias with the name '%s' already exists", alias)
	}

	// Add alias
	am.aliases[alias] = command
	return nil
}

// RemoveAlias removes a local alias
func (am *AliasManager) RemoveAlias(alias string) error {
	// Check if the alias exists
	if _, exists := am.aliases[alias]; !exists {
		return fmt.Errorf("no alias with the name '%s' found", alias)
	}

	// Remove alias
	delete(am.aliases, alias)
	return nil
}

// GetAlias returns an alias if it exists
func (am *AliasManager) GetAlias(alias string) (string, bool) {
	command, exists := am.aliases[alias]
	return command, exists
}

// GetAllAliases returns all local aliases
func (am *AliasManager) GetAllAliases() map[string]string {
	// Create a copy to avoid modifying the internal map
	result := make(map[string]string, len(am.aliases))
	for alias, command := range am.aliases {
		result[alias] = command
	}
	return result
}

// SaveAliases saves all aliases to a file
func (am *AliasManager) SaveAliases() error {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	configDir := filepath.Join(userConfigDir, "nexuflex")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Create file
	aliasPath := filepath.Join(configDir, "local_aliases.txt")
	f, err := os.Create(aliasPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write aliases to file
	for alias, command := range am.aliases {
		if _, err := f.WriteString(fmt.Sprintf("%s=%s\n", alias, command)); err != nil {
			return err
		}
	}

	return nil
}

// LoadAliases loads aliases from a file
func (am *AliasManager) LoadAliases() error {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	aliasPath := filepath.Join(userConfigDir, "nexuflex", "local_aliases.txt")

	// Check if file exists
	if _, err := os.Stat(aliasPath); os.IsNotExist(err) {
		return nil // File doesn't exist, but that's not an error
	}

	// Open file
	f, err := os.Open(aliasPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Clear aliases
	am.aliases = make(map[string]string)

	// Read file line by line
	buffer := make([]byte, 4096)
	var line string
	for {
		n, err := f.Read(buffer)
		if err != nil {
			break // EOF or other error
		}

		// Process buffer
		for i := 0; i < n; i++ {
			if buffer[i] == '\n' {
				// Line end found, process alias
				if line != "" {
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 && len(parts[0]) > 0 {
						// Add alias, but only if the maximum count hasn't been reached
						if len(am.aliases) < am.maxCount {
							am.aliases[parts[0]] = parts[1]
						}
					}
				}
				line = ""
			} else if buffer[i] != '\r' { // Ignore CR
				line += string(buffer[i])
			}
		}
	}

	// Process last line if present
	if line != "" {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && len(parts[0]) > 0 {
			// Add alias, but only if the maximum count hasn't been reached
			if len(am.aliases) < am.maxCount {
				am.aliases[parts[0]] = parts[1]
			}
		}
	}

	return nil
}

// ExpandCommand replaces an alias with the full command
func (am *AliasManager) ExpandCommand(command string) string {
	// Trim command
	command = strings.TrimSpace(command)

	// Split command into parts
	parts := strings.SplitN(command, " ", 2)
	firstWord := parts[0]

	// Check if the first word is an alias
	if expandedCommand, ok := am.aliases[firstWord]; ok {
		// Add rest of command if present
		if len(parts) > 1 {
			return expandedCommand + " " + parts[1]
		}
		return expandedCommand
	}

	// No alias found, return original command
	return command
}

// IsReservedKeyword checks if a word is a reserved keyword
func IsReservedKeyword(word string) bool {
	// List of reserved keywords
	reservedKeywords := map[string]bool{
		"help":       true,
		"login":      true,
		"logout":     true,
		"alias":      true,
		"unalias":    true,
		"exit":       true,
		"quit":       true,
		"clear":      true,
		"history":    true,
		"use":        true,
		"connect":    true,
		"disconnect": true,
		"status":     true,
	}

	return reservedKeywords[strings.ToLower(word)]
}
