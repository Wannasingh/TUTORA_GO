package ws

import (
	"encoding/json"
	"log"
	"sync"
)

// WSMessage represents a WebSocket message envelope
type WSMessage struct {
	Type    string      `json:"type"`    // "message", "notification", "typing"
	Payload interface{} `json:"payload"`
}

// Hub maintains active WebSocket clients and broadcasts messages
type Hub struct {
	clients    map[int]map[*Client]bool // userID -> set of clients
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main event loop — call this as a goroutine
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()
			log.Printf("[WS Hub] User %d connected (total connections: %d)", client.UserID, len(h.clients[client.UserID]))

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.UserID]; ok {
				delete(clients, client)
				close(client.Send)
				if len(clients) == 0 {
					delete(h.clients, client.UserID)
				}
			}
			h.mu.Unlock()
			log.Printf("[WS Hub] User %d disconnected", client.UserID)
		}
	}
}

// SendToUser sends a message to all active connections for a specific user
func (h *Hub) SendToUser(userID int, msg *WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[WS Hub] Failed to marshal message: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[userID]; ok {
		for client := range clients {
			select {
			case client.Send <- data:
			default:
				// Client send buffer is full, skip
			}
		}
	}
}

// IsOnline checks if a user has any active WebSocket connections
func (h *Hub) IsOnline(userID int) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[userID]) > 0
}
