package websocket

import (
	"net/http"

	"github.com/Tabhi109/investwise/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for the portfolio simulator (can be restricted in production config)
		return true
	},
}

// HandleWS upgrades HTTP connection and registers client with WebSocket Hub.
func HandleWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error("Failed to upgrade connection to WebSocket", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
			return
		}

		client := &Client{
			Hub:  hub,
			Conn: conn,
			Send: make(chan []byte, 256), // Buffered outbound queue
		}

		// Register client with the Hub
		client.Hub.Register <- client

		// Start read and write pumps as concurrent routines
		go client.WritePump()
		go client.ReadPump()
	}
}
