package event_handler

import (
	"fmt"
	"os/exec"
	"runtime"
)

func Logout() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("shutdown", "/l")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("loginctl", "terminate-user", "current")
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("osascript", "-e", "tell application \"System Events\" to log out")
		return cmd.Run()
	default:
		return fmt.Errorf("logging out not supported on %s", runtime.GOOS)
	}
}
