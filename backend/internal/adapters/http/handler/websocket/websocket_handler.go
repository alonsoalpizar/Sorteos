package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	ws "github.com/sorteos-platform/backend/internal/infrastructure/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: In production, validate origin against allowed domains
		// origin := r.Header.Get("Origin")
		// return origin == "https://sorteos.com" || origin == "https://www.sorteos.com"
		return true // Allow all origins in development
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub *ws.Hub
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleConnection upgrades HTTP connection to WebSocket for a specific raffle
// Route: GET /api/v1/raffles/:id/ws
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	raffleID := c.Param("id")
	if raffleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "raffle_id is required"})
		return
	}

	// Optional: Get user ID from JWT token if authenticated
	var userID *string
	if userIDVal, exists := c.Get("user_id"); exists {
		if uid, ok := userIDVal.(string); ok {
			userID = &uid
		}
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WebSocket Handler] Failed to upgrade connection: %v", err)
		return
	}

	// Create new client
	client := ws.NewClient(h.hub, conn, raffleID, userID)

	// Register client with hub
	h.hub.Register <- client

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()

	log.Printf("[WebSocket Handler] Client %s connected to raffle %s", client.ID, raffleID)
}

// GetConnectionStats returns statistics about WebSocket connections
// Route: GET /api/v1/raffles/:id/ws/stats
func (h *WebSocketHandler) GetConnectionStats(c *gin.Context) {
	raffleID := c.Param("id")

	stats := gin.H{
		"raffle_id":         raffleID,
		"connected_clients": h.hub.GetConnectedClients(raffleID),
	}

	c.JSON(http.StatusOK, stats)
}

// GetGlobalStats returns global WebSocket statistics
// Route: GET /api/v1/admin/websocket/stats
func (h *WebSocketHandler) GetGlobalStats(c *gin.Context) {
	stats := gin.H{
		"total_clients":  h.hub.GetTotalClients(),
		"active_raffles": h.hub.GetActiveRaffles(),
	}

	c.JSON(http.StatusOK, stats)
}
