// output.go
/**
* Nexuflex Client - Output Field Implementation
*
* This file contains extensions for the output field of the user interface.
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
	"github.com/rivo/tview"
)

// EnhancedTextView extends the standard TextView from tview
// with additional features like timestamps and formatting
type EnhancedTextView struct {
	*tview.TextView
	maxLines      int
	showTimestamp bool
	lineCount     int
}

// NewEnhancedTextView creates an enhanced output field
func NewEnhancedTextView(maxLines int, showTimestamp bool) *EnhancedTextView {
	output := &EnhancedTextView{
		TextView:      tview.NewTextView(),
		maxLines:      maxLines,
		showTimestamp: showTimestamp,
		lineCount:     0,
	}

	// Configure TextView
	output.
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)

	output.SetBorder(true).
		SetTitle("Output").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 1)

	return output
}

// WriteLine writes a line to the output field
func (o *EnhancedTextView) WriteLine(line string) {
	// Add timestamp if enabled
	if o.showTimestamp {
		timestamp := time.Now().Format("15:04:05")
		line = fmt.Sprintf("[gray]%s[white] %s", timestamp, line)
	}

	// Add line with line break
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}

	// Increment line counter
	o.lineCount++

	// Remove excess lines
	if o.maxLines > 0 && o.lineCount > o.maxLines {
		content := o.GetText(true)
		lines := strings.Split(content, "\n")

		// Calculate number of lines to remove
		removeCount := o.lineCount - o.maxLines
		if removeCount > len(lines) {
			removeCount = len(lines) - 1
		}

		// Remove oldest lines
		newContent := strings.Join(lines[removeCount:], "\n")
		o.SetText(newContent)

		// Adjust line counter
		o.lineCount -= removeCount
	}

	// Add line and scroll to end
	o.Write([]byte(line))
	row, _ := o.TextView.GetScrollOffset()
	_, _, _, height := o.TextView.GetInnerRect()
	o.TextView.ScrollTo(row+height, 0)
}

// WriteCommand writes a user-entered command to the output field
func (o *EnhancedTextView) WriteCommand(command string) {
	o.WriteLine(fmt.Sprintf("> [yellow]%s[white]", command))
}

// WriteError writes an error message to the output field
func (o *EnhancedTextView) WriteError(message string) {
	o.WriteLine(fmt.Sprintf("[red]Error: %s[white]", message))
}

// WriteSuccess writes a success message to the output field
func (o *EnhancedTextView) WriteSuccess(message string) {
	o.WriteLine(fmt.Sprintf("[green]%s[white]", message))
}

// WriteInfo writes an information message to the output field
func (o *EnhancedTextView) WriteInfo(message string) {
	o.WriteLine(fmt.Sprintf("[blue]%s[white]", message))
}

// WriteWarning writes a warning message to the output field
func (o *EnhancedTextView) WriteWarning(message string) {
	o.WriteLine(fmt.Sprintf("[yellow]%s[white]", message))
}

// ClearOutput clears the content of the output field
func (o *EnhancedTextView) ClearOutput() {
	o.SetText("")
	o.lineCount = 0
}

// SetMaxLines sets the maximum number of lines in the output field
func (o *EnhancedTextView) SetMaxLines(maxLines int) {
	o.maxLines = maxLines
}

// SetShowTimestamp enables or disables timestamp display
func (o *EnhancedTextView) SetShowTimestamp(show bool) {
	o.showTimestamp = show
}

// ScrollToTop scrolls to the top of the output field
func (o *EnhancedTextView) ScrollToTop() {
	o.ScrollTo(0, 0)
}

// ScrollToBottom scrolls to the bottom of the output field
func (o *EnhancedTextView) ScrollToBottom() {
	o.ScrollToHighlight()
}

// AddKeyboardHandlers adds keyboard handlers for scrolling
func (o *EnhancedTextView) AddKeyboardHandlers(inputCapture func(event *tcell.EventKey) *tcell.EventKey) {
	// Save previous handler
	prevHandler := o.GetInputCapture()

	// Set new handler that calls the previous one
	o.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyPgUp:
			// Page up
			row, _ := o.GetScrollOffset()
			_, _, _, height := o.GetInnerRect()
			o.ScrollTo(row-height, 0)
			return nil

		case tcell.KeyPgDn:
			// Page down
			row, _ := o.GetScrollOffset()
			_, _, _, height := o.GetInnerRect()
			o.ScrollTo(row+height, 0)
			return nil

		case tcell.KeyHome:
			// To top
			o.ScrollToTop()
			return nil

		case tcell.KeyEnd:
			// To bottom
			o.ScrollToBottom()
			return nil
		}

		// If a previous handler exists, call it
		if prevHandler != nil {
			return prevHandler(event)
		}

		// If an external handler exists, call it
		if inputCapture != nil {
			return inputCapture(event)
		}

		return event
	})
}
