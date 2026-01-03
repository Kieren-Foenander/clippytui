package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// StartMonitor starts monitoring the clipboard for changes.
// This function runs in a goroutine (think of it as a lightweight thread).
// Goroutines allow us to do things concurrently without blocking the main program.
func StartMonitor(program *tea.Program) {
	// #region agent log
	writeLog("clipboard.go:16", "StartMonitor called", map[string]interface{}{"program": "not nil"}, "H2")
	// #endregion

	// Start wl-paste --watch type in a subprocess
	// This command watches for clipboard changes and outputs the MIME type when it changes
	cmd := exec.Command("wl-paste", "--watch", "type")

	// #region agent log
	writeLog("clipboard.go:20", "Command created", map[string]interface{}{"command": "wl-paste --watch type"}, "H1")
	// #endregion

	// Get a pipe to read the command's stdout
	stdout, err := cmd.StdoutPipe()
	// #region agent log
	if err != nil {
		writeLog("clipboard.go:23", "StdoutPipe failed", map[string]interface{}{"error": err.Error()}, "H2")
	} else {
		writeLog("clipboard.go:23", "StdoutPipe success", map[string]interface{}{}, "H2")
	}
	// #endregion
	if err != nil {
		// If we can't start monitoring, send an error message to the program
		program.Send(ErrorMsg{Err: fmt.Errorf("failed to create stdout pipe: %w", err)})
		return
	}

	// Start the command (non-blocking)
	err = cmd.Start()
	// #region agent log
	if err != nil {
		writeLog("clipboard.go:30", "Command Start failed", map[string]interface{}{"error": err.Error()}, "H2")
	} else {
		writeLog("clipboard.go:30", "Command Start success", map[string]interface{}{"pid": cmd.Process.Pid}, "H2")
	}
	// #endregion
	if err != nil {
		program.Send(ErrorMsg{Err: fmt.Errorf("failed to start wl-paste: %w", err)})
		return
	}

	// Create a scanner to read line-by-line from stdout
	scanner := bufio.NewScanner(stdout)

	// #region agent log
	writeLog("clipboard.go:36", "Scanner created, entering scan loop", map[string]interface{}{}, "H3,H4")
	// #endregion

	// Read lines continuously
	// When wl-paste detects a clipboard change, it outputs a line
	scanCount := 0
	for scanner.Scan() {
		scanCount++
		line := scanner.Text()
		// #region agent log
		writeLog("clipboard.go:40", "Scanner read line", map[string]interface{}{"line": line, "scanCount": scanCount}, "H3,H4")
		// #endregion

		// When we detect a change, get the full clipboard content
		content, err := GetClipboardContent()
		// #region agent log
		if err != nil {
			writeLog("clipboard.go:42", "GetClipboardContent failed", map[string]interface{}{"error": err.Error()}, "H5")
		} else {
			writeLog("clipboard.go:42", "GetClipboardContent success", map[string]interface{}{"contentLength": len(content), "contentPreview": content[:min(50, len(content))]}, "H5")
		}
		// #endregion
		if err != nil {
			program.Send(ErrorMsg{Err: fmt.Errorf("failed to get clipboard content: %w", err)})
			continue
		}

		// Skip empty content
		trimmed := strings.TrimSpace(content)
		// #region agent log
		writeLog("clipboard.go:49", "Content check", map[string]interface{}{"trimmedLength": len(trimmed), "isEmpty": trimmed == ""}, "H8")
		// #endregion
		if trimmed == "" {
			continue
		}

		// Send a ClipboardMsg to the Bubble Tea program
		// This message will be received by the Update() function
		// #region agent log
		writeLog("clipboard.go:55", "Sending ClipboardMsg", map[string]interface{}{"contentLength": len(content)}, "H6")
		// #endregion
		program.Send(ClipboardMsg{Content: content})
	}

	// If the scanner encounters an error, report it
	err = scanner.Err()
	// #region agent log
	if err != nil {
		writeLog("clipboard.go:59", "Scanner error", map[string]interface{}{"error": err.Error(), "scanCount": scanCount}, "H4")
	} else {
		writeLog("clipboard.go:59", "Scanner loop exited", map[string]interface{}{"scanCount": scanCount}, "H4")
	}
	// #endregion
	if err != nil {
		program.Send(ErrorMsg{Err: fmt.Errorf("scanner error: %w", err)})
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

