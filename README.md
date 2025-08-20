# Torrent VPN Manager

A lightweight cross-platform GUI application built with Go and Fyne to manage your Docker-based torrent VPN solution.

![Application Icon](assets/icon.svg)

## Features

- **🚀 One-Click Control**: Start/stop your entire Docker Compose stack
- **📊 Real-time Monitoring**: Live status updates for all services
- **🌐 Quick Access**: Direct browser links to web interfaces
- **📝 Live Logs**: View container logs in real-time
- **🔒 VPN Status**: Monitor VPN connection and forwarded ports
- **💻 Cross-Platform**: Works on macOS, Windows, and Linux

## Support the Project

If you find this project helpful, consider sponsoring its development:

[![Sponsor](https://img.shields.io/badge/Sponsor-GitHub%20Sponsors-pink?style=for-the-badge&logo=github)](https://github.com/sponsors/raisen)

Your support helps maintain and improve this project! 🙏

## Services Managed

- **Gluetun VPN**: PIA VPN with port forwarding and kill switch
- **qBittorrent**: Torrent client with web interface (Port 8080)
- **Jellyfin**: Media server (Port 8096)
- **Port Updater**: Automatic port forwarding updates for qBittorrent
- **Streaming Optimizer**: Auto-enables sequential download for video files

## Quick Start

### Prerequisites
- Go 1.19 or later
- Docker and Docker Compose
- Private Internet Access (PIA) VPN account

### Setup

1. **Configure Docker Compose**:
   ```bash
   # Copy the example file
   cp docker-compose.example.yml docker-compose.yml

   # Edit docker-compose.yml and update:
   # - OPENVPN_USER: Your PIA username
   # - OPENVPN_PASSWORD: Your PIA password
   # - Volume paths for downloads directory
   ```

2. **Build and Run**:

```bash
# Download dependencies
go mod download

# Run directly (from gui directory)
cd gui && go run cmd/main.go

# Or build executable
cd gui && go build -o torrent-vpn-manager cmd/main.go
./torrent-vpn-manager
```

### Cross-Platform Building

```bash
# From gui directory
cd gui

# Windows
GOOS=windows GOARCH=amd64 go build -o torrent-vpn-manager.exe cmd/main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o torrent-vpn-manager-linux cmd/main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o torrent-vpn-manager-mac cmd/main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o torrent-vpn-manager-mac-arm64 cmd/main.go
```

## GUI Interface

### Main Window
- **System Status**: Overall health and control buttons
- **Services Grid**: Individual service status and health indicators
- **Web Interfaces**: Quick access to all web UIs with one-click browser opening
- **Logs**: Real-time application and container logs

### Quick Access Ports
- **qBittorrent**: http://localhost:8080 (admin/admin)
- **Jellyfin**: http://localhost:8096

## Architecture

```
gui/
  cmd/
    main.go            # GUI application entry point
  internal/
    main_window.go     # Main GUI interface
  docker/
    manager.go         # Docker Compose management
  assets/
    icon.go            # Embedded application icon
    icon.svg           # Application icon (SVG)
    icon.png           # Application icon (PNG for packaging)
  genicon/
    main.go            # Icon generation utility
  go.mod               # GUI module dependencies
dockerfiles/           # Custom Docker images
scripts/               # Docker helper scripts
docker-compose.yml     # Service definitions
```

## Why Fyne?

- **Lightweight**: ~5-10MB executable vs 100MB+ Electron apps
- **Native Performance**: True native GUI, not web-based
- **Cross-Platform**: Single codebase for all platforms
- **No Runtime**: Self-contained executable, no dependencies
- **Modern UI**: Clean, responsive interface

## Development

### Adding New Features
1. Service management logic goes in `gui/docker/manager.go`
2. UI components go in `gui/internal/main_window.go`
3. Build and test with `cd gui && go run cmd/main.go`

### Dependencies
- `fyne.io/fyne/v2`: Cross-platform GUI framework
- `gopkg.in/yaml.v3`: YAML parsing for docker-compose.yml

## Troubleshooting

### Common Issues
- **Docker not found**: Ensure Docker is installed and in PATH
- **Permission denied**: Run with appropriate Docker permissions
- **Ports in use**: Check if ports 8080, 8096 are available

### Logs
The application shows both GUI logs and container logs in the interface. For debugging, you can also check:
```bash
docker-compose logs -f
docker ps
```

## Additional Resources

- Gluetun documentation: https://github.com/qdm12/gluetun
- PIA setup guide: https://github.com/qdm12/gluetun-wiki/blob/main/setup/providers/privateinternetaccess.md
