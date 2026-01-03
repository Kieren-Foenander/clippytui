# ClippyTUI - Wayland Clipboard Manager

A beautiful terminal-based clipboard manager for Wayland systems, built with Go and Bubble Tea.

## Features

- 📋 Automatic clipboard history tracking
- 🎨 Beautiful TUI matching Omarchy aesthetic
- ⌨️ Vim-style navigation (j/k, arrow keys)
- 🚀 Fast and lightweight
- 🔄 Real-time clipboard monitoring

## Requirements

- Go 1.21 or later
- `wl-clipboard` package (usually pre-installed on Wayland systems)
- A Wayland compositor (e.g., Sway, Hyprland, etc.)

## Installation

1. Clone this repository:

```bash
git clone <your-repo-url>
cd clippytui
```

2. Install dependencies:

```bash
go mod tidy
```

3. Build the application:

```bash
go build -o clippytui
```

4. Run it:

```bash
./clippytui
```

## Usage

- **↑/k**: Move cursor up
- **↓/j**: Move cursor down
- **Enter**: Copy selected item to clipboard
- **q**: Quit the application

## Architecture

The application follows the Model-View-Update (MVU) pattern:

- **Model** (`model.go`): Holds application state (clipboard history, cursor position)
- **View** (`view.go`): Renders the TUI based on model state
- **Update** (`update.go`): Handles messages (keyboard input, clipboard changes) and updates the model
- **Clipboard** (`clipboard.go`): Monitors clipboard using `wl-paste --watch` and handles copying

## Project Structure

```
clippytui/
├── go.mod          # Go module definition
├── main.go         # Entry point
├── model.go        # Data structures and Init()
├── update.go       # Message handling and state updates
├── view.go         # UI rendering
├── clipboard.go    # wl-clipboard integration
└── styles.go       # Lip Gloss styling definitions
```

## Future Improvements

- Custom C implementation for clipboard access (removing wl-clipboard dependency)
- Persistent history storage
- Search/filter functionality
- Image clipboard support
