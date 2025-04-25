package handler

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// MediaCommand represents a media control command
type MediaCommand struct {
	Name        string
	Icon        string
	Description string
	Handler     func() error
}

// GetMediaCommands returns all available media control commands
func GetMediaCommands() []MediaCommand {
	return []MediaCommand{
		{
			Name:        "Play/Pause",
			Icon:        "mdi:play-pause",
			Description: "Toggle media playback",
			Handler:     PlayPause,
		},
		{
			Name:        "Next Track",
			Icon:        "mdi:skip-next",
			Description: "Play next track",
			Handler:     NextTrack,
		},
		{
			Name:        "Previous Track",
			Icon:        "mdi:skip-previous",
			Description: "Play previous track",
			Handler:     PreviousTrack,
		},
		{
			Name:        "Volume Up",
			Icon:        "mdi:volume-plus",
			Description: "Increase volume",
			Handler:     VolumeUp,
		},
		{
			Name:        "Volume Down",
			Icon:        "mdi:volume-minus",
			Description: "Decrease volume",
			Handler:     VolumeDown,
		},
		{
			Name:        "Mute",
			Icon:        "mdi:volume-mute",
			Description: "Toggle mute",
			Handler:     ToggleMute,
		},
	}
}

// GetMediaButtonConfig returns the Home Assistant button configuration for a media command
func GetMediaButtonConfig(device map[string]any, uniqueID string, baseTopic string, cmd MediaCommand) (string, map[string]interface{}) {
	nameAsId := strings.ReplaceAll(strings.ToLower(cmd.Name), " ", "_")
	nameAsId = strings.ReplaceAll(nameAsId, "/", "_")
	return nameAsId, map[string]any{
		"name":               fmt.Sprintf("%s", cmd.Name),
		"unique_id":          fmt.Sprintf("%s_media_%s", uniqueID, nameAsId),
		"command_topic":      fmt.Sprintf("%s/media/%s", baseTopic, nameAsId),
		"availability_topic": fmt.Sprintf("%s/availability", baseTopic),
		"icon":               cmd.Icon,
		"enabled_by_default": true, // Set to false to disable by default in Home Assistant
		"device":             device,
	}
}

// PlayPause toggles media playback
func PlayPause() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]179)")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("dbus-send", "--type=method_call", "--dest=org.mpris.MediaPlayer2.playerctld", "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2.Player.PlayPause")
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("osascript", "-e", "tell application \"System Events\" to key code 16 using {command down}")
		return cmd.Run()
	default:
		return fmt.Errorf("media control not supported on %s", runtime.GOOS)
	}
}

// NextTrack plays the next track
func NextTrack() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]176)")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("dbus-send", "--type=method_call", "--dest=org.mpris.MediaPlayer2.playerctld", "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2.Player.Next")
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("osascript", "-e", "tell application \"System Events\" to key code 17 using {command down}")
		return cmd.Run()
	default:
		return fmt.Errorf("media control not supported on %s", runtime.GOOS)
	}
}

// PreviousTrack plays the previous track
func PreviousTrack() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]177)")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("dbus-send", "--type=method_call", "--dest=org.mpris.MediaPlayer2.playerctld", "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2.Player.Previous")
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("osascript", "-e", "tell application \"System Events\" to key code 16 using {command down}")
		return cmd.Run()
	default:
		return fmt.Errorf("media control not supported on %s", runtime.GOOS)
	}
}

// VolumeUp increases the system volume
func VolumeUp() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]175)")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", "+5%")
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("osascript", "-e", "set volume output volume (output volume of (get volume settings) + 6)")
		return cmd.Run()
	default:
		return fmt.Errorf("volume control not supported on %s", runtime.GOOS)
	}
}

// VolumeDown decreases the system volume
func VolumeDown() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]174)")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", "-5%")
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("osascript", "-e", "set volume output volume (output volume of (get volume settings) - 6)")
		return cmd.Run()
	default:
		return fmt.Errorf("volume control not supported on %s", runtime.GOOS)
	}
}

// ToggleMute toggles system mute state
func ToggleMute() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]173)")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("pactl", "set-sink-mute", "@DEFAULT_SINK@", "toggle")
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("osascript", "-e", "set volume with output muted (not output muted of (get volume settings))")
		return cmd.Run()
	default:
		return fmt.Errorf("mute control not supported on %s", runtime.GOOS)
	}
}
