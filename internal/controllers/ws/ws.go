package ws

import (
	"financing-aggregator/internal/exchange"
	"go.uber.org/zap"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketHandler interface {
	SubscribeToApplicationUpdates(c *gin.Context)
	BroadcastNewOffer(appID string, offer exchange.OfferResponse)
	CloseAll()
}

type webSocketHandler struct {
	mu          sync.RWMutex
	logger      *zap.Logger
	upgrader    websocket.Upgrader
	connections map[string][]*websocket.Conn
}

func NewWebSocketHandler(logger *zap.Logger) WebSocketHandler {
	return &webSocketHandler{
		logger: logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		connections: make(map[string][]*websocket.Conn),
	}
}

// SubscribeToApplicationUpdates
//
// @Summary		Upgrades the HTTP connection to a WebSocket
// @Description Upgrades the HTTP connection to a WebSocket and subscribes the client
// @Description to real-time application updates. The client must provide the application ID
// @Description as a URL parameter. The connection is kept open until the client disconnects or an error occurs.
// @Security 	BearerAuth
// @Tags		wss
// @Accept		json
// @Param 		id path string true "Application ID"
// @Failure		400 {object} exchange.ErrorResponse
// @Router 		/ws/applications/{id} [get]
func (h *webSocketHandler) SubscribeToApplicationUpdates(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer conn.Close()

	appID := c.Param("id")
	if appID == "" {
		conn.WriteJSON(exchange.NewErrorResponse("application id is required"))
		return
	}

	h.mu.Lock()
	h.connections[appID] = append(h.connections[appID], conn)
	h.mu.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	h.mu.Lock()
	conns := h.connections[appID]
	for i, c := range conns {
		if c == conn {
			h.connections[appID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(h.connections[appID]) == 0 {
		delete(h.connections, appID)
	}
	h.mu.Unlock()
}

func (h *webSocketHandler) BroadcastNewOffer(appID string, application exchange.OfferResponse) {
	h.mu.RLock()
	conns := h.connections[appID]
	h.mu.RUnlock()

	for _, conn := range conns {
		if err := conn.WriteJSON(application); err != nil {
			h.logger.Error("failed to write application update", zap.Error(err))
			continue
		}
	}
}

func (h *webSocketHandler) CloseAll() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, conns := range h.connections {
		for _, conn := range conns {
			_ = conn.Close()
		}
	}
	h.connections = make(map[string][]*websocket.Conn)
}
