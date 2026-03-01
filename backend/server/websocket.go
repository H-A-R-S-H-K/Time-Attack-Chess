package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// RegisterWebSocket registers the WebSocket handler.
func RegisterWebSocket(mux *http.ServeMux, gm *GameManager) {
	mux.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, gm)
	})
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, gm *GameManager) {
	// Extract game ID from URL: /ws/{gameId}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "missing game ID", http.StatusBadRequest)
		return
	}
	gameID := parts[len(parts)-1]

	// Get color from query param
	colorParam := r.URL.Query().Get("color")

	session, ok := gm.GetGame(gameID)
	if !ok {
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Determine player color
	var playerColor PlayerColor
	switch colorParam {
	case "white":
		playerColor = ColorWhite
	case "black":
		playerColor = ColorBlack
	default:
		playerColor = ColorSpectator
	}

	pc := &PlayerConn{
		Conn:  conn,
		Color: playerColor,
	}

	// Register connection
	session.mu.Lock()
	switch playerColor {
	case ColorWhite:
		session.Players[0] = pc
	case ColorBlack:
		session.Players[1] = pc
		// Both players connected — start the game timer
		if session.Players[0] != nil {
			session.LastMoveAt = time.Now()
			go session.StartTimerTicker()
		}
	default:
		session.Spectators = append(session.Spectators, pc)
	}
	session.mu.Unlock()

	log.Printf("WebSocket connected: game=%s color=%s", gameID, colorParam)

	// Send initial game state
	session.mu.RLock()
	initialState := BuildGameStateResponse(session.GameState, -1, -1)
	session.mu.RUnlock()
	if err := pc.SendJSON(initialState); err != nil {
		log.Printf("Error sending initial state: %v", err)
	}

	// Read messages from client
	defer func() {
		conn.Close()
		session.mu.Lock()
		switch playerColor {
		case ColorWhite:
			session.Players[0] = nil
		case ColorBlack:
			session.Players[1] = nil
		}
		session.mu.Unlock()
		log.Printf("WebSocket disconnected: game=%s color=%s", gameID, colorParam)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg ClientMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			_ = pc.SendJSON(ErrorResponse{
				Type:    "ERROR",
				Message: "invalid message format",
			})
			continue
		}

		switch msg.Type {
		case "MOVE":
			resp, err := session.HandleMove(msg.From, msg.To, msg.Promotion, playerColor)
			if err != nil {
				_ = pc.SendJSON(ErrorResponse{
					Type:    "ERROR",
					Message: err.Error(),
				})
				continue
			}

			// Broadcast the updated game state to all
			session.BroadcastToAll(resp)

			// Check if game is over
			if resp.IsGameOver {
				gameOver := GameOverResponse{
					Type:   "GAME_OVER",
					Result: resp.GameResult,
					Reason: resp.GameOverReason,
				}
				session.BroadcastToAll(gameOver)
				session.StopTimerTicker()
			}

		case "LEAVE":
			return

		default:
			_ = pc.SendJSON(ErrorResponse{
				Type:    "ERROR",
				Message: "unknown message type: " + msg.Type,
			})
		}
	}
}
