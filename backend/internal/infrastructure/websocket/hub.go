package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// MessageType represents the type of WebSocket message
type MessageType string

const (
	MessageTypeNumberUpdate       MessageType = "number_update"
	MessageTypeReservationExpired MessageType = "reservation_expired"
	MessageTypeReservationCreated MessageType = "reservation_created"
	MessageTypeError              MessageType = "error"
)

// Message represents a WebSocket message
type Message struct {
	Type     MessageType            `json:"type"`
	RaffleID string                 `json:"raffle_id"`
	Data     map[string]interface{} `json:"data"`
}

// Hub maintains active WebSocket connections and broadcasts messages
type Hub struct {
	// Clients organized by raffle_id -> set of clients
	raffles map[string]map[*Client]bool
	mu      sync.RWMutex

	// Channels for hub operations (exported for external use)
	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		raffles:    make(map[string]map[*Client]bool),
		Broadcast:  make(chan *Message, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop (should be called in a goroutine)
func (h *Hub) Run() {
	log.Println("[WebSocket Hub] Starting...")
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.broadcastToRaffle(message)
		}
	}
}

// registerClient registers a new client to a raffle
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.raffles[client.RaffleID] == nil {
		h.raffles[client.RaffleID] = make(map[*Client]bool)
	}

	h.raffles[client.RaffleID][client] = true

	log.Printf("[WebSocket Hub] Client %s registered to raffle %s (total: %d)",
		client.ID, client.RaffleID, len(h.raffles[client.RaffleID]))
}

// unregisterClient removes a client from a raffle
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.raffles[client.RaffleID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.Send)

			// If no clients remain for this raffle, remove the raffle entry
			if len(clients) == 0 {
				delete(h.raffles, client.RaffleID)
			}

			log.Printf("[WebSocket Hub] Client %s unregistered from raffle %s (remaining: %d)",
				client.ID, client.RaffleID, len(clients))
		}
	}
}

// broadcastToRaffle sends a message to all clients connected to a specific raffle
func (h *Hub) broadcastToRaffle(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.raffles[message.RaffleID]
	if !ok || len(clients) == 0 {
		// No clients connected to this raffle
		return
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("[WebSocket Hub] Error marshaling message: %v", err)
		return
	}

	// Send to all clients in this raffle
	for client := range clients {
		select {
		case client.Send <- messageJSON:
			// Message sent successfully
		default:
			// Channel is full, close the client
			log.Printf("[WebSocket Hub] Client %s channel full, closing", client.ID)
			close(client.Send)
			delete(clients, client)
		}
	}

	log.Printf("[WebSocket Hub] Broadcast %s to raffle %s (%d clients)",
		message.Type, message.RaffleID, len(clients))
}

// BroadcastNumberUpdate notifies all clients about a number status change
func (h *Hub) BroadcastNumberUpdate(raffleID, numberID, status string, userID *string) {
	data := map[string]interface{}{
		"number_id": numberID,
		"status":    status,
	}

	if userID != nil {
		data["user_id"] = *userID
	}

	h.Broadcast <- &Message{
		Type:     MessageTypeNumberUpdate,
		RaffleID: raffleID,
		Data:     data,
	}
}

// BroadcastReservationExpired notifies all clients that a reservation expired
func (h *Hub) BroadcastReservationExpired(raffleID string, numberIDs []string) {
	h.Broadcast <- &Message{
		Type:     MessageTypeReservationExpired,
		RaffleID: raffleID,
		Data: map[string]interface{}{
			"number_ids": numberIDs,
		},
	}
}

// BroadcastReservationCreated notifies all clients about a new reservation
func (h *Hub) BroadcastReservationCreated(raffleID string, numberIDs []string, userID string) {
	h.Broadcast <- &Message{
		Type:     MessageTypeReservationCreated,
		RaffleID: raffleID,
		Data: map[string]interface{}{
			"number_ids": numberIDs,
			"user_id":    userID,
		},
	}
}

// GetConnectedClients returns the number of clients connected to a raffle
func (h *Hub) GetConnectedClients(raffleID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.raffles[raffleID]; ok {
		return len(clients)
	}
	return 0
}

// GetTotalClients returns the total number of connected clients across all raffles
func (h *Hub) GetTotalClients() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	total := 0
	for _, clients := range h.raffles {
		total += len(clients)
	}
	return total
}

// GetActiveRaffles returns the number of raffles with active connections
func (h *Hub) GetActiveRaffles() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.raffles)
}
