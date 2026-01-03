package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Update is the "U" in MVU. It processes messages and updates the model.
// In Bubble Tea, every event (key press, clipboard change, etc.) becomes a message.
// Update receives these messages and returns an updated model and optionally a command.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Handle keyboard input
	case tea.KeyMsg:
		switch msg.String() {
		// Quit the application
		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		// Move cursor up (k or up arrow)
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
			m.status = "" // Clear status when navigating

		// Move cursor down (j or down arrow)
		case "j", "down":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
			m.status = "" // Clear status when navigating

		// Copy selected item to clipboard
		case "enter":
			if len(m.items) > 0 && m.cursor >= 0 && m.cursor < len(m.items) {
				// Get the selected item's content
				selectedItem := m.items[m.cursor]
				// Return updated model and a command to copy to clipboard
				// The command will send a CopiedMsg when done
				return m, CopyToClipboard(selectedItem.Content)
			}
		}

	// Handle terminal window resize
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update previews for all items with new width
		for i := range m.items {
			maxWidth := m.width - 10
			if maxWidth < 50 {
				maxWidth = 50
			}
			m.items[i].Preview = createPreview(m.items[i].Content, maxWidth)
		}

	// Handle new clipboard content detected
	case ClipboardMsg:
		// Calculate max width for preview
		maxWidth := m.width - 10
		if maxWidth < 50 {
			maxWidth = 50
		}
		// Add the new clipboard item to history
		// Since we're modifying the slice, we need to assign the returned value
		m.items = addClipboardItem(m.items, msg.Content, maxWidth)
		// Reset cursor to top (most recent item)
		m.cursor = 0
		return m, nil

	// Handle successful copy operation
	case CopiedMsg:
		m.status = "Copied!"
		// Clear status after a delay (we'll handle this with a timer command)
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
			return ClearStatusMsg{}
		})

	// Handle errors
	case ErrorMsg:
		m.status = "Error: " + msg.Err.Error()
		return m, tea.Tick(3*time.Second, func(time.Time) tea.Msg {
			return ClearStatusMsg{}
		})

	// Clear status message (sent after timer expires)
	case ClearStatusMsg:
		m.status = ""
		return m, nil
	}

	// Return the (possibly updated) model and no command
	return m, nil
}

// ClearStatusMsg is a message type to clear the status after a delay
type ClearStatusMsg struct{}

