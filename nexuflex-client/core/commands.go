// commands.go
/**
* Nexuflex Client - Command History and Processing
*
* This file contains functions for managing command history
* and command preprocessing.
*
* @author msto63
* @version 1.0.0
* @date 2025-03-12
 */

package core

import (
	"os"
	"path/filepath"
	"strings"
)

// CommandHistory manages the command history
type CommandHistory struct {
	entries      []string
	maxEntries   int
	currentIndex int
	savePath     string
}

// NewCommandHistory creates a new command history
func NewCommandHistory(maxEntries int) *CommandHistory {
	return &CommandHistory{
		entries:      make([]string, 0, maxEntries),
		maxEntries:   maxEntries,
		currentIndex: -1,
	}
}

// Add adds a command to the history
func (h *CommandHistory) Add(command string) {
	// Don't add empty commands or commands that start with whitespace
	command = strings.TrimSpace(command)
	if command == "" {
		return
	}

	// Check if the command is already the last element in the history
	if len(h.entries) > 0 && h.entries[len(h.entries)-1] == command {
		return
	}

	// Add command to history
	h.entries = append(h.entries, command)

	// If history becomes too large, remove oldest entries
	if len(h.entries) > h.maxEntries {
		h.entries = h.entries[len(h.entries)-h.maxEntries:]
	}

	// Set index to end of history
	h.currentIndex = len(h.entries)
}

// Previous returns the previous command in the history
func (h *CommandHistory) Previous() (string, bool) {
	if len(h.entries) == 0 || h.currentIndex <= 0 {
		return "", false
	}

	h.currentIndex--
	return h.entries[h.currentIndex], true
}

// Next returns the next command in the history
func (h *CommandHistory) Next() (string, bool) {
	if len(h.entries) == 0 || h.currentIndex >= len(h.entries) {
		return "", false
	}

	h.currentIndex++
	if h.currentIndex == len(h.entries) {
		h.currentIndex = len(h.entries)
		return "", true // Empty string, but successful (for clearing the input line)
	}

	return h.entries[h.currentIndex], true
}

// ResetNavigation resets the navigation index
func (h *CommandHistory) ResetNavigation() {
	h.currentIndex = len(h.entries)
}

// GetEntries returns all entries in the history
func (h *CommandHistory) GetEntries() []string {
	return h.entries
}

// SetSavePath sets the path where the history is saved
func (h *CommandHistory) SetSavePath(path string) {
	h.savePath = path
}

// Save saves the history to a file
func (h *CommandHistory) Save() error {
	if h.savePath == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		// Ensure directory exists
		configDir := filepath.Join(userConfigDir, "nexuflex")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
		h.savePath = filepath.Join(configDir, "history.txt")
	}

	// Create directory for the file if it doesn't exist
	dir := filepath.Dir(h.savePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create history file and write
	f, err := os.Create(h.savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write commands line by line to the file
	for _, entry := range h.entries {
		if _, err := f.WriteString(entry + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// Load loads the history from a file
func (h *CommandHistory) Load() error {
	if h.savePath == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		h.savePath = filepath.Join(userConfigDir, "nexuflex", "history.txt")
	}

	// Check if file exists
	if _, err := os.Stat(h.savePath); os.IsNotExist(err) {
		return nil // File doesn't exist, but that's not an error
	}

	// Open file
	f, err := os.Open(h.savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Clear history
	h.entries = make([]string, 0, h.maxEntries)

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
				// End of line found, add command to history
				if line != "" {
					h.Add(line)
				}
				line = ""
			} else if buffer[i] != '\r' { // Ignore CR
				line += string(buffer[i])
			}
		}
	}

	// Add last line if present
	if line != "" {
		h.Add(line)
	}

	// Set index to end of history
	h.currentIndex = len(h.entries)

	return nil
}

// CommandProcessor processes commands before execution
type CommandProcessor struct {
	localAliases map[string]string
}

// NewCommandProcessor creates a new command processor
func NewCommandProcessor() *CommandProcessor {
	return &CommandProcessor{
		localAliases: make(map[string]string),
	}
}

// AddLocalAlias adds a local alias
func (p *CommandProcessor) AddLocalAlias(alias, command string) {
	p.localAliases[alias] = command
}

// RemoveLocalAlias removes a local alias
func (p *CommandProcessor) RemoveLocalAlias(alias string) {
	delete(p.localAliases, alias)
}

// GetLocalAliases returns all local aliases
func (p *CommandProcessor) GetLocalAliases() map[string]string {
	return p.localAliases
}

// ProcessCommand processes a command before execution
func (p *CommandProcessor) ProcessCommand(command string, useLocalAliases bool) string {
	// Trim command
	command = strings.TrimSpace(command)

	// Resolve local aliases
	if useLocalAliases {
		parts := strings.SplitN(command, " ", 2)
		firstWord := parts[0]

		// Check if the first word is an alias
		if expandedCommand, ok := p.localAliases[firstWord]; ok {
			// Add rest of command if present
			if len(parts) > 1 {
				command = expandedCommand + " " + parts[1]
			} else {
				command = expandedCommand
			}
		}
	}

	return command
}

// SaveLocalAliases saves the local aliases to a file
func (p *CommandProcessor) SaveLocalAliases() error {
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
	aliasPath := filepath.Join(configDir, "aliases.txt")
	f, err := os.Create(aliasPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write aliases to file
	for alias, command := range p.localAliases {
		if _, err := f.WriteString(alias + "=" + command + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// LoadLocalAliases loads the local aliases from a file
func (p *CommandProcessor) LoadLocalAliases() error {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	aliasPath := filepath.Join(userConfigDir, "nexuflex", "aliases.txt")

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
	p.localAliases = make(map[string]string)

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
				// End of line found, process alias
				if line != "" {
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						p.localAliases[parts[0]] = parts[1]
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
		if len(parts) == 2 {
			p.localAliases[parts[0]] = parts[1]
		}
	}

	return nil
}
