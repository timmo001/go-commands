package event_handler

import (
	"fmt"
	"os/exec"
	"runtime"
)

func Restart() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("shutdown", "/r")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("systemctl", "reboot")
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("osascript", "-e", "tell application \"System Events\" to restart")
		return cmd.Run()
	default:
		return fmt.Errorf("restarting not supported on %s", runtime.GOOS)
	}
}
