package docker

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Manager struct {
	projectDir string
}

// Helper function to detect available docker compose command
func (m *Manager) getComposeCommand() []string {
	// Check if modern 'docker compose' is available
	cmd := exec.Command("docker", "compose", "version")
	if err := cmd.Run(); err == nil {
		return []string{"docker", "compose"}
	}

	// Fallback to legacy 'docker-compose'
	return []string{"docker-compose"}
}

type ServiceStatus struct {
	Name   string
	Status string
	Health string
}

type SystemStatus struct {
	Overall  string
	Services []ServiceStatus
	VPNInfo  VPNInfo
}

type VPNInfo struct {
	Connected     bool
	IP            string
	ForwardedPort string
	Location      string
}

func NewManager(projectDir string) *Manager {
	return &Manager{
		projectDir: projectDir,
	}
}

func (m *Manager) Start() error {
	composeCmd := m.getComposeCommand()
	args := append(composeCmd, "up", "-d", "--build")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = m.projectDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start services: %v\nOutput: %s", err, string(output))
	}

	return nil
}

func (m *Manager) Stop() error {
	composeCmd := m.getComposeCommand()
	args := append(composeCmd, "down")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = m.projectDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop services: %v\nOutput: %s", err, string(output))
	}

	return nil
}

func (m *Manager) GetStatus() SystemStatus {
	services := m.getServiceStatuses()
	vpnInfo := m.getVPNInfo()

	overall := "Stopped"
	runningCount := 0

	for _, service := range services {
		if service.Status == "running" {
			runningCount++
		}
	}

	if runningCount == len(services) {
		overall = "All Running"
	} else if runningCount > 0 {
		overall = fmt.Sprintf("Partial (%d/%d)", runningCount, len(services))
	}

	return SystemStatus{
		Overall:  overall,
		Services: services,
		VPNInfo:  vpnInfo,
	}
}

func (m *Manager) getServiceStatuses() []ServiceStatus {
	services := []ServiceStatus{
		{"gluetun", "unknown", "unknown"},
		{"qbittorrent", "unknown", "unknown"},
		{"jellyfin", "unknown", "unknown"},
		{"qbit-port-updater", "unknown", "unknown"},
		{"qbit-streaming-optimizer", "unknown", "unknown"},
	}

	composeCmd := m.getComposeCommand()
	args := append(composeCmd, "ps")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = m.projectDir

	output, err := cmd.Output()
	if err != nil {
		return services
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Look for service names in the output
		for i := range services {
			serviceName := services[i].Name
			// The service name appears in the middle of the line after container name
			if strings.Contains(line, serviceName) && strings.Contains(line, "Up ") {
				services[i].Status = "running"
				// Check for health status
				if strings.Contains(line, "(healthy)") {
					services[i].Health = "healthy"
				} else {
					services[i].Health = "running"
				}
			} else if strings.Contains(line, serviceName) && strings.Contains(line, "Exit") {
				services[i].Status = "stopped"
				services[i].Health = "stopped"
			}
		}
	}

	return services
}

func (m *Manager) getVPNInfo() VPNInfo {
	info := VPNInfo{
		Connected:     false,
		IP:            "Unknown",
		ForwardedPort: "Unknown",
		Location:      "Unknown",
	}

	// Check if gluetun container is running
	cmd := exec.Command("docker", "ps", "--filter", "name=gluetun", "--format", "{{.Status}}")
	output, err := cmd.Output()
	if err != nil || !strings.Contains(string(output), "Up") {
		return info
	}

	info.Connected = true
	info.Location = "Netherlands (PIA)" // From your docker-compose.yml

	// Try to get forwarded port from the volume
	portFile := filepath.Join(m.projectDir, "gluetun_data", "forwarded_port")
	if data, err := os.ReadFile(portFile); err == nil {
		info.ForwardedPort = strings.TrimSpace(string(data))
	}

	// Get external IP
	cmd = exec.Command("docker", "exec", "torrent-vpn-gluetun-1", "curl", "-s", "ifconfig.me")
	if output, err := cmd.Output(); err == nil {
		info.IP = strings.TrimSpace(string(output))
	}

	return info
}

func (m *Manager) GetLogs(serviceName string) ([]string, error) {
	containerName := fmt.Sprintf("torrent-vpn-%s-1", serviceName)
	cmd := exec.Command("docker", "logs", "--tail", "50", containerName)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func (m *Manager) IsServiceHealthy(serviceName string) bool {
	// This could be enhanced with actual health check logic
	// For now, just check if container is running
	containerName := fmt.Sprintf("torrent-vpn-%s-1", serviceName)
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s", containerName), "--format", "{{.Status}}")

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), "Up")
}

func (m *Manager) GetVPNStatus() (bool, string, error) {
	// Check VPN connection by testing external IP
	cmd := exec.Command("docker", "exec", "torrent-vpn-gluetun-1", "curl", "-s", "--max-time", "5", "ifconfig.me")

	output, err := cmd.Output()
	if err != nil {
		return false, "", err
	}

	ip := strings.TrimSpace(string(output))

	// A very basic check - if we get an IP, assume VPN is working
	// You could enhance this by checking if the IP is from PIA Netherlands
	connected := len(ip) > 0 && !strings.Contains(ip, "curl:")

	return connected, ip, nil
}
