package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ClipboardItem represents a single entry in the clipboard history.
// In Go, structs are like classes but simpler - they just hold data.
// We'll add methods to them later (like String() for display).
type ClipboardItem struct {
	Content   string    // The actual clipboard text content
	Timestamp time.Time // When this item was captured
	Preview   string    // Truncated single-line preview for the list view
}

// Model holds all the application state. This is the "M" in MVU (Model-View-Update).
// In Bubble Tea, the Model must implement the tea.Model interface, which requires
// three methods: Init(), Update(), and View().
type Model struct {
	items  []ClipboardItem // Slice (like an array) of clipboard history items
	cursor int             // Index of currently selected item (0-based)
	width  int             // Terminal width for responsive layout
	height int             // Terminal height for responsive layout
	status string          // Status message (e.g., "Copied!" after selecting an item)
}

// Init is called when the program starts. It returns an initial command to run.
// Commands in Bubble Tea are asynchronous operations that return messages.
// Here we return nil because we don't need to do anything immediately.
// The clipboard monitoring will be started separately in main.go.
func (m Model) Init() tea.Cmd {
	return nil
}

// These are custom message types. In Bubble Tea, messages are how different
// parts of your program communicate. When something happens (like clipboard
// changes or user input), a message is sent to Update().

// ClipboardMsg is sent when new clipboard content is detected
type ClipboardMsg struct {
	Content string
}

// CopiedMsg is sent after successfully copying an item back to clipboard
type CopiedMsg struct{}

// ErrorMsg is sent when something goes wrong
type ErrorMsg struct {
	Err error
}

