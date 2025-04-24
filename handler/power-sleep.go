package event_handler

import (
	"fmt"
	"os/exec"
	"runtime"
)

func Sleep() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("systemctl", "suspend")
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("pmset", "sleepnow")
		return cmd.Run()
	default:
		return fmt.Errorf("sleeping not supported on %s", runtime.GOOS)
	}
}
