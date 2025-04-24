package event_handler

import (
	"fmt"
	"os/exec"
	"runtime"
)

func Lock() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("rundll32.exe", "user32.dll,LockWorkStation")
		return cmd.Run()
	case "linux":
		// Try Wayland first with loginctl
		cmd := exec.Command("loginctl", "lock-session")
		if err := cmd.Run(); err == nil {
			return nil
		}

		// If Wayland fails, try X11
		cmd = exec.Command("xscreensaver-command", "-lock")
		if err := cmd.Run(); err == nil {
			return nil
		}

		// If xscreensaver fails, try xlock as last resort
		cmd = exec.Command("xlock")
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("pmset", "displaysleepnow")
		return cmd.Run()
	default:
		return fmt.Errorf("locking not supported on %s", runtime.GOOS)
	}
}
