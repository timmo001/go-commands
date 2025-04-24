package event_handler

import (
	"fmt"
	"os/exec"
	"runtime"
)

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
