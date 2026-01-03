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
	// Start wl-paste --watch type in a subprocess
	// This command watches for clipboard changes and outputs the MIME type when it changes
	cmd := exec.Command("wl-paste", "--watch", "type")

	// Get a pipe to read the command's stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		// If we can't start monitoring, send an error message to the program
		program.Send(ErrorMsg{Err: fmt.Errorf("failed to create stdout pipe: %w", err)})
		return
	}

	// Start the command (non-blocking)
	if err := cmd.Start(); err != nil {
		program.Send(ErrorMsg{Err: fmt.Errorf("failed to start wl-paste: %w", err)})
		return
	}

	// Create a scanner to read line-by-line from stdout
	scanner := bufio.NewScanner(stdout)

	// Read lines continuously
	// When wl-paste detects a clipboard change, it outputs a line
	for scanner.Scan() {
		// When we detect a change, get the full clipboard content
		content, err := GetClipboardContent()
		if err != nil {
			program.Send(ErrorMsg{Err: fmt.Errorf("failed to get clipboard content: %w", err)})
			continue
		}

		// Skip empty content
		if strings.TrimSpace(content) == "" {
			continue
		}

		// Send a ClipboardMsg to the Bubble Tea program
		// This message will be received by the Update() function
		program.Send(ClipboardMsg{Content: content})
	}

	// If the scanner encounters an error, report it
	if err := scanner.Err(); err != nil {
		program.Send(ErrorMsg{Err: fmt.Errorf("scanner error: %w", err)})
	}
}

// GetClipboardContent retrieves the current clipboard content using wl-paste.
// This is called when we detect a clipboard change to get the actual text.
func GetClipboardContent() (string, error) {
	// Run wl-paste without --watch to get current clipboard content
	cmd := exec.Command("wl-paste", "--no-newline")
	output, err := cmd.Output()
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

