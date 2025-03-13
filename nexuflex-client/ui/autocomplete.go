// autocomplete.go
/**
* Nexuflex Client - Auto-completion Implementation
*
* This file contains the implementation of command completion
* for the user interface.
*
* @author msto63
* @version 1.0.0
* @date 2025-03-12
 */

package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

// AutoCompleter provides functions for command completion
type AutoCompleter struct {
	output            *tview.TextView
	localCommands     map[string]bool
	fallbackHandler   func(text string) ([]string, string, error)
	cachedSuggestions map[string][]string
}

// NewAutoCompleter creates a new AutoCompleter
func NewAutoCompleter(output *tview.TextView, fallbackHandler func(text string) ([]string, string, error)) *AutoCompleter {
	// Register standard commands
	localCommands := map[string]bool{
		"help":       true,
		"?":          true,
		"exit":       true,
		"quit":       true,
		"clear":      true,
		"cls":        true,
		"connect":    true,
		"disconnect": true,
		"login":      true,
		"logout":     true,
		"alias":      true,
		"unalias":    true,
		"history":    true,
		"use":        true,
	}

	return &AutoCompleter{
		output:            output,
		localCommands:     localCommands,
		fallbackHandler:   fallbackHandler,
		cachedSuggestions: make(map[string][]string),
	}
}

// Complete attempts to complete the entered text
func (ac *AutoCompleter) Complete(text string) ([]string, string) {
	// Trim whitespace at beginning and end
	text = strings.TrimSpace(text)

	// Handle empty text
	if text == "" {
		// Return local commands and server context
		suggestions := make([]string, 0, len(ac.localCommands))
		for cmd := range ac.localCommands {
			suggestions = append(suggestions, cmd)
		}

		// Ask server for completions
		if ac.fallbackHandler != nil {
			serverSuggestions, _, _ := ac.fallbackHandler(text)
			suggestions = append(suggestions, serverSuggestions...)
		}

		return suggestions, ""
	}

	// Complete local commands
	if !strings.Contains(text, ".") {
		// Look for possible local commands
		localSuggestions := make([]string, 0)
		for cmd := range ac.localCommands {
			if strings.HasPrefix(cmd, text) {
				localSuggestions = append(localSuggestions, cmd)
			}
		}

		// If local completions found
		if len(localSuggestions) > 0 {
			return localSuggestions, findCommonPrefix(localSuggestions)
		}
	}

	// Try server-side completion
	if ac.fallbackHandler != nil {
		// First check cache
		if suggestions, ok := ac.cachedSuggestions[text]; ok {
			return suggestions, findCommonPrefix(suggestions)
		}

		// Ask server
		suggestions, commonPrefix, err := ac.fallbackHandler(text)
		if err == nil && len(suggestions) > 0 {
			// Store in cache
			ac.cachedSuggestions[text] = suggestions
			return suggestions, commonPrefix
		}
	}

	return []string{}, ""
}

// ShowSuggestions displays completion suggestions in the output area
func (ac *AutoCompleter) ShowSuggestions(suggestions []string) {
	if len(suggestions) == 0 {
		return
	}

	// Prepare text output
	var sb strings.Builder
	sb.WriteString("[blue]Possible completions:[white]\n")

	// Group suggestions
	groups := groupSuggestions(suggestions)

	// Output groups
	for group, items := range groups {
		if group != "" {
			sb.WriteString(fmt.Sprintf("[yellow]%s:[white]\n", group))
		}

		// Format entries in columns
		columns := formatInColumns(items, 4, 20)
		sb.WriteString(columns)
		sb.WriteString("\n")
	}

	// Output
	if ac.output != nil {
		ac.output.Write([]byte(sb.String()))
	}
}

// InvalidateCache clears the suggestions cache
func (ac *AutoCompleter) InvalidateCache() {
	ac.cachedSuggestions = make(map[string][]string)
}

// AddLocalCommand adds a local command to completion
func (ac *AutoCompleter) AddLocalCommand(command string) {
	ac.localCommands[command] = true
}

// RemoveLocalCommand removes a local command from completion
func (ac *AutoCompleter) RemoveLocalCommand(command string) {
	delete(ac.localCommands, command)
}

// Helper functions

// findCommonPrefix finds the common prefix in a list of strings
func findCommonPrefix(strings []string) string {
	if len(strings) == 0 {
		return ""
	}
	if len(strings) == 1 {
		return strings[0]
	}

	prefix := strings[0]
	for i := 1; i < len(strings); i++ {
		j := 0
		for j < len(prefix) && j < len(strings[i]) && prefix[j] == strings[i][j] {
			j++
		}
		prefix = prefix[:j]
	}

	return prefix
}

// groupSuggestions groups suggestions by service
func groupSuggestions(suggestions []string) map[string][]string {
	groups := make(map[string][]string)

	for _, suggestion := range suggestions {
		var group, item string

		if strings.Contains(suggestion, ".") {
			parts := strings.SplitN(suggestion, ".", 2)
			group = parts[0]
			item = suggestion
		} else {
			group = ""
			item = suggestion
		}

		if _, ok := groups[group]; !ok {
			groups[group] = make([]string, 0)
		}

		groups[group] = append(groups[group], item)
	}

	return groups
}

// formatInColumns formats a list of strings in columns
func formatInColumns(items []string, numColumns, columnWidth int) string {
	if len(items) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, item := range items {
		// Start new line
		if i%numColumns == 0 && i > 0 {
			sb.WriteString("\n")
		}

		// Format item and add to line
		format := fmt.Sprintf("%%-%ds", columnWidth)
		formattedItem := fmt.Sprintf(format, item)

		// Limit to maximum length
		if len(formattedItem) > columnWidth {
			formattedItem = formattedItem[:columnWidth-3] + "..."
		}

		sb.WriteString(formattedItem)
	}

	return sb.String()
}
