package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View is the "V" in MVU. It renders the UI based on the current model state.
// In Bubble Tea, View returns a string that gets printed to the terminal.
// Lip Gloss helps us style this string with colors and borders.
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	// Calculate available space for content
	// We need space for: header (1 line), footer (1 line), borders (2 lines)
	contentHeight := m.height - 4
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Build the header
	header := headerStyle.Render("📋 Clipboard Manager")
	header = lipgloss.PlaceHorizontal(m.width-4, lipgloss.Left, header)

	// Build the list of items
	var listItems []string

	// Determine which items to show (handle scrolling)
	startIdx := 0
	endIdx := len(m.items)

	// If we have more items than can fit, show a window around the cursor
	if len(m.items) > contentHeight {
		// Center the cursor in the visible area
		halfHeight := contentHeight / 2
		startIdx = m.cursor - halfHeight
		if startIdx < 0 {
			startIdx = 0
		}
		endIdx = startIdx + contentHeight
		if endIdx > len(m.items) {
			endIdx = len(m.items)
			startIdx = endIdx - contentHeight
			if startIdx < 0 {
				startIdx = 0
			}
		}
	}

	// Render each visible item
	for i := startIdx; i < endIdx; i++ {
		item := m.items[i]
		var line string

		// Add cursor indicator (▸) if this is the selected item
		if i == m.cursor {
			line = cursorStyle.Render("▸") + " "
			// Apply selected style to the item text
			line += selectedStyle.Render(item.Preview)
		} else {
			line = "  " // Two spaces for alignment when no cursor
			line += itemStyle.Render(item.Preview)
		}

		// Truncate if too long (account for cursor/indent)
		maxLineWidth := m.width - 6 // Account for borders and padding
		if len(line) > maxLineWidth {
			line = line[:maxLineWidth-3] + "..."
		}

		listItems = append(listItems, line)
	}

	// If no items yet, show a message
	if len(listItems) == 0 {
		listItems = append(listItems, itemStyle.Render("  No clipboard history yet. Copy something to get started!"))
	}

	// Join list items with newlines
	listContent := strings.Join(listItems, "\n")

	// Add status message if present
	if m.status != "" {
		var statusText string
		if strings.HasPrefix(m.status, "Error:") {
			statusText = errorStyle.Render(m.status)
		} else {
			statusText = statusStyle.Render(m.status)
		}
		// Place status at the bottom of the list area
		listContent += "\n\n" + lipgloss.PlaceHorizontal(m.width-4, lipgloss.Center, statusText)
	}

	// Build the footer with help text
	footer := footerStyle.Render("↑/k Up  ↓/j Down  Enter Copy  q Quit")

	// Combine all parts
	content := header + "\n\n" + listContent + "\n\n" + footer

	// Wrap everything in a bordered box
	return renderBox(content, m.width, m.height)
}

// Helper function to ensure we have a valid model with default values
func NewModel() Model {
	return Model{
		items:  []ClipboardItem{},
		cursor: 0,
		width:  80,  // Default width
		height: 24,  // Default height
		status: "",
	}
}

