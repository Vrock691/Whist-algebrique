package event

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// MessageHandler receives client-to-server messages after their envelope and
// channel have been validated.
type MessageHandler func(gameID uuid.UUID, envelope Envelope) error

type WebSocketServer struct {
	upgrader      websocket.Upgrader
	handleMessage MessageHandler

	mu    sync.RWMutex
	rooms map[uuid.UUID]map[*webSocketClient]struct{}
}

type webSocketClient struct {
	conn      *websocket.Conn
	gameID    uuid.UUID
	send      chan []byte
	done      chan struct{}
	closeOnce sync.Once
}

func NewWebSocketServer(handleMessage MessageHandler) *WebSocketServer {
	return &WebSocketServer{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
		},
		handleMessage: handleMessage,
		rooms:         make(map[uuid.UUID]map[*webSocketClient]struct{}),
	}
}

// Handle upgrades GET /ws?gameId=<uuid> to a WebSocket connection.
func (s *WebSocketServer) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameID, err := uuid.Parse(r.URL.Query().Get("gameId"))
	if err != nil {
		http.Error(w, "gameId must be a valid UUID", http.StatusBadRequest)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &webSocketClient{
		conn: conn, gameID: gameID, send: make(chan []byte, 16), done: make(chan struct{}),
	}
	s.addClient(client)
	defer s.removeClient(client)

	go s.writePump(client)
	s.readPump(client)
}

func (s *WebSocketServer) readPump(client *webSocketClient) {
	client.conn.SetReadLimit(64 * 1024)
	_ = client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		return client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	for {
		messageType, data, err := client.conn.ReadMessage()
		if err != nil {
			return
		}
		if messageType != websocket.TextMessage {
			s.sendError(client, "INVALID_MESSAGE", "WebSocket messages must be text")
			continue
		}

		var envelope Envelope
		if err := json.Unmarshal(data, &envelope); err != nil || envelope.Channel == "" || len(envelope.Payload) == 0 {
			s.sendError(client, "INVALID_ENVELOPE", "message must contain channel and payload")
			continue
		}
		if !isClientChannel(envelope.Channel) {
			s.sendError(client, "UNSUPPORTED_CHANNEL", "channel is not a client-to-server channel")
			continue
		}
		if s.handleMessage != nil {
			if err := s.handleMessage(client.gameID, envelope); err != nil {
				s.sendError(client, "MESSAGE_REJECTED", err.Error())
			}
		}

		ack := map[string]string{
			"channel": envelope.Channel,
			"status":  "received",
		}
		ackData, _ := json.Marshal(ack)
		client.send <- ackData


	}
}

func (s *WebSocketServer) writePump(client *webSocketClient) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-client.send:
			_ = client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				_ = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			_ = client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-client.done:
			return
		}
	}
}

func (s *WebSocketServer) Broadcast(gameID uuid.UUID, envelope Envelope) error {
	data, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	s.mu.RLock()
	clients := make([]*webSocketClient, 0, len(s.rooms[gameID]))
	for client := range s.rooms[gameID] {
		clients = append(clients, client)
	}
	s.mu.RUnlock()

	var fullClient *webSocketClient
	for _, client := range clients {
		select {
		case client.send <- data:
		default:
			fullClient = client
		}
	}
	if fullClient != nil {
		s.removeClient(fullClient)
		return errors.New("WebSocket client send buffer is full")
	}
	return nil
}

func (s *WebSocketServer) addClient(client *webSocketClient) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.rooms[client.gameID] == nil {
		s.rooms[client.gameID] = make(map[*webSocketClient]struct{})
	}
	s.rooms[client.gameID][client] = struct{}{}
}

func (s *WebSocketServer) removeClient(client *webSocketClient) {
	client.closeOnce.Do(func() {
		s.mu.Lock()
		if room := s.rooms[client.gameID]; room != nil {
			delete(room, client)
			if len(room) == 0 {
				delete(s.rooms, client.gameID)
			}
		}
		close(client.done)
		s.mu.Unlock()
		_ = client.conn.Close()
	})
}

func (s *WebSocketServer) sendError(client *webSocketClient, code, message string) {
	payload, err := json.Marshal(map[string]string{"code": code, "message": message})
	if err == nil {
		select {
		case client.send <- payload:
		default:
		}
	}
}

func isClientChannel(channel string) bool {
	return channel == "bet/placed" || channel == "play/card"
}
