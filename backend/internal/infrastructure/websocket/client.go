package websocket

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// Client represents a WebSocket client connection
type Client struct {
	ID       string
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	RaffleID string
	UserID   *string // Optional: for authenticated connections
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WebSocket Client %s] Unexpected close error: %v", c.ID, err)
			}
			break
		}
		// We don't process messages from clients, this is just to keep connection alive
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Send the message
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

			// Send any queued messages as separate WebSocket frames
			// (each message should be a complete JSON object)
			n := len(c.Send)
			for i := 0; i < n; i++ {
				queuedMsg := <-c.Send
				if err := c.Conn.WriteMessage(websocket.TextMessage, queuedMsg); err != nil {
					return
				}
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// NewClient creates a new WebSocket client
func NewClient(hub *Hub, conn *websocket.Conn, raffleID string, userID *string) *Client {
	return &Client{
		ID:       uuid.New().String(),
		Hub:      hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		RaffleID: raffleID,
		UserID:   userID,
	}
}
