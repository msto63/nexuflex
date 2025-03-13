// input.go
/**
* Nexuflex Client - Input Field Implementation
*
* This file contains extensions for the input field of the user interface,
* including auto-completion and history navigation.
*
* @author msto63
* @version 1.0.0
* @date 2025-03-12
 */

package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/msto63/nexuflex/nexuflex-client/core"
	"github.com/rivo/tview"
)

// EnhancedInputField extends the standard InputField from tview
// with additional features like auto-completion and history navigation
type EnhancedInputField struct {
	*tview.InputField
	history          *core.CommandHistory
	aliasManager     *core.AliasManager
	autoCompleteFunc func(text string) ([]string, string)
	showCompletions  func([]string)
}

// NewEnhancedInputField creates an enhanced input field
func NewEnhancedInputField(
	history *core.CommandHistory,
	aliasManager *core.AliasManager,
	autoCompleteFunc func(text string) ([]string, string),
	showCompletions func([]string),
) *EnhancedInputField {
	input := &EnhancedInputField{
		InputField:       tview.NewInputField(),
		history:          history,
		aliasManager:     aliasManager,
		autoCompleteFunc: autoCompleteFunc,
		showCompletions:  showCompletions,
	}

	// Configure input field
	input.
		SetLabel("> ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack)

	// Enable custom keyboard handling
	input.SetInputCapture(input.handleKeyPress)

	return input
}

// handleKeyPress handles keyboard input
func (i *EnhancedInputField) handleKeyPress(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:
		// Previous command from history
		if cmd, ok := i.history.Previous(); ok {
			i.SetText(cmd)
		}
		return nil

	case tcell.KeyDown:
		// Next command from history
		if cmd, ok := i.history.Next(); ok {
			i.SetText(cmd)
		} else {
			// If there's no next command, clear input field
			i.SetText("")
		}
		return nil

	case tcell.KeyTab:
		// Auto-completion
		currentText := i.GetText()
		if i.autoCompleteFunc != nil {
			suggestions, commonPrefix := i.autoCompleteFunc(currentText)

			if len(suggestions) == 0 {
				// No suggestions
				return nil
			} else if len(suggestions) == 1 {
				// Exactly one suggestion - complete directly
				i.SetText(suggestions[0])
			} else if commonPrefix != "" && commonPrefix != currentText {
				// Common prefix found - complete to that
				i.SetText(commonPrefix)
			} else {
				// Multiple suggestions - show them
				if i.showCompletions != nil {
					i.showCompletions(suggestions)
				}
			}
		}
		return nil

	case tcell.KeyCtrlA:
		// Jump to start of line
		i.SetText(i.GetText())
		i.SetCursorPos(0)
		return nil

	case tcell.KeyCtrlE:
		// Jump to end of line
		text := i.GetText()
		i.SetCursorPos(len(text))
		return nil

	case tcell.KeyCtrlU:
		// Delete line up to cursor
		text := i.GetText()
		pos := i.GetCursorPos()
		if pos > 0 {
			i.SetText(text[pos:])
			i.SetCursorPos(0)
		}
		return nil

	case tcell.KeyCtrlK:
		// Delete line from cursor
		text := i.GetText()
		pos := i.GetCursorPos()
		if pos < len(text) {
			i.SetText(text[:pos])
		}
		return nil

	case tcell.KeyCtrlW:
		// Delete word
		text := i.GetText()
		pos := i.GetCursorPos()
		if pos > 0 {
			// Go backwards to space
			newPos := pos - 1
			for newPos > 0 && text[newPos] == ' ' {
				newPos--
			}
			for newPos > 0 && text[newPos-1] != ' ' {
				newPos--
			}

			if newPos < pos {
				newText := text[:newPos] + text[pos:]
				i.SetText(newText)
				i.SetCursorPos(newPos)
			}
		}
		return nil
	}

	// Default handling for other keys
	return event
}

// ProcessCommand processes the entered command
func (i *EnhancedInputField) ProcessCommand() string {
	command := i.GetText()

	// Ignore empty command
	if strings.TrimSpace(command) == "" {
		return ""
	}

	// Add command to history
	i.history.Add(command)
	i.history.ResetNavigation()

	// Clear input field
	i.SetText("")

	// Resolve local aliases
	if i.aliasManager != nil {
		command = i.aliasManager.ExpandCommand(command)
	}

	return command
}
