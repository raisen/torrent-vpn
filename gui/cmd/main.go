package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"torrent-vpn-gui/assets"
	"torrent-vpn-gui/docker"
	"torrent-vpn-gui/internal"
)

func main() {
	myApp := app.New()
	myApp.SetIcon(assets.IconResource)
	myWindow := myApp.NewWindow("Torrent VPN Manager")
	myWindow.Resize(fyne.NewSize(1000, 750)) // Increased width to accommodate larger tables

	dockerManager := docker.NewManager("..")
	content := internal.NewMainWindow(dockerManager)
	myWindow.SetContent(content)

	myWindow.ShowAndRun()
}
