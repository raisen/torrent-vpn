# GUI Application

Cross-platform GUI for the Torrent VPN Manager built with Go and Fyne.

## Requirements

- Go 1.19 or later
- Make (for build automation)
- Docker and Docker Compose (for the underlying services)
- Platform-specific GUI libraries (automatically handled by Fyne)

## Quick Start

```bash
# Clone and navigate to GUI directory
cd gui

# Install dependencies and run
make deps
make run
```

## Structure

```
gui/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   └── main_window.go       # Main GUI interface
├── docker/
│   └── manager.go           # Docker Compose management
├── assets/
│   ├── icon.go              # Embedded application icon
│   ├── icon.svg             # Application icon source (SVG)
│   └── icon.png             # Generated icon (PNG, ignored by git)
├── genicon/
│   └── main.go              # Icon generation utility
├── build/                   # Build output directory (ignored by git)
├── dist/                    # Distribution packages (ignored by git)
├── Makefile                 # Build automation
├── go.mod                   # GUI module dependencies
├── go.sum                   # Dependency checksums
└── README.md                # This file
```

## Building

From the GUI directory, use the provided Makefile:

```bash
# Run directly for development
make run

# Build executable for current platform
make build

# Cross-platform builds
make build-windows    # Windows executable
make build-linux      # Linux executable
make build-darwin     # macOS executable
make build-all        # All platforms

# Show all available targets
make help
```

## Packaging and Installation

Create platform-specific packages:

```bash
# Generate icon (PNG from SVG)
make icon

# Package for specific platforms
make package-darwin   # macOS .app bundle
make package-linux    # Linux AppImage (future)
make package-windows  # Windows installer (future)

# Install to system applications folder
make install          # Current platform
make uninstall        # Remove from system

# Clean build artifacts
make clean
```

## Development

```bash
# Install dependencies
make deps

# Format code
make fmt

# Run tests (if any)
make test
```

## Dependencies

### Runtime Dependencies
- `fyne.io/fyne/v2`: Cross-platform GUI framework
- `golang.org/x/image/vector` and `oksvg`: SVG processing for icon generation
- Local `docker` package: Docker Compose management

### Build Dependencies
- `fyne.io/tools/cmd/fyne`: Fyne packaging tool (auto-installed by Makefile)
- Platform-specific packaging tools (handled automatically)

### System Dependencies
- Docker and Docker Compose V2
- Platform GUI libraries (X11/Wayland on Linux, automatically available on macOS/Windows)

## Installation Locations

- **macOS**: `/Applications/TorrentVPNManager.app`
- **Linux**: `~/.local/share/applications/` (planned)
- **Windows**: `%PROGRAMFILES%\TorrentVPNManager\` (planned)
