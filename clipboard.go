package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// StartMonitor starts monitoring the clipboard for changes using polling.
// This function runs in a goroutine (think of it as a lightweight thread).
// Instead of using wl-paste --watch (which wasn't working), we poll the clipboard
// periodically to detect changes.
func StartMonitor(program *tea.Program) {
	// #region agent log
	writeLog("clipboard.go:16", "StartMonitor called (polling mode)", map[string]interface{}{"program": "not nil"}, "H2")
	// #endregion

	// Track the last clipboard content to detect changes
	lastContent := ""
	ticker := time.NewTicker(500 * time.Millisecond) // Check every 500ms
	defer ticker.Stop()

	// #region agent log
	writeLog("clipboard.go:24", "Starting polling loop", map[string]interface{}{"interval": "500ms"}, "H4")
	// #endregion

	// Poll the clipboard periodically
	for range ticker.C {
		// Get current clipboard content
		content, err := GetClipboardContent()
		// #region agent log
		if err != nil {
			// Only log errors occasionally to avoid spam
			if time.Now().Unix()%5 == 0 {
				writeLog("clipboard.go:32", "GetClipboardContent failed (polling)", map[string]interface{}{"error": err.Error()}, "H5")
			}
		}
		// #endregion
		if err != nil {
			// Don't spam errors - only send if it's a new type of error
			continue
		}

		// Check if content has changed
		if content != lastContent {
			// #region agent log
			writeLog("clipboard.go:42", "Clipboard content changed", map[string]interface{}{"oldLength": len(lastContent), "newLength": len(content)}, "H4")
			// #endregion

			// Skip empty content
			trimmed := strings.TrimSpace(content)
			// #region agent log
			writeLog("clipboard.go:47", "Content check", map[string]interface{}{"trimmedLength": len(trimmed), "isEmpty": trimmed == ""}, "H8")
			// #endregion

			// Update last content (even if empty, to avoid re-triggering)
			lastContent = content

			if trimmed != "" {
				// Send a ClipboardMsg to the Bubble Tea program
				// This message will be received by the Update() function
				// #region agent log
				writeLog("clipboard.go:55", "Sending ClipboardMsg (polling)", map[string]interface{}{"contentLength": len(content)}, "H6")
				// #endregion
				program.Send(ClipboardMsg{Content: content})
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetClipboardContent retrieves the current clipboard content using wl-paste.
// This is called when we detect a clipboard change to get the actual text.
func GetClipboardContent() (string, error) {
	// #region agent log
	writeLog("clipboard.go:66", "GetClipboardContent called", map[string]interface{}{}, "H5")
	// #endregion

	// Run wl-paste without --watch to get current clipboard content
	cmd := exec.Command("wl-paste", "--no-newline")
	// #region agent log
	writeLog("clipboard.go:69", "wl-paste command created", map[string]interface{}{"command": "wl-paste --no-newline"}, "H1,H5")
	// #endregion

	output, err := cmd.Output()
	// #region agent log
	if err != nil {
		writeLog("clipboard.go:70", "wl-paste Output failed", map[string]interface{}{"error": err.Error()}, "H1,H5")
	} else {
		writeLog("clipboard.go:70", "wl-paste Output success", map[string]interface{}{"outputLength": len(output)}, "H5")
	}
	// #endregion
	if err != nil {
		return "", fmt.Errorf("wl-paste failed: %w", err)
	}
	return string(output), nil
}

// CopyToClipboard copies the given content to the clipboard using wl-copy.
// This is called when the user presses Enter to copy a selected item.
func CopyToClipboard(content string) tea.Cmd {
	// Return a command function. Commands in Bubble Tea are functions that
	// perform async operations and return messages when done.
	return func() tea.Msg {
		// Create wl-copy command
		cmd := exec.Command("wl-copy")

		// Set stdin to the content we want to copy
		// In Go, we can set stdin using StdinPipe() or by writing to it
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to create stdin pipe: %w", err)}
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to start wl-copy: %w", err)}
		}

		// Write content to stdin
		if _, err := stdin.Write([]byte(content)); err != nil {
			stdin.Close()
			cmd.Wait()
			return ErrorMsg{Err: fmt.Errorf("failed to write to stdin: %w", err)}
		}

		// Close stdin to signal we're done writing
		stdin.Close()

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			return ErrorMsg{Err: fmt.Errorf("wl-copy failed: %w", err)}
		}

		// Return success message
		return CopiedMsg{}
	}
}

// createPreview creates a single-line preview of clipboard content.
// This is used to display items in the list without taking up too much space.
func createPreview(content string, maxWidth int) string {
	// Remove newlines and replace with spaces for single-line display
	preview := strings.ReplaceAll(content, "\n", " ")
	preview = strings.ReplaceAll(preview, "\r", " ")

	// Trim whitespace
	preview = strings.TrimSpace(preview)

	// If preview is too long, truncate it
	if len(preview) > maxWidth {
		preview = preview[:maxWidth-3] + "..."
	}

	return preview
}

// addClipboardItem creates a new ClipboardItem and adds it to the items slice.
// This is a helper function that returns the updated slice.
// In Go, slices are reference types, but when we append, we might get a new underlying array,
// so we return the updated slice to ensure the caller has the correct reference.
func addClipboardItem(items []ClipboardItem, content string, maxWidth int) []ClipboardItem {
	// Create preview
	if maxWidth < 50 {
		maxWidth = 50 // Minimum width
	}

	item := ClipboardItem{
		Content:   content,
		Timestamp: time.Now(),
		Preview:   createPreview(content, maxWidth),
	}

	// Add to the beginning of the slice (most recent first)
	// In Go, we use append() to add to slices
	items = append([]ClipboardItem{item}, items...)

	// Limit history to prevent memory issues (keep last 100 items)
	maxItems := 100
	if len(items) > maxItems {
		items = items[:maxItems]
	}

	return items
}

