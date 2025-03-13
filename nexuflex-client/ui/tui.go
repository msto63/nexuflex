// tui.go
/**
 * Nexuflex Client - Text User Interface Main Class
 *
 * This file contains the main class for the text-based user interface (TUI)
 * of the nexuflex client.
 *
 * @author msto63
 * @version 1.0.0
 * @date 2025-03-12
 */

package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/msto63/nexuflex/nexuflex-client/core"
	"github.com/msto63/nexuflex/nexuflex-client/i18n"
	"github.com/msto63/nexuflex/shared/proto"
	"github.com/rivo/tview"
)

// TUI represents the complete text-based user interface
type TUI struct {
	// Main components
	app        *tview.Application
	pages      *tview.Pages
	layout     *tview.Flex
	header     *tview.TextView
	output     *tview.TextView
	input      *tview.InputField
	statusBar  *tview.Flex
	statusText *tview.TextView
	statusInfo *tview.TextView

	// Dialogs
	loginForm  *tview.Form
	serverList *tview.List
	helpText   *tview.TextView

	// Client and other components
	client         *core.Client
	commandHistory *core.CommandHistory
	aliasManager   *core.AliasManager

	// Status
	lastCommand   string
	statusMessage string
}

// NewTUI creates a new TUI instance
func NewTUI(client *core.Client) *TUI {
	// Create new TUI instance
	tui := &TUI{
		app:            tview.NewApplication(),
		pages:          tview.NewPages(),
		client:         client,
		commandHistory: core.NewCommandHistory(100), // 100 entries in history
		aliasManager:   core.NewAliasManager(50),    // 50 aliases maximum
	}

	// Initialize user interface
	tui.initUI()

	// Set callbacks for the client
	client.SetCallbacks(
		tui.handleStatusChanged,
		tui.handleServerList,
		tui.handleOutput,
	)

	// Load command history and aliases
	tui.commandHistory.Load()
	tui.aliasManager.LoadAliases()

	return tui
}

// initUI initializes the user interface
func (t *TUI) initUI() {
	// Create header
	t.header = tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(i18n.GetMessage("ui.header")).
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorBlue)

	// Create output area
	t.output = tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			t.app.Draw()
		})
	t.output.SetBorder(true).SetTitle(i18n.GetMessage("ui.output_title"))

	// Create input field
	t.input = tview.NewInputField().
		SetLabel(i18n.GetMessage("ui.command_prompt")).
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetDoneFunc(t.handleCommand)

	// Create status bar
	t.statusText = tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(tcell.ColorWhite)
	t.statusInfo = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetTextColor(tcell.ColorWhite)

	t.statusBar = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(t.statusText, 0, 3, false).
		AddItem(t.statusInfo, 0, 1, false)
	t.statusBar.SetBackgroundColor(tcell.ColorDarkGray)

	// Create layout
	t.layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.header, 1, 0, false).
		AddItem(t.output, 0, 1, false).
		AddItem(t.input, 1, 0, true).
		AddItem(t.statusBar, 1, 0, false)

	// Create login form
	t.loginForm = tview.NewForm().
		AddInputField(i18n.GetMessage("ui.username"), "", 20, nil, nil).
		AddPasswordField(i18n.GetMessage("ui.password"), "", 20, '*', nil).
		AddButton(i18n.GetMessage("ui.login_button"), t.handleLogin).
		AddButton(i18n.GetMessage("ui.cancel_button"), func() {
			t.pages.SwitchToPage("main")
		})
	t.loginForm.SetBorder(true).SetTitle(i18n.GetMessage("ui.login_title")).SetTitleAlign(tview.AlignCenter)
	t.loginForm.SetBackgroundColor(tcell.ColorBlack)

	// Create server list
	t.serverList = tview.NewList().
		ShowSecondaryText(true).
		SetSecondaryTextColor(tcell.ColorDimGray)
	t.serverList.SetBorder(true).SetTitle(i18n.GetMessage("ui.available_servers")).SetTitleAlign(tview.AlignCenter)
	t.serverList.SetDoneFunc(func() {
		t.pages.SwitchToPage("main")
	})

	// Create help text
	t.helpText = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetText(getHelpText())
	t.helpText.SetBorder(true).SetTitle(i18n.GetMessage("ui.help_title")).SetTitleAlign(tview.AlignCenter)
	t.helpText.SetDoneFunc(func(key tcell.Key) {
		t.pages.SwitchToPage("main")
	})

	// Add pages
	t.pages.AddPage("main", t.layout, true, true)
	t.pages.AddPage("login", centeredFlex(t.loginForm, 40, 10), true, false)
	t.pages.AddPage("servers", centeredFlex(t.serverList, 60, 20), true, false)
	t.pages.AddPage("help", centeredFlex(t.helpText, 70, 20), true, false)

	// Keyboard shortcuts
	t.app.SetInputCapture(t.handleGlobalKeys)
	t.input.SetInputCapture(t.handleInputKeys)
}

// Run starts the user interface
func (t *TUI) Run() error {
	// Set status
	t.updateStatus(i18n.GetMessage("general.ready"), &proto.StatusInfo{
		ConnectionStatus: proto.StatusInfo_OFFLINE,
		SessionStatus:    proto.StatusInfo_NOT_LOGGED_IN,
	})

	// Display initial text
	t.output.SetText(i18n.GetMessage("general.welcome_message"))

	// Start the application
	return t.app.SetRoot(t.pages, true).EnableMouse(true).Run()
}

// ShowError displays an error message in the status bar
func (t *TUI) ShowError(message string) {
	t.statusText.SetText(fmt.Sprintf("[red]%s[white]", message))
	t.app.Draw()

	// Clear message after 5 seconds
	go func() {
		time.Sleep(5 * time.Second)
		t.app.QueueUpdateDraw(func() {
			// Only clear if the text is still the same
			if strings.Contains(t.statusText.GetText(true), message) {
				t.statusText.SetText("")
			}
		})
	}()
}

// ShowInfo displays an information message in the status bar
func (t *TUI) ShowInfo(message string) {
	t.statusText.SetText(fmt.Sprintf("[green]%s[white]", message))
	t.app.Draw()

	// Clear message after 3 seconds
	go func() {
		time.Sleep(3 * time.Second)
		t.app.QueueUpdateDraw(func() {
			// Only clear if the text is still the same
			if strings.Contains(t.statusText.GetText(true), message) {
				t.statusText.SetText("")
			}
		})
	}()
}

// handleCommand processes the entered command line
func (t *TUI) handleCommand(key tcell.Key) {
	// Get command
	command := t.input.GetText()

	// Ignore empty command
	if strings.TrimSpace(command) == "" {
		return
	}

	// Resolve aliases
	command = t.aliasManager.ExpandCommand(command)

	// Add command to history
	t.commandHistory.Add(command)

	// Clear input field
	t.input.SetText("")

	// Display output in terminal
	t.output.Write([]byte(fmt.Sprintf("> [yellow]%s[white]\n", command)))

	// Process special client commands
	if t.handleSpecialCommand(command) {
		return
	}

	// Send command to server
	if t.client.IsConnected() {
		err := t.client.ExecuteCommand(command)
		if err != nil {
			t.ShowError(err.Error())
		}
	} else {
		t.ShowError(i18n.GetMessage("error.not_connected"))
	}
}

// handleSpecialCommand processes special client-side commands
func (t *TUI) handleSpecialCommand(command string) bool {
	command = strings.TrimSpace(command)
	parts := strings.SplitN(command, " ", 2)
	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "help", "?":
		// Show help
		t.pages.SwitchToPage("help")
		return true

	case "exit", "quit":
		// Exit application
		t.app.Stop()
		return true

	case "clear", "cls":
		// Clear output
		t.output.SetText("")
		return true

	case "connect":
		// Connect to server
		if len(parts) < 2 {
			t.ShowError(fmt.Sprintf(i18n.GetMessage("commands.syntax"), "connect <host> [port]"))
			return true
		}

		connectParts := strings.Split(parts[1], " ")
		host := connectParts[0]
		port := 50051 // Default port

		if len(connectParts) > 1 {
			if _, err := fmt.Sscanf(connectParts[1], "%d", &port); err != nil {
				t.ShowError(fmt.Sprintf("Invalid port: %s", connectParts[1]))
				return true
			}
		}

		err := t.client.Connect(host, port, false)
		if err != nil {
			t.ShowError(err.Error())
		} else {
			t.ShowInfo(fmt.Sprintf(i18n.GetMessage("success.connected"), host, port))
		}
		return true

	case "disconnect":
		// Disconnect from server
		t.client.Close()
		t.updateStatus(i18n.GetMessage("success.disconnected"), &proto.StatusInfo{
			ConnectionStatus: proto.StatusInfo_OFFLINE,
			SessionStatus:    proto.StatusInfo_NOT_LOGGED_IN,
		})
		return true

	case "login":
		// Show login dialog
		t.pages.SwitchToPage("login")
		return true

	case "logout":
		// Log out
		if !t.client.IsConnected() {
			t.ShowError(i18n.GetMessage("error.not_connected"))
			return true
		}

		if !t.client.IsLoggedIn() {
			t.ShowError(i18n.GetMessage("error.not_logged_in"))
			return true
		}

		err := t.client.Logout()
		if err != nil {
			t.ShowError(err.Error())
		} else {
			t.ShowInfo(i18n.GetMessage("success.logged_out"))
		}
		return true

	case "alias":
		// Define or show aliases
		if len(parts) < 2 {
			// Show aliases
			aliases := t.aliasManager.GetAllAliases()
			if len(aliases) == 0 {
				t.output.Write([]byte(i18n.GetMessage("commands.no_aliases") + "\n"))
			} else {
				t.output.Write([]byte(i18n.GetMessage("commands.local_aliases") + "\n"))
				for alias, command := range aliases {
					t.output.Write([]byte(fmt.Sprintf("  %s = %s\n", alias, command)))
				}
			}
		} else {
			// Define alias
			aliasParts := strings.SplitN(parts[1], "=", 2)
			if len(aliasParts) != 2 {
				t.ShowError(fmt.Sprintf(i18n.GetMessage("commands.syntax"), "alias <name>=<command>"))
				return true
			}

			alias := strings.TrimSpace(aliasParts[0])
			command := strings.TrimSpace(aliasParts[1])

			if alias == "" {
				t.ShowError(i18n.GetMessage("error.empty_alias"))
				return true
			}

			if command == "" {
				t.ShowError(i18n.GetMessage("error.empty_command"))
				return true
			}

			if isReservedKeyword(alias) {
				t.ShowError(fmt.Sprintf(i18n.GetMessage("error.reserved_keyword"), alias))
				return true
			}

			err := t.aliasManager.AddAlias(alias, command)
			if err != nil {
				t.ShowError(err.Error())
			} else {
				t.ShowInfo(fmt.Sprintf(i18n.GetMessage("success.alias_created"), alias, command))
				t.aliasManager.SaveAliases()
			}
		}
		return true

	case "unalias":
		// Delete alias
		if len(parts) < 2 {
			t.ShowError(fmt.Sprintf(i18n.GetMessage("commands.syntax"), "unalias <name>"))
			return true
		}

		alias := strings.TrimSpace(parts[1])
		err := t.aliasManager.RemoveAlias(alias)
		if err != nil {
			t.ShowError(err.Error())
		} else {
			t.ShowInfo(fmt.Sprintf(i18n.GetMessage("success.alias_deleted"), alias))
			t.aliasManager.SaveAliases()
		}
		return true

	case "history":
		// Show history
		entries := t.commandHistory.GetEntries()
		if len(entries) == 0 {
			t.output.Write([]byte(i18n.GetMessage("commands.no_history") + "\n"))
		} else {
			t.output.Write([]byte(i18n.GetMessage("commands.command_history") + "\n"))
			for i, entry := range entries {
				t.output.Write([]byte(fmt.Sprintf("  %d: %s\n", i+1, entry)))
			}
		}
		return true

	case "use":
		// Set service context
		if len(parts) < 2 {
			t.output.Write([]byte(fmt.Sprintf(i18n.GetMessage("commands.current_context"),
				t.client.GetLastServiceUsed())))
			return true
		}

		service := strings.TrimSpace(parts[1])
		t.client.SetLastServiceUsed(service)
		t.ShowInfo(fmt.Sprintf(i18n.GetMessage("commands.context_set"), service))
		return true
	}

	return false
}

// handleLogin processes the login
func (t *TUI) handleLogin() {
	username := t.loginForm.GetFormItem(0).(*tview.InputField).GetText()
	password := t.loginForm.GetFormItem(1).(*tview.InputField).GetText()

	// Reset form
	t.loginForm.GetFormItem(1).(*tview.InputField).SetText("")

	// Return to main page
	t.pages.SwitchToPage("main")

	// Check if connected to server
	if !t.client.IsConnected() {
		t.ShowError(i18n.GetMessage("error.not_connected"))
		return
	}

	// Login
	err := t.client.Login(username, password)
	if err != nil {
		t.ShowError(err.Error())
	}
}

// handleServerList processes the server list
func (t *TUI) handleServerList(servers []*proto.ServerInfo) (int, error) {
	// Clear list
	t.serverList.Clear()

	// Add servers to list
	for i, server := range servers {
		title := fmt.Sprintf("%s (%s)", server.ShortName, server.Address)
		secondary := fmt.Sprintf("Version: %s, TLS: %v", server.Version, server.TlsEnabled)

		t.serverList.AddItem(title, secondary, rune('1'+i), func(index int) func() {
			return func() {
				t.pages.SwitchToPage("main")
				// Return selected index
				// (processed later)
			}
		}(i))
	}

	// Show list
	t.pages.SwitchToPage("servers")

	// Wait for selection
	selectedIndex := -1

	// Since we need a return value, we have to wait here
	// In a real implementation, we would probably use a channel
	// or perform discovery asynchronously in the background

	return selectedIndex, nil
}

// handleOutput processes output from the server
func (t *TUI) handleOutput(output string) {
	t.output.Write([]byte(output + "\n"))
}

// handleStatusChanged processes status changes
func (t *TUI) handleStatusChanged(statusInfo *proto.StatusInfo) {
	t.updateStatus("", statusInfo)
}

// updateStatus updates the status display
func (t *TUI) updateStatus(message string, statusInfo *proto.StatusInfo) {
	if message != "" {
		t.statusText.SetText(message)
	}

	if statusInfo == nil {
		return
	}

	// Create status text
	var statusText strings.Builder

	// Connection status
	switch statusInfo.ConnectionStatus {
	case proto.StatusInfo_OFFLINE:
		statusText.WriteString("[red]" + i18n.GetMessage("status.offline") + "[white]")
	case proto.StatusInfo_CONNECTING:
		statusText.WriteString("[yellow]" + i18n.GetMessage("status.connecting") + "[white]")
	case proto.StatusInfo_CONNECTED:
		if statusInfo.ServerName != "" {
			statusText.WriteString(fmt.Sprintf("[green]%s[white]",
				fmt.Sprintf(i18n.GetMessage("status.connected"), statusInfo.ServerName)))
		} else {
			statusText.WriteString("[green]" + i18n.GetMessage("status.connected") + "[white]")
		}
	case proto.StatusInfo_CONNECTION_ERROR:
		statusText.WriteString("[red]" + i18n.GetMessage("status.connection_error") + "[white]")
	}

	// Separator
	statusText.WriteString(" | ")

	// Session status
	switch statusInfo.SessionStatus {
	case proto.StatusInfo_NOT_LOGGED_IN:
		statusText.WriteString("[yellow]" + i18n.GetMessage("status.not_logged_in") + "[white]")
	case proto.StatusInfo_AUTHENTICATED:
		if statusInfo.Username != "" {
			statusText.WriteString(fmt.Sprintf("[green]%s[white]",
				fmt.Sprintf(i18n.GetMessage("status.logged_in"), statusInfo.Username)))
		} else {
			statusText.WriteString("[green]" + i18n.GetMessage("status.logged_in") + "[white]")
		}
	case proto.StatusInfo_LOGIN_REQUIRED:
		statusText.WriteString("[yellow]" + i18n.GetMessage("status.login_required") + "[white]")
	case proto.StatusInfo_SESSION_EXPIRING:
		remaining := statusInfo.SessionRemainingMinutes
		statusText.WriteString(fmt.Sprintf("[yellow]%s[white]",
			fmt.Sprintf(i18n.GetMessage("status.session_expiring"), remaining)))
	case proto.StatusInfo_SESSION_EXPIRED:
		statusText.WriteString("[red]" + i18n.GetMessage("status.session_expired") + "[white]")
	}

	// Service context
	if statusInfo.CurrentService != "" {
		statusText.WriteString(fmt.Sprintf(" | %s",
			fmt.Sprintf(i18n.GetMessage("status.service_context"), statusInfo.CurrentService)))
	}

	// Update status display
	t.statusInfo.SetText(statusText.String())
	t.app.Draw()
}

// handleGlobalKeys processes global keyboard shortcuts
func (t *TUI) handleGlobalKeys(event *tcell.EventKey) *tcell.EventKey {
	// If a modal dialog is active, only process Escape
	if t.pages.HasPage("modal") {
		if event.Key() == tcell.KeyEscape {
			t.pages.RemovePage("modal")
			return nil
		}
		return event
	}

	// Global keyboard shortcuts
	switch event.Key() {
	case tcell.KeyCtrlC:
		// Exit application
		t.app.Stop()
		return nil

	case tcell.KeyCtrlL:
		// Show login dialog
		if t.pages.HasPage("login") {
			t.pages.SwitchToPage("login")
			return nil
		}

	case tcell.KeyCtrlH:
		// Show help
		if t.pages.HasPage("help") {
			t.pages.SwitchToPage("help")
			return nil
		}

	case tcell.KeyCtrlD:
		// Start server discovery
		go func() {
			err := t.client.DiscoverServer(5 * time.Second)
			if err != nil {
				t.app.QueueUpdateDraw(func() {
					t.ShowError(fmt.Sprintf(i18n.GetMessage("error.discovery"), err))
				})
			}
		}()
		return nil
	}

	return event
}

// handleInputKeys processes keyboard shortcuts in the input field
func (t *TUI) handleInputKeys(event *tcell.EventKey) *tcell.EventKey {
	// History navigation
	switch event.Key() {
	case tcell.KeyUp:
		// Previous command
		if cmd, ok := t.commandHistory.Previous(); ok {
			t.input.SetText(cmd)
		}
		return nil

	case tcell.KeyDown:
		// Next command
		if cmd, ok := t.commandHistory.Next(); ok {
			t.input.SetText(cmd)
		}
		return nil

	case tcell.KeyTab:
		// Auto-completion
		currentText := t.input.GetText()
		if t.client.IsConnected() {
			suggestions, commonPrefix, err := t.client.AutoComplete(currentText, len(currentText))
			if err == nil && len(suggestions) > 0 {
				if len(suggestions) == 1 {
					// Only one suggestion - complete directly
					t.input.SetText(suggestions[0])
				} else if commonPrefix != "" && commonPrefix != currentText {
					// Complete common prefix
					t.input.SetText(commonPrefix)
				} else {
					// Multiple suggestions - show them
					t.output.Write([]byte("Possible completions:\n"))
					for _, suggestion := range suggestions {
						t.output.Write([]byte(fmt.Sprintf("  %s\n", suggestion)))
					}
				}
			}
		}
		return nil
	}

	return event
}

// centeredFlex centers a flex element on the screen
func centeredFlex(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false),
			width, 1, true).
		AddItem(nil, 0, 1, false)
}

// getHelpText returns the help text for the application
func getHelpText() string {
	return fmt.Sprintf(`[yellow]%s[white]
 
 [blue]%s:[white]
   [yellow]help[white] or [yellow]?[white]          %s
   [yellow]exit[white] or [yellow]quit[white]       %s
   [yellow]clear[white] or [yellow]cls[white]       %s
   [yellow]history[white]               %s
 
 [blue]%s:[white]
   [yellow]connect <host> [port][white]  %s
   [yellow]disconnect[white]             %s
 
 [blue]%s:[white]
   [yellow]login[white]                  %s
   [yellow]logout[white]                 %s
 
 [blue]%s:[white]
   [yellow]alias[white]                  %s
   [yellow]alias <n>=<command>[white]    %s
   [yellow]unalias <n>[white]            %s
 
 [blue]%s:[white]
   [yellow]use <service>[white]          %s
 
 [blue]%s:[white]
   [yellow]Ctrl+H[white]                 %s
   [yellow]Ctrl+L[white]                 %s
   [yellow]Ctrl+D[white]                 %s
   [yellow]Ctrl+C[white]                 %s
   [yellow]↑/↓[white]                    %s
   [yellow]Tab[white]                    %s
 
 [blue]%s:[white]
   [yellow]<Service>.<Action>.<SubAction> <Parameters>[white]
 
   %s: [yellow]Finance.Create.Report Q4_2024 "Profit and Loss"[white]
 
 %s`,
		i18n.GetMessage("help.title"),
		i18n.GetMessage("help.general_commands"),
		i18n.GetMessage("help.help_command"),
		i18n.GetMessage("help.exit_command"),
		i18n.GetMessage("help.clear_command"),
		i18n.GetMessage("help.history_command"),
		i18n.GetMessage("help.connection_management"),
		i18n.GetMessage("help.connect_command"),
		i18n.GetMessage("help.disconnect_command"),
		i18n.GetMessage("help.authentication"),
		i18n.GetMessage("help.login_command"),
		i18n.GetMessage("help.logout_command"),
		i18n.GetMessage("help.alias_management"),
		i18n.GetMessage("help.alias_list_command"),
		i18n.GetMessage("help.alias_create_command"),
		i18n.GetMessage("help.alias_delete_command"),
		i18n.GetMessage("help.context"),
		i18n.GetMessage("help.context_command"),
		i18n.GetMessage("help.keyboard_shortcuts"),
		i18n.GetMessage("help.ctrl_h"),
		i18n.GetMessage("help.ctrl_l"),
		i18n.GetMessage("help.ctrl_d"),
		i18n.GetMessage("help.ctrl_c"),
		i18n.GetMessage("help.arrow_keys"),
		i18n.GetMessage("help.tab_key"),
		i18n.GetMessage("help.command_format"),
		"Example",
		"Press any key to return to the main application.")
}

// isReservedKeyword checks if a word is a reserved keyword
func isReservedKeyword(word string) bool {
	// List of reserved keywords
	reservedKeywords := map[string]bool{
		"help":       true,
		"?":          true,
		"login":      true,
		"logout":     true,
		"alias":      true,
		"unalias":    true,
		"exit":       true,
		"quit":       true,
		"clear":      true,
		"cls":        true,
		"history":    true,
		"use":        true,
		"connect":    true,
		"disconnect": true,
		"status":     true,
	}

	return reservedKeywords[strings.ToLower(word)]
}
