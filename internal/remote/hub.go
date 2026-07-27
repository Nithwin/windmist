package remote

import (
	"fmt"
	"log"
	"sync"
	"github.com/Nithwin/WindMist/internal/config"
)

// Controller interface defines the standard operations for any remote interface (Telegram, Web, WhatsApp).
type Controller interface {
	Start(hub *Hub) error
	Stop() error
	SendMessage(text string) error
	Name() string
}

// Hub manages all active remote controllers.
type Hub struct {
	mu          sync.Mutex
	config      *config.RemoteConfig
	controllers map[string]Controller
	// Channel to receive messages from the agent to broadcast
	Broadcast chan string
	stopChan  chan struct{}
}

var globalHub *Hub

// InitHub initializes the global remote hub.
func InitHub(cfg *config.RemoteConfig) *Hub {
	globalHub = &Hub{
		config:      cfg,
		controllers: make(map[string]Controller),
		Broadcast:   make(chan string, 100),
		stopChan:    make(chan struct{}),
	}
	go globalHub.listen()
	return globalHub
}

// GetHub returns the active global hub.
func GetHub() *Hub {
	return globalHub
}

// Register adds a new controller to the hub and starts it.
func (h *Hub) Register(c Controller) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.controllers[c.Name()]; exists {
		return fmt.Errorf("controller %s is already registered", c.Name())
	}

	err := c.Start(h)
	if err != nil {
		return err
	}

	h.controllers[c.Name()] = c
	return nil
}

// Unregister stops and removes a controller.
func (h *Hub) Unregister(name string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	c, exists := h.controllers[name]
	if !exists {
		return fmt.Errorf("controller %s not found", name)
	}

	err := c.Stop()
	delete(h.controllers, name)
	return err
}

// listen waits for broadcast messages and sends them to all registered controllers.
func (h *Hub) listen() {
	for {
		select {
		case msg := <-h.Broadcast:
			h.mu.Lock()
			for name, c := range h.controllers {
				if err := c.SendMessage(msg); err != nil {
					log.Printf("Error sending message to %s: %v", name, err)
				}
			}
			h.mu.Unlock()
		case <-h.stopChan:
			return
		}
	}
}

// Stop shuts down the hub and all controllers.
func (h *Hub) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, c := range h.controllers {
		_ = c.Stop()
	}
	close(h.stopChan)
}
