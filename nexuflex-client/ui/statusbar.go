// statusbar.go
/**
* Nexuflex Client - Status Bar Implementation
*
* This file contains the implementation of the status bar for the user interface.
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
	"github.com/msto63/nexuflex/nexuflex-client/i18n"
	"github.com/msto63/nexuflex/shared/proto"
	"github.com/rivo/tview"
)

// StatusBar is an extended status bar with two areas
type StatusBar struct {
	flex       *tview.Flex
	statusMsg  *tview.TextView
	statusInfo *tview.TextView
	app        *tview.Application
	msgTimer   *time.Timer
}

// NewStatusBar creates a new status bar
func NewStatusBar(app *tview.Application) *StatusBar {
	// Status message area on the left
	statusMsg := tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(tcell.ColorWhite)

	// Status information area on the right
	statusInfo := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetTextColor(tcell.ColorWhite)

	// Flex container
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(statusMsg, 0, 3, false).
		AddItem(statusInfo, 0, 1, false)

	flex.SetBackgroundColor(tcell.ColorDarkGray)

	return &StatusBar{
		flex:       flex,
		statusMsg:  statusMsg,
		statusInfo: statusInfo,
		app:        app,
	}
}

// GetPrimitive returns the tview.Primitive flex container
func (s *StatusBar) GetPrimitive() tview.Primitive {
	return s.flex
}

// ShowError displays a temporary error message in the status bar
func (s *StatusBar) ShowError(message string) {
	s.statusMsg.SetText(fmt.Sprintf("[red]%s[white]", message))
	s.app.Draw()

	// Clear message after 5 seconds
	s.startMessageTimer(5 * time.Second)
}

// ShowInfo displays a temporary information message in the status bar
func (s *StatusBar) ShowInfo(message string) {
	s.statusMsg.SetText(fmt.Sprintf("[green]%s[white]", message))
	s.app.Draw()

	// Clear message after 3 seconds
	s.startMessageTimer(3 * time.Second)
}

// ShowWarning displays a temporary warning message in the status bar
func (s *StatusBar) ShowWarning(message string) {
	s.statusMsg.SetText(fmt.Sprintf("[yellow]%s[white]", message))
	s.app.Draw()

	// Clear message after 4 seconds
	s.startMessageTimer(4 * time.Second)
}

// startMessageTimer starts a timer to clear the message
func (s *StatusBar) startMessageTimer(duration time.Duration) {
	// If a timer is already running, stop it
	if s.msgTimer != nil {
		s.msgTimer.Stop()
	}

	// Start new timer
	s.msgTimer = time.AfterFunc(duration, func() {
		s.app.QueueUpdateDraw(func() {
			s.statusMsg.SetText("")
		})
	})
}

// SetMessage sets a permanent message in the status bar
func (s *StatusBar) SetMessage(message string) {
	// If a timer is running, stop it
	if s.msgTimer != nil {
		s.msgTimer.Stop()
		s.msgTimer = nil
	}

	s.statusMsg.SetText(message)
	s.app.Draw()
}

// UpdateStatus updates the status display with information from the Proto-StatusInfo
func (s *StatusBar) UpdateStatus(statusInfo *proto.StatusInfo) {
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
	s.statusInfo.SetText(statusText.String())
	s.app.Draw()
}

// Clear clears both text areas of the status bar
func (s *StatusBar) Clear() {
	s.statusMsg.SetText("")
	s.statusInfo.SetText("")
	s.app.Draw()
}

// SetBackgroundColor changes the background color of the status bar
func (s *StatusBar) SetBackgroundColor(color tcell.Color) {
	s.flex.SetBackgroundColor(color)
	s.app.Draw()
}
