package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Client WebSocket 客户端
type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID uint
	RoomID string
}

// Hub 管理所有 WebSocket 连接
type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	rooms      map[string]map[*Client]bool
	mu         sync.RWMutex
}

// Message WebSocket 消息结构
type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
	From uint            `json:"from,omitempty"`
	Room string          `json:"room,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境应检查 Origin
	},
}

// NewHub 创建 Hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		rooms:      make(map[string]map[*Client]bool),
	}
}

// Run 运行 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			log.Printf("[WS] Client connected. Total: %d", len(h.Clients))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				h.removeFromRoom(client)
			}
			h.mu.Unlock()
			log.Printf("[WS] Client disconnected. Total: %d", len(h.Clients))

		case message := <-h.Broadcast:
			h.mu.RLock()
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// JoinRoom 加入房间
func (h *Hub) JoinRoom(client *Client, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][client] = true
	client.RoomID = roomID

	log.Printf("[WS] Client joined room %s", roomID)
}

// LeaveRoom 离开房间
func (h *Hub) LeaveRoom(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.removeFromRoom(client)
}

// BroadcastToRoom 向房间广播
func (h *Hub) BroadcastToRoom(roomID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.rooms[roomID]; ok {
		for client := range clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.rooms[roomID], client)
			}
		}
	}
}

func (h *Hub) removeFromRoom(client *Client) {
	if roomClients, ok := h.rooms[client.RoomID]; ok {
		delete(roomClients, client)
		if len(roomClients) == 0 {
			delete(h.rooms, client.RoomID)
		}
	}
	client.RoomID = ""
}

// ReadPump 读取客户端消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("[WS] Invalid message from client %d: %v", c.UserID, err)
			continue
		}

		// 根据消息类型分发
		switch msg.Type {
		case "join_room":
			var data struct{ RoomID string `json:"room_id"` }
			json.Unmarshal(msg.Data, &data)
			c.Hub.JoinRoom(c, data.RoomID)

		case "game_action":
			if c.RoomID != "" {
				c.Hub.BroadcastToRoom(c.RoomID, message)
			}

		case "chat":
			c.Hub.Broadcast <- message

		default:
			log.Printf("[WS] Unknown message type: %s", msg.Type)
		}
	}
}

// WritePump 向客户端发送消息
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}

// ServeWS 处理 WebSocket 连接
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade error: %v", err)
		return
	}

	client := &Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
