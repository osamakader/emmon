package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"emmon/monitor"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// WebServer handles the web interface
type WebServer struct {
	port     string
	log      *logrus.Logger
	monitor  *monitor.SystemMonitor
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	mu       sync.RWMutex
}

// NewWebServer creates a new web server instance
func NewWebServer(port string, log *logrus.Logger, monitor *monitor.SystemMonitor) *WebServer {
	return &WebServer{
		port:    port,
		log:     log,
		monitor: monitor,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for embedded use
			},
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

// Start starts the web server
func (ws *WebServer) Start() error {
	// Serve static files
	http.HandleFunc("/", ws.handleIndex)
	http.HandleFunc("/ws", ws.handleWebSocket)
	http.HandleFunc("/api/stats", ws.handleStats)

	// Start WebSocket broadcast goroutine
	go ws.broadcastStats()

	ws.log.Infof("Starting web server on port %s", ws.port)
	return http.ListenAndServe(":"+ws.port, nil)
}

// handleIndex serves the main HTML page
func (ws *WebServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, GetHTML())
}

// handleWebSocket handles WebSocket connections
func (ws *WebServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.log.Errorf("WebSocket upgrade failed: %v", err)
		return
	}

	ws.mu.Lock()
	ws.clients[conn] = true
	ws.mu.Unlock()

	ws.log.Infof("New WebSocket client connected")

	// Handle client disconnect
	go func() {
		defer func() {
			conn.Close()
			ws.mu.Lock()
			delete(ws.clients, conn)
			ws.mu.Unlock()
			ws.log.Infof("WebSocket client disconnected")
		}()

		for {
			// Read messages (we don't need to handle them for now)
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}

// handleStats serves current system stats as JSON
func (ws *WebServer) handleStats(w http.ResponseWriter, r *http.Request) {
	stats, err := ws.monitor.GetSystemStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// broadcastStats broadcasts system stats to all connected WebSocket clients
func (ws *WebServer) broadcastStats() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats, err := ws.monitor.GetSystemStats()
		if err != nil {
			ws.log.Errorf("Failed to get system stats: %v", err)
			continue
		}

		ws.mu.RLock()
		for client := range ws.clients {
			err := client.WriteJSON(stats)
			if err != nil {
				ws.log.Errorf("Failed to send stats to client: %v", err)
				client.Close()
				delete(ws.clients, client)
			}
		}
		ws.mu.RUnlock()
	}
}
