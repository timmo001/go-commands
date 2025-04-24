package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/timmo001/go-commands/mqtt"
)

func init() {
	log.SetLevel(log.DebugLevel)

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", "error", err)
	}
}

func main() {
	// Get MQTT configuration from environment variables
	mqttHost := os.Getenv("MQTT_HOST")
	mqttPort := os.Getenv("MQTT_PORT")
	mqttUser := os.Getenv("MQTT_USER")
	mqttPassword := os.Getenv("MQTT_PASSWORD")

	// Create MQTT broker URL
	brokerURL := fmt.Sprintf("tcp://%s:%s", mqttHost, mqttPort)

	// Create a new MQTT client
	client := mqtt.NewClient(brokerURL, mqttUser, mqttPassword)

	// Connect to the broker
	if err := client.Connect(); err != nil {
		log.Fatal("Failed to connect to MQTT broker", "error", err)
	}
	defer client.Disconnect()

	hostname := getHostname()
	deviceName := fmt.Sprintf("Go Commands - %s", hostname)
	uniqueID := fmt.Sprintf("go_commands_%s", hostname)
	baseTopic := fmt.Sprintf("go-commands/%s", uniqueID)

	// Publish the Home Assistant discovery message for the server status
	sensorConfig := map[string]interface{}{
		"name":               "Status",
		"unique_id":          fmt.Sprintf("%s_status", uniqueID),
		"state_topic":        fmt.Sprintf("%s/status", baseTopic),
		"availability_topic": fmt.Sprintf("%s/availability", baseTopic),
		"icon":               "mdi:server",
		"device": map[string]interface{}{
			"identifiers":  []string{uniqueID},
			"name":         deviceName,
			"model":        "Go Commands Service",
			"manufacturer": "Timmo",
		},
	}

	// Publish discovery configuration
	err := client.PublishDiscovery("sensor", uniqueID, "status", sensorConfig)
	if err != nil {
		log.Error("Failed to publish discovery message", "error", err)
	}

	// Publish initial availability
	err = client.Publish(fmt.Sprintf("%s/availability", baseTopic), 1, true, "online")
	if err != nil {
		log.Error("Failed to publish availability", "error", err)
	}

	// Start publishing status periodically
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			<-ticker.C
			err := client.Publish(fmt.Sprintf("%s/status", baseTopic), 1, false, "online")
			if err != nil {
				log.Error("Failed to publish status", "error", err)
			}
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal is received
	<-sigChan

	// Publish offline status before exiting
	err = client.Publish(fmt.Sprintf("%s/availability", baseTopic), 1, true, "offline")
	if err != nil {
		log.Error("Failed to publish offline status", "error", err)
	}

	log.Info("Shutting down...")
}
