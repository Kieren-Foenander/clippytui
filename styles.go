package main

import (
	"github.com/charmbracelet/lipgloss"
)

// Lip Gloss is a library for styling terminal output. It's like CSS for the terminal.
// We define styles as variables that can be reused throughout the application.

var (
	// BorderStyle defines the box-drawing characters for borders
	// This creates the nice box around our UI
	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")) // Purple border color

	// HeaderStyle styles the title bar
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("62")).
			Padding(0, 1)

	// ItemStyle styles regular list items
	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")) // Light gray text

	// SelectedStyle styles the currently selected item
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")). // Bright yellow
			Bold(true)

	// CursorStyle is the arrow indicator (▸) for selected items
	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Bold(true)

	// FooterStyle styles the help text at the bottom
	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")). // Dark gray
			Align(lipgloss.Center)

	// StatusStyle styles status messages (like "Copied!")
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")). // Green
			Bold(true)

	// ErrorStyle styles error messages
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")). // Red
			Bold(true)
)

// renderBox wraps content in a bordered box
func renderBox(content string, width, height int) string {
	return borderStyle.
		Width(width - 2).  // Subtract 2 for border
		Height(height - 2). // Subtract 2 for border
		Render(content)
}

