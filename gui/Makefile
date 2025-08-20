# Torrent VPN Manager GUI - Makefile
# Cross-platform build system for Go + Fyne application

# Application info
APP_NAME := torrent-vpn-manager
APP_ID := com.weydson.torrentvpn
APP_DISPLAY_NAME := TorrentVPNManager

# Build info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Directories
BUILD_DIR := build
ASSETS_DIR := assets
DIST_DIR := dist

# Go build flags
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.commit=$(COMMIT)
GO_BUILD_FLAGS := -ldflags="$(LDFLAGS)" -trimpath

# Platform detection
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# Platform-specific settings
ifeq ($(UNAME_S),Darwin)
    PLATFORM := darwin
    INSTALL_DIR := /Applications
    BINARY_EXT :=
    APP_EXT := .app
endif
ifeq ($(UNAME_S),Linux)
    PLATFORM := linux
    INSTALL_DIR := $(HOME)/.local/bin
    DESKTOP_DIR := $(HOME)/.local/share/applications
    ICON_DIR := $(HOME)/.local/share/icons/hicolor/256x256/apps
    BINARY_EXT :=
    APP_EXT :=
endif
ifeq ($(OS),Windows_NT)
    PLATFORM := windows
    INSTALL_DIR := $(USERPROFILE)/AppData/Local/Programs
    BINARY_EXT := .exe
    APP_EXT :=
endif

# Targets
.PHONY: all build run clean install uninstall package deps icon help

all: build

help: ## Show this help message
	@echo "Torrent VPN Manager GUI - Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

icon: $(ASSETS_DIR)/icon.png ## Generate PNG icon from SVG

$(ASSETS_DIR)/icon.png: $(ASSETS_DIR)/icon.svg
	@echo "🎨 Generating PNG icon..."
	go run genicon/main.go

run: deps ## Run the application directly
	@echo "🚀 Running application..."
	go run $(GO_BUILD_FLAGS) cmd/main.go

build: deps $(BUILD_DIR)/$(APP_NAME)$(BINARY_EXT) ## Build the application binary

$(BUILD_DIR)/$(APP_NAME)$(BINARY_EXT): cmd/main.go
	@echo "🔨 Building $(APP_NAME) for $(PLATFORM)..."
	@mkdir -p $(BUILD_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)$(BINARY_EXT) cmd/main.go

# Cross-platform builds
build-windows: ## Build for Windows
	@echo "🔨 Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows.exe cmd/main.go

build-linux: ## Build for Linux
	@echo "🔨 Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux cmd/main.go

build-mac-intel: ## Build for macOS Intel
	@echo "🔨 Building for macOS Intel..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-mac-intel cmd/main.go

build-mac-arm: ## Build for macOS Apple Silicon
	@echo "🔨 Building for macOS Apple Silicon..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-mac-arm64 cmd/main.go

build-all: build-windows build-linux build-mac-intel build-mac-arm ## Build for all platforms

# Platform-specific packaging
package: package-$(PLATFORM) ## Package for current platform

package-darwin: icon ## Package macOS .app bundle
	@echo "📦 Packaging macOS application..."
	@mkdir -p $(DIST_DIR)
	@command -v fyne >/dev/null 2>&1 || { echo "Installing fyne packager..."; go install fyne.io/tools/cmd/fyne@latest; }
	fyne package --target darwin --icon "$(PWD)/$(ASSETS_DIR)/icon.png" --name $(APP_DISPLAY_NAME) --app-id $(APP_ID) --src "$(PWD)/cmd" --release
	@mv "$(APP_DISPLAY_NAME).app" $(DIST_DIR)/ 2>/dev/null || true

package-linux: build icon ## Package Linux application with desktop entry
	@echo "📦 Packaging Linux application..."
	@mkdir -p $(DIST_DIR)/$(APP_NAME)
	@cp $(BUILD_DIR)/$(APP_NAME) $(DIST_DIR)/$(APP_NAME)/
	@cp $(ASSETS_DIR)/icon.png $(DIST_DIR)/$(APP_NAME)/
	@echo "[Desktop Entry]" > $(DIST_DIR)/$(APP_NAME)/$(APP_NAME).desktop
	@echo "Version=1.0" >> $(DIST_DIR)/$(APP_NAME)/$(APP_NAME).desktop
	@echo "Type=Application" >> $(DIST_DIR)/$(APP_NAME)/$(APP_NAME).desktop
	@echo "Name=$(APP_DISPLAY_NAME)" >> $(DIST_DIR)/$(APP_NAME)/$(APP_NAME).desktop
	@echo "Comment=Torrent VPN Manager GUI" >> $(DIST_DIR)/$(APP_NAME)/$(APP_NAME).desktop
	@echo "Exec=$(APP_NAME)" >> $(DIST_DIR)/$(APP_NAME)/$(APP_NAME).desktop
	@echo "Icon=$(APP_NAME)" >> $(DIST_DIR)/$(APP_NAME)/$(APP_NAME).desktop
	@echo "Categories=Network;FileTransfer;" >> $(DIST_DIR)/$(APP_NAME)/$(APP_NAME).desktop
	@echo "Terminal=false" >> $(DIST_DIR)/$(APP_NAME)/$(APP_NAME).desktop

package-windows: build icon ## Package Windows application
	@echo "📦 Packaging Windows application..."
	@mkdir -p $(DIST_DIR)/$(APP_NAME)
	@cp $(BUILD_DIR)/$(APP_NAME).exe $(DIST_DIR)/$(APP_NAME)/
	@cp $(ASSETS_DIR)/icon.png $(DIST_DIR)/$(APP_NAME)/

# Installation targets
install: install-$(PLATFORM) ## Install application for current platform

install-darwin: package-darwin ## Install on macOS
	@echo "📲 Installing $(APP_DISPLAY_NAME) to Applications..."
	@sudo rm -rf "$(INSTALL_DIR)/$(APP_DISPLAY_NAME).app"
	@sudo cp -R "$(DIST_DIR)/$(APP_DISPLAY_NAME).app" "$(INSTALL_DIR)/"
	@echo "✅ Installed to $(INSTALL_DIR)/$(APP_DISPLAY_NAME).app"
	@echo "💡 You may need to run 'killall Dock' to refresh the Dock icon"

refresh-dock: ## Refresh macOS Dock to show updated icons
	@echo "🔄 Refreshing Dock..."
	@killall Dock
	@echo "✅ Dock refreshed"

install-linux: package-linux ## Install on Linux
	@echo "📲 Installing $(APP_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR) $(DESKTOP_DIR) $(ICON_DIR)
	@cp $(DIST_DIR)/$(APP_NAME)/$(APP_NAME) $(INSTALL_DIR)/
	@chmod +x $(INSTALL_DIR)/$(APP_NAME)
	@cp $(DIST_DIR)/$(APP_NAME)/$(APP_NAME).desktop $(DESKTOP_DIR)/
	@cp $(DIST_DIR)/$(APP_NAME)/icon.png $(ICON_DIR)/$(APP_NAME).png
	@echo "✅ Installed to $(INSTALL_DIR)/$(APP_NAME)"
	@echo "💡 Desktop entry created at $(DESKTOP_DIR)/$(APP_NAME).desktop"

install-windows: package-windows ## Install on Windows
	@echo "📲 Installing $(APP_NAME) to $(INSTALL_DIR)..."
	@mkdir -p "$(INSTALL_DIR)/$(APP_NAME)"
	@cp "$(DIST_DIR)/$(APP_NAME)/$(APP_NAME).exe" "$(INSTALL_DIR)/$(APP_NAME)/"
	@cp "$(DIST_DIR)/$(APP_NAME)/icon.png" "$(INSTALL_DIR)/$(APP_NAME)/"
	@echo "✅ Installed to $(INSTALL_DIR)/$(APP_NAME)"

# Uninstall targets
uninstall: uninstall-$(PLATFORM) ## Uninstall application from current platform

uninstall-darwin: ## Uninstall from macOS
	@echo "🗑️  Uninstalling $(APP_DISPLAY_NAME)..."
	@sudo rm -rf "$(INSTALL_DIR)/$(APP_DISPLAY_NAME).app"
	@echo "✅ Uninstalled from macOS"

uninstall-linux: ## Uninstall from Linux
	@echo "🗑️  Uninstalling $(APP_NAME)..."
	@rm -f $(INSTALL_DIR)/$(APP_NAME)
	@rm -f $(DESKTOP_DIR)/$(APP_NAME).desktop
	@rm -f $(ICON_DIR)/$(APP_NAME).png
	@echo "✅ Uninstalled from Linux"

uninstall-windows: ## Uninstall from Windows
	@echo "🗑️  Uninstalling $(APP_NAME)..."
	@rm -rf "$(INSTALL_DIR)/$(APP_NAME)"
	@echo "✅ Uninstalled from Windows"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(ASSETS_DIR)/icon.png
	@rm -f "$(APP_DISPLAY_NAME).app" 2>/dev/null || true
	@echo "✅ Cleaned"

test: ## Run tests
	@echo "🧪 Running tests..."
	go test ./...

fmt: ## Format code
	@echo "🎨 Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	go vet ./...

lint: fmt vet ## Run linting tools

dev: ## Development mode (clean, build, run)
	$(MAKE) clean
	$(MAKE) run

# Version info
version: ## Show version information
	@echo "App: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Commit: $(COMMIT)"
	@echo "Platform: $(PLATFORM)"
