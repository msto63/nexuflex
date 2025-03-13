// layout.go
/**
* Nexuflex Client - UI Layout
*
* This file contains helper functions for the user interface layout.
*
* @author msto63
* @version 1.0.0
* @date 2025-03-12
 */

package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateHeader creates a header for the TUI
func CreateHeader(title string, backgroundColor, textColor tcell.Color) *tview.TextView {
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(title).
		SetTextColor(textColor).
		SetBackgroundColor(backgroundColor)

	return header
}

// CreateOutput creates the output area for the TUI
func CreateOutput(title string) *tview.TextView {
	output := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			// Will be set later
		})

	output.SetBorder(true).
		SetTitle(title).
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 1)

	return output
}

// CreateInput creates the input field for the TUI
func CreateInput(label string, doneFunc func(key tcell.Key) bool) *tview.InputField {
	input := tview.NewInputField().
		SetLabel(label).
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetDoneFunc(doneFunc)

	return input
}

// CreateStatusBar creates the status bar for the TUI
func CreateStatusBar(backgroundColor tcell.Color) (*tview.Flex, *tview.TextView, *tview.TextView) {
	// Status text (left)
	statusText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(tcell.ColorWhite)

	// Status info (right)
	statusInfo := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetTextColor(tcell.ColorWhite)

	// Flex container
	statusBar := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(statusText, 0, 3, false).
		AddItem(statusInfo, 0, 1, false)

	statusBar.SetBackgroundColor(backgroundColor)

	return statusBar, statusText, statusInfo
}

// CreateMainLayout creates the main layout for the TUI
func CreateMainLayout(header, output, input, statusBar tview.Primitive) *tview.Flex {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, false).
		AddItem(output, 0, 1, false).
		AddItem(input, 1, 0, true).
		AddItem(statusBar, 1, 0, false)

	return layout
}

// CreateLoginForm creates the login form
func CreateLoginForm(loginFunc func()) *tview.Form {
	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Login", loginFunc).
		AddButton("Cancel", nil) // Will be set later

	form.SetBorder(true).
		SetTitle("Login").
		SetTitleAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorBlack)

	return form
}

// CreateModal creates a modal window
func CreateModal(title string, text string, buttons []string, callbacks []func()) *tview.Modal {
	modal := tview.NewModal()

	modal.SetTitle(title).
		SetText(text).
		SetBackgroundColor(tcell.ColorBlack)

	for i, button := range buttons {
		modal.AddButtons([]string{button})

		if i < len(callbacks) {
			modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonIndex >= 0 && buttonIndex < len(callbacks) {
					callbacks[buttonIndex]()
				}
			})
		}
	}

	return modal
}

// CenteredFlex centers a primitive on the screen
func CenteredFlex(p tview.Primitive, width, height int) tview.Primitive {
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
