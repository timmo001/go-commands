package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Client represents an MQTT client instance
type Client struct {
	client    MQTT.Client
	brokerURL string
	username  string
	password  string
	clientID  string
	connected bool
}

// NewClient creates a new MQTT client instance
func NewClient(brokerURL, username, password string) *Client {
	return &Client{
		brokerURL: brokerURL,
		username:  username,
		password:  password,
		clientID:  fmt.Sprintf("go-commands-%d", time.Now().Unix()),
	}
}

// Connect establishes connection to the MQTT broker
func (c *Client) Connect() error {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(c.brokerURL)
	opts.SetClientID(c.clientID)
	opts.SetUsername(c.username)
	opts.SetPassword(c.password)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(c.onConnect)
	opts.SetConnectionLostHandler(c.onConnectionLost)

	c.client = MQTT.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %v", token.Error())
	}

	return nil
}

// Disconnect cleanly disconnects from the MQTT broker
func (c *Client) Disconnect() {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250)
	}
}

// Publish sends a message to a specific topic with QoS
func (c *Client) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if !c.client.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	var payloadBytes []byte
	switch p := payload.(type) {
	case string:
		payloadBytes = []byte(p)
	case []byte:
		payloadBytes = p
	default:
		var err error
		payloadBytes, err = json.Marshal(p)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %v", err)
		}
	}

	token := c.client.Publish(topic, qos, retained, payloadBytes)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish message: %v", token.Error())
	}

	log.Info("Published message", "topic", topic, "payload", payload)

	return nil
}

// PublishDiscovery publishes a Home Assistant discovery message
func (c *Client) PublishDiscovery(component, nodeID, objectID string, config interface{}) error {
	topic := fmt.Sprintf("homeassistant/%s/%s/%s/config", component, nodeID, objectID)
	err := c.Publish(topic, 1, true, config)
	if err == nil {
		log.Info("Published discovery message", "topic", topic)
	}
	return err
}

// Subscribe subscribes to a topic with specified QoS and message handler
func (c *Client) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) error {
	if token := c.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic: %v", token.Error())
	}
	log.Info("Subscribed to topic", "topic", topic)
	return nil
}

func (c *Client) onConnect(client MQTT.Client) {
	c.connected = true
	log.Info("Connected to MQTT broker", "broker", c.brokerURL)
}

func (c *Client) onConnectionLost(client MQTT.Client, err error) {
	c.connected = false
	log.Error("Lost connection to MQTT broker", "error", err)
}

// IsConnected returns the current connection status
func (c *Client) IsConnected() bool {
	return c.connected && c.client.IsConnected()
}
