package handler

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// PowerCommand represents a power-related command
type PowerCommand struct {
	Name        string
	Icon        string
	Description string
	Handler     func() error
}

// GetPowerCommands returns all available power commands
func GetPowerCommands() []PowerCommand {
	commands := []PowerCommand{
		{
			Name:        "Shutdown",
			Icon:        "mdi:power",
			Description: "Shutdown the system",
			Handler:     Shutdown,
		},
		{
			Name:        "Restart",
			Icon:        "mdi:restart",
			Description: "Restart the system",
			Handler:     Restart,
		},
		{
			Name:        "Sleep",
			Icon:        "mdi:power-sleep",
			Description: "Put the system to sleep",
			Handler:     Sleep,
		},
		{
			Name:        "Hibernate",
			Icon:        "mdi:power-sleep",
			Description: "Hibernate the system",
			Handler:     Hibernate,
		},
		{
			Name:        "Lock",
			Icon:        "mdi:lock",
			Description: "Lock the system",
			Handler:     Lock,
		},
		{
			Name:        "Logout",
			Icon:        "mdi:logout",
			Description: "Log out the current user",
			Handler:     Logout,
		},
	}

	if runtime.GOOS == "linux" {
		commands = append(commands, PowerCommand{
			Name:        "Restart to Windows",
			Icon:        "mdi:microsoft-windows",
			Description: "Restart the system to Windows",
			Handler:     RestartToWindows,
		})
	}

	return commands
}

// GetButtonConfig returns the Home Assistant button configuration for a power command
func GetButtonConfig(device map[string]interface{}, uniqueID string, baseTopic string, cmd PowerCommand) (string, map[string]interface{}) {
	nameAsId := strings.ReplaceAll(strings.ToLower(cmd.Name), " ", "_")
	return nameAsId, map[string]interface{}{
		"name":               fmt.Sprintf("%s", cmd.Name),
		"unique_id":          fmt.Sprintf("%s_power_%s", uniqueID, nameAsId),
		"command_topic":      fmt.Sprintf("%s/power/%s", baseTopic, nameAsId),
		"availability_topic": fmt.Sprintf("%s/availability", baseTopic),
		"icon":               cmd.Icon,
		"device":             device,
	}
}

// Shutdown shuts down the system
func Shutdown() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("shutdown", "/s")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("shutdown", "-h", "now")
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("osascript", "-e", "tell application \"System Events\" to shut down")
		return cmd.Run()
	default:
		return fmt.Errorf("shutdown not supported on %s", runtime.GOOS)
	}
}

// Restart restarts the system
func Restart() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("shutdown", "/r")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("shutdown", "-r", "now")
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("osascript", "-e", "tell application \"System Events\" to restart")
		return cmd.Run()
	default:
		return fmt.Errorf("restart not supported on %s", runtime.GOOS)
	}
}

// RestartToWindows restarts the system to Windows using the Windows Boot Manager efi entry
func RestartToWindows() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("restarting to Windows is only supported on Linux")
	}

	// Find the Windows Boot Manager entry
	cmd := exec.Command("efibootmgr")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run efibootmgr: %v", err)
	}

	// Parse the output to find Windows Boot Manager entry
	bootEntry := ""
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Windows Boot Manager") {
			// Extract the boot number (e.g., "Boot0001" -> "0001")
			bootEntry = line[4:8] // Assuming format is consistent
			break
		}
	}

	if bootEntry == "" {
		return fmt.Errorf("Windows Boot Manager not found")
	}

	// Set Windows Boot Manager as next boot option
	cmd = exec.Command("efibootmgr", "--bootnext", bootEntry)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set Windows Boot Manager as next boot option: %v", err)
	}

	// Reboot the system
	cmd = exec.Command("reboot")
	return cmd.Run()
}

// Sleep puts the system to sleep
func Sleep() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "0,1,0")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("systemctl", "suspend")
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("osascript", "-e", "tell application \"System Events\" to sleep")
		return cmd.Run()
	default:
		return fmt.Errorf("sleep not supported on %s", runtime.GOOS)
	}
}

// Hibernate hibernates the system
func Hibernate() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("shutdown", "/h")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("systemctl", "hibernate")
		return cmd.Run()
	case "darwin":
		return fmt.Errorf("hibernate not supported on macOS")
	default:
		return fmt.Errorf("hibernate not supported on %s", runtime.GOOS)
	}
}

// Lock locks the system
func Lock() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("rundll32.exe", "user32.dll,LockWorkStation")
		return cmd.Run()
	case "linux":
		// Try different commands for different desktop environments
		commands := [][]string{
			{"loginctl", "lock-session"},                                     // systemd
			{"gnome-screensaver-command", "-l"},                              // GNOME
			{"qdbus", "org.freedesktop.ScreenSaver", "/ScreenSaver", "Lock"}, // KDE
			{"xdg-screensaver", "lock"},                                      // Generic
		}

		for _, cmdArgs := range commands {
			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
			if err := cmd.Run(); err == nil {
				return nil
			}
		}
		return fmt.Errorf("failed to lock screen")
	case "darwin":
		cmd := exec.Command("osascript", "-e", "tell application \"System Events\" to keystroke \"q\" using {command down, control down}")
		return cmd.Run()
	default:
		return fmt.Errorf("lock not supported on %s", runtime.GOOS)
	}
}

// Logout logs out the current user
func Logout() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("shutdown", "/l")
		return cmd.Run()
	case "linux":
		// Try different commands for different desktop environments
		commands := [][]string{
			{"gnome-session-quit", "--no-prompt"},                                 // GNOME
			{"qdbus", "org.kde.ksmserver", "/KSMServer", "logout", "0", "0", "0"}, // KDE
			{"loginctl", "terminate-user", "$USER"},                               // systemd
		}

		for _, cmdArgs := range commands {
			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
			if err := cmd.Run(); err == nil {
				return nil
			}
		}
		return fmt.Errorf("failed to logout")
	case "darwin":
		cmd := exec.Command("osascript", "-e", "tell application \"System Events\" to log out")
		return cmd.Run()
	default:
		return fmt.Errorf("logout not supported on %s", runtime.GOOS)
	}
}
