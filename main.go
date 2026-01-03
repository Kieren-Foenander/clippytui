package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// main is the entry point of every Go program.
// It's like the "if __name__ == '__main__'" block in Python.
func main() {
	// Create a new model with default values
	m := NewModel()

	// Create a new Bubble Tea program
	// WithAltScreen() gives us a full-screen terminal experience
	// (clears the screen and restores it when we quit)
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Start the clipboard monitor in a goroutine
	// This runs concurrently with the TUI, watching for clipboard changes
	go StartMonitor(p)

	// Try to get initial clipboard content
	// This populates the history with whatever is currently in the clipboard
	if initialContent, err := GetClipboardContent(); err == nil {
		if initialContent != "" {
			// Send initial content as a message
			p.Send(ClipboardMsg{Content: initialContent})
		}
	}

	// Run the program. This blocks until the user quits (presses 'q' or Ctrl+C)
	// It returns the final model and any error
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

