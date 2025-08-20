package internal

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"torrent-vpn-gui/docker"
)

type MainWindow struct {
	dockerManager *docker.Manager
	statusCard    *widget.Card
	servicesCard  *widget.Card
	portsCard     *widget.Card
	logsCard      *widget.Card

	startBtn      *widget.Button
	stopBtn       *widget.Button
	statusLabel   *widget.Label
	servicesGrid  *fyne.Container
	portsGrid     *fyne.Container
	logsEntry     *widget.Entry
}

func NewMainWindow(dockerManager *docker.Manager) fyne.CanvasObject {
	mw := &MainWindow{
		dockerManager: dockerManager,
	}

	return mw.buildUI()
}

func (mw *MainWindow) buildUI() fyne.CanvasObject {
	// Control buttons
	mw.startBtn = widget.NewButton("Start Services", mw.startServices)
	mw.startBtn.Importance = widget.HighImportance

	mw.stopBtn = widget.NewButton("Stop Services", mw.stopServices)
	mw.stopBtn.Importance = widget.MediumImportance

	refreshBtn := widget.NewButton("Refresh Status", mw.refreshStatus)

	controlsContainer := container.NewHBox(
		mw.startBtn,
		mw.stopBtn,
		refreshBtn,
	)

	// Status section
	mw.statusLabel = widget.NewLabel("Status: Unknown")
	mw.statusCard = widget.NewCard("System Status", "", container.NewVBox(
		mw.statusLabel,
		controlsContainer,
	))

	// Services status grid (using labels instead of table)
	mw.servicesGrid = mw.createServicesGrid()
	mw.servicesCard = widget.NewCard("Services", "", mw.servicesGrid)

	// Available ports grid (using buttons instead of table)
	mw.portsGrid = mw.createPortsGrid()
	mw.portsCard = widget.NewCard("Available Web Interfaces", "", mw.portsGrid)

	// Logs viewer with larger size
	mw.logsEntry = widget.NewMultiLineEntry()
	mw.logsEntry.SetText("Application logs will appear here...")
	mw.logsEntry.Wrapping = fyne.TextWrapWord
	logsScroll := container.NewScroll(mw.logsEntry)
	logsScroll.Resize(fyne.NewSize(800, 200)) // Set explicit height
	mw.logsCard = widget.NewCard("Logs", "", logsScroll)

	// Layout with more space for logs
	topRow := container.NewGridWithColumns(2, mw.statusCard, mw.servicesCard)
	middleRow := mw.portsCard
	bottomRow := mw.logsCard

	// Use border layout to give logs more space
	upperContent := container.NewVBox(topRow, middleRow)
	content := container.NewBorder(
		upperContent, // top
		nil,         // bottom
		nil,         // left
		nil,         // right
		bottomRow,   // center (logs will get remaining space)
	)

	// Start periodic status updates
	go mw.periodicStatusUpdate()

	// Initial status check
	mw.refreshStatus()

	return content
}

func (mw *MainWindow) createServicesGrid() *fyne.Container {
	// Header row
	header := container.NewGridWithColumns(3,
		widget.NewLabelWithStyle("Service", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Status", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Health", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	// Service rows
	services := []string{"Gluetun VPN", "qBittorrent", "Jellyfin", "Port Updater", "Streaming Optimizer"}
	var rows []fyne.CanvasObject
	rows = append(rows, header)

	for _, service := range services {
		row := container.NewGridWithColumns(3,
			widget.NewLabel(service),
			widget.NewLabel("Unknown"),
			widget.NewLabel("Unknown"),
		)
		rows = append(rows, row)
	}

	return container.NewVBox(rows...)
}

func (mw *MainWindow) createPortsGrid() *fyne.Container {
	// Header row
	header := container.NewGridWithColumns(4,
		widget.NewLabelWithStyle("Service", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Port", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("URL", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Action", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	// Port rows
	ports := [][]string{
		{"qBittorrent", "8080", "http://localhost:8080"},
		{"Jellyfin", "8096", "http://localhost:8096"},
	}

	var rows []fyne.CanvasObject
	rows = append(rows, header)

	for _, port := range ports {
		// Capture the URL in a local variable to avoid closure issues
		url := port[2]
		openBtn := widget.NewButton("Open", func() {
			mw.openBrowser(url)
		})

		row := container.NewGridWithColumns(4,
			widget.NewLabel(port[0]),
			widget.NewLabel(port[1]),
			widget.NewLabel(port[2]),
			openBtn,
		)
		rows = append(rows, row)
	}

	return container.NewVBox(rows...)
}

func (mw *MainWindow) startServices() {
	mw.addLog("Starting Docker Compose services...")
	mw.startBtn.Disable()

	go func() {
		err := mw.dockerManager.Start()
		if err != nil {
			mw.addLog(fmt.Sprintf("Error starting services: %v", err))
		} else {
			mw.addLog("Services started successfully!")
		}
		mw.startBtn.Enable()
		mw.refreshStatus()
	}()
}

func (mw *MainWindow) stopServices() {
	mw.addLog("Stopping Docker Compose services...")
	mw.stopBtn.Disable()

	go func() {
		err := mw.dockerManager.Stop()
		if err != nil {
			mw.addLog(fmt.Sprintf("Error stopping services: %v", err))
		} else {
			mw.addLog("Services stopped successfully!")
		}
		mw.stopBtn.Enable()
		mw.refreshStatus()
	}()
}

func (mw *MainWindow) refreshStatus() {
	go func() {
		status := mw.dockerManager.GetStatus()
		mw.statusLabel.SetText(fmt.Sprintf("Status: %s", status.Overall))

		// Update services grid with real data
		if len(mw.servicesGrid.Objects) > 1 { // Skip header row
			for i, service := range status.Services {
				if i+1 < len(mw.servicesGrid.Objects) {
					row := mw.servicesGrid.Objects[i+1].(*fyne.Container)
					if len(row.Objects) >= 3 {
						// Update status and health labels
						row.Objects[1].(*widget.Label).SetText(strings.Title(service.Status))
						row.Objects[2].(*widget.Label).SetText(strings.Title(service.Health))
					}
				}
			}
			mw.servicesGrid.Refresh()
		}

		// Update web interfaces grid with real data
		if len(mw.portsGrid.Objects) > 1 { // Skip header row
			// The ports grid shows static information, so we don't need to update it
			// but we could add VPN port information here if needed
		}

		mw.addLog(fmt.Sprintf("Status refreshed: %s", status.Overall))
	}()
}

func (mw *MainWindow) periodicStatusUpdate() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		mw.refreshStatus()
	}
}

func (mw *MainWindow) addLog(message string) {
	timestamp := time.Now().Format("15:04:05")
	logLine := fmt.Sprintf("[%s] %s\n", timestamp, message)

	currentText := mw.logsEntry.Text
	newText := currentText + logLine

	// Keep only last 25 lines to prevent overflow without scrollbars
	lines := strings.Split(newText, "\n")
	if len(lines) > 25 {
		lines = lines[len(lines)-25:]
		newText = strings.Join(lines, "\n")
	}

	mw.logsEntry.SetText(newText)

	// Keep cursor at end
	mw.logsEntry.CursorRow = len(lines) - 1
}

func (mw *MainWindow) openBrowser(urlStr string) {
	if u, err := url.Parse(urlStr); err == nil {
		fyne.CurrentApp().OpenURL(u)
		mw.addLog(fmt.Sprintf("Opening browser: %s", urlStr))
	} else {
		mw.addLog(fmt.Sprintf("Invalid URL: %s", urlStr))
	}
}
