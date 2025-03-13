// keybindings.go
/**
* Nexuflex Client - Key Bindings
*
* This file contains the definitions and processing of key bindings
* for the user interface.
*
* @author msto63
* @version 1.0.0
* @date 2025-03-12
 */

package ui

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

// KeyHandler is a type for keyboard handling functions
type KeyHandler func() bool

// KeyBindings manages the key bindings of the application
type KeyBindings struct {
	globalHandlers map[tcell.Key]KeyHandler
	inputHandlers  map[tcell.Key]KeyHandler
	outputHandlers map[tcell.Key]KeyHandler
	helpText       map[tcell.Key]string
}

// NewKeyBindings creates a new instance of key binding management
func NewKeyBindings() *KeyBindings {
	return &KeyBindings{
		globalHandlers: make(map[tcell.Key]KeyHandler),
		inputHandlers:  make(map[tcell.Key]KeyHandler),
		outputHandlers: make(map[tcell.Key]KeyHandler),
		helpText:       make(map[tcell.Key]string),
	}
}

// AddGlobalHandler adds a global keyboard handler
func (kb *KeyBindings) AddGlobalHandler(key tcell.Key, handler KeyHandler, helpText string) {
	kb.globalHandlers[key] = handler
	if helpText != "" {
		kb.helpText[key] = helpText
	}
}

// AddInputHandler adds a keyboard handler for the input field
func (kb *KeyBindings) AddInputHandler(key tcell.Key, handler KeyHandler, helpText string) {
	kb.inputHandlers[key] = handler
	if helpText != "" {
		kb.helpText[key] = helpText
	}
}

// AddOutputHandler adds a keyboard handler for the output field
func (kb *KeyBindings) AddOutputHandler(key tcell.Key, handler KeyHandler, helpText string) {
	kb.outputHandlers[key] = handler
	if helpText != "" {
		kb.helpText[key] = helpText
	}
}

// HandleGlobalKey processes a keypress in the global context
func (kb *KeyBindings) HandleGlobalKey(event *tcell.EventKey) *tcell.EventKey {
	if handler, ok := kb.globalHandlers[event.Key()]; ok {
		if handler() {
			return nil // Key was processed
		}
	}

	return event // Pass key on
}

// HandleInputKey processes a keypress in the input field
func (kb *KeyBindings) HandleInputKey(event *tcell.EventKey) *tcell.EventKey {
	if handler, ok := kb.inputHandlers[event.Key()]; ok {
		if handler() {
			return nil // Key was processed
		}
	}

	return event // Pass key on
}

// HandleOutputKey processes a keypress in the output field
func (kb *KeyBindings) HandleOutputKey(event *tcell.EventKey) *tcell.EventKey {
	if handler, ok := kb.outputHandlers[event.Key()]; ok {
		if handler() {
			return nil // Key was processed
		}
	}

	return event // Pass key on
}

// GetHelpText returns the help text for a key
func (kb *KeyBindings) GetHelpText(key tcell.Key) string {
	if text, ok := kb.helpText[key]; ok {
		return text
	}
	return ""
}

// GetAllHelpTexts returns all help texts
func (kb *KeyBindings) GetAllHelpTexts() map[tcell.Key]string {
	return kb.helpText
}

// SetupDefaultKeyBindings configures the default key bindings for the application
func SetupDefaultKeyBindings(tui *TUI) *KeyBindings {
	kb := NewKeyBindings()

	// Global key bindings
	kb.AddGlobalHandler(tcell.KeyCtrlC, func() bool {
		tui.app.Stop()
		return true
	}, "Exits the application")

	kb.AddGlobalHandler(tcell.KeyCtrlL, func() bool {
		tui.pages.SwitchToPage("login")
		return true
	}, "Opens the login dialog")

	kb.AddGlobalHandler(tcell.KeyCtrlH, func() bool {
		tui.pages.SwitchToPage("help")
		return true
	}, "Shows the help")

	kb.AddGlobalHandler(tcell.KeyCtrlD, func() bool {
		go func() {
			err := tui.client.DiscoverServer(5 * time.Second)
			if err != nil {
				tui.app.QueueUpdateDraw(func() {
					tui.ShowError(err.Error())
				})
			}
		}()
		return true
	}, "Starts server discovery")

	kb.AddGlobalHandler(tcell.KeyEscape, func() bool {
		// If a modal dialog is active, close it
		if tui.pages.HasPage("modal") {
			tui.pages.RemovePage("modal")
			return true
		}
		// Otherwise, if not on main page, return
		if tui.pages.GetCurrentPage() != "main" {
			tui.pages.SwitchToPage("main")
			return true
		}
		return false
	}, "Closes dialogs or returns to main view")

	// Input field key bindings
	kb.AddInputHandler(tcell.KeyUp, func() bool {
		// Get previous command from history
		return true
	}, "Previous command from history")

	kb.AddInputHandler(tcell.KeyDown, func() bool {
		// Get next command from history
		return true
	}, "Next command from history")

	kb.AddInputHandler(tcell.KeyTab, func() bool {
		// Auto-completion
		return true
	}, "Command completion")

	// Output field key bindings
	kb.AddOutputHandler(tcell.KeyPgUp, func() bool {
		// Scroll page up
		return true
	}, "Scroll page up")

	kb.AddOutputHandler(tcell.KeyPgDn, func() bool {
		// Scroll page down
		return true
	}, "Scroll page down")

	kb.AddOutputHandler(tcell.KeyHome, func() bool {
		// Scroll to start
		return true
	}, "Scroll to start of output")

	kb.AddOutputHandler(tcell.KeyEnd, func() bool {
		// Scroll to end
		return true
	}, "Scroll to end of output")

	return kb
}
