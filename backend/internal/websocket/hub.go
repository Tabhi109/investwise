package websocket

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/Tabhi109/investwise/internal/logger"
)

// SubscriptionRequest carries details about subscribing/unsubscribing clients
type SubscriptionRequest struct {
	Client *Client
	Ticker string
}

// BroadcastMessage defines a price tick update
type BroadcastMessage struct {
	Ticker string
	Data   []byte
}

// Hub maintains the set of active clients and handles routing message ticks to subscribers
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Subscriptions map: Ticker -> set of Clients
	subscriptions map[string]map[*Client]bool

	// Inbound messages from the simulation feed.
	Broadcast chan *BroadcastMessage

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	// Subscribe requests
	Subscribe chan *SubscriptionRequest

	// Unsubscribe requests
	Unsubscribe chan *SubscriptionRequest

	// Context and cancel for shutdown
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewHub creates and initializes a WebSocket Hub
func NewHub() *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Hub{
		Broadcast:     make(chan *BroadcastMessage, 1024),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Subscribe:     make(chan *SubscriptionRequest),
		Unsubscribe:   make(chan *SubscriptionRequest),
		clients:       make(map[*Client]bool),
		subscriptions: make(map[string]map[*Client]bool),
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Run starts the central event loop of the Hub.
// Runs inside its own goroutine managed by composition root.
func (h *Hub) Run() {
	h.wg.Add(1)
	defer h.wg.Done()

	logger.Info("WebSocket Hub is running")
	for {
		select {
		case <-h.ctx.Done():
			logger.Info("WebSocket Hub stopping, cleaning up active connections")
			for client := range h.clients {
				close(client.Send)
				delete(h.clients, client)
			}
			return

		case client := <-h.Register:
			h.clients[client] = true
			logger.Debug("New WebSocket client connected")

		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				h.removeClientFromAllSubscriptions(client)
				logger.Debug("WebSocket client disconnected")
			}

		case req := <-h.Subscribe:
			if _, ok := h.clients[req.Client]; ok {
				if h.subscriptions[req.Ticker] == nil {
					h.subscriptions[req.Ticker] = make(map[*Client]bool)
				}
				h.subscriptions[req.Ticker][req.Client] = true
				logger.Debug("Client subscribed to ticker", "ticker", req.Ticker)
				
				// Send confirmation message to client
				confirmMsg, _ := json.Marshal(map[string]string{
					"type":   "info",
					"status": "subscribed",
					"ticker": req.Ticker,
				})
				select {
				case req.Client.Send <- confirmMsg:
				default:
				}
			}

		case req := <-h.Unsubscribe:
			if clients, ok := h.subscriptions[req.Ticker]; ok {
				delete(clients, req.Client)
				if len(clients) == 0 {
					delete(h.subscriptions, req.Ticker)
				}
				logger.Debug("Client unsubscribed from ticker", "ticker", req.Ticker)
				
				// Send confirmation message to client
				confirmMsg, _ := json.Marshal(map[string]string{
					"type":   "info",
					"status": "unsubscribed",
					"ticker": req.Ticker,
				})
				select {
				case req.Client.Send <- confirmMsg:
				default:
				}
			}

		case msg := <-h.Broadcast:
			// Route price updates to subscribers of specific ticker
			if clients, ok := h.subscriptions[msg.Ticker]; ok {
				for client := range clients {
					select {
					case client.Send <- msg.Data:
					default:
						// Buffer is full (backpressure): drop slow client
						logger.Warn("Client send queue full (backpressure), unregistering client")
						h.Unregister <- client
					}
				}
			}
		}
	}
}

// Close signals the hub to close and waits for active routines to finish
func (h *Hub) Close() {
	h.cancel()
	h.wg.Wait()
}

func (h *Hub) removeClientFromAllSubscriptions(client *Client) {
	for ticker, clients := range h.subscriptions {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			if len(clients) == 0 {
				delete(h.subscriptions, ticker)
			}
		}
	}
}
