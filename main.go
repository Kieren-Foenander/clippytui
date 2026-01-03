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
	// #region agent log
	writeLog("main.go:23", "Starting StartMonitor goroutine", map[string]interface{}{}, "H2")
	// #endregion
	go StartMonitor(p)

	// Try to get initial clipboard content
	// This populates the history with whatever is currently in the clipboard
	// #region agent log
	writeLog("main.go:27", "Getting initial clipboard content", map[string]interface{}{}, "H7")
	// #endregion
	initialContent, err := GetClipboardContent()
	// #region agent log
	if err != nil {
		writeLog("main.go:27", "Initial GetClipboardContent failed", map[string]interface{}{"error": err.Error()}, "H7")
	} else {
		writeLog("main.go:27", "Initial GetClipboardContent success", map[string]interface{}{"contentLength": len(initialContent), "isEmpty": initialContent == ""}, "H7")
	}
	// #endregion
	if err == nil {
		if initialContent != "" {
			// #region agent log
			writeLog("main.go:30", "Sending initial ClipboardMsg", map[string]interface{}{"contentLength": len(initialContent)}, "H6")
			// #endregion
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

