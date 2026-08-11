package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

// Client WebSocket 客户端
type Client struct {
	Hub         *Hub
	Conn        *websocket.Conn
	Send        chan []byte
	UserID      uint
	Username    string
	RoomID      string
	mu          sync.RWMutex
}

// Hub 管理所有 WebSocket 连接
type Hub struct {
	Clients     map[*Client]bool
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan []byte
	rooms       map[string]map[*Client]bool
	RoomMgr     *RoomManager
	JWTSecret   string
	mu          sync.RWMutex
	clientByUID map[uint]*Client // userID -> client

	// 匹配队列
	matchQueueMu sync.Mutex
	matchQueue   []*Client
}

// Message WebSocket 消息结构
type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NewHub 创建 Hub
func NewHub(jwtSecret string) *Hub {
	return &Hub{
		Clients:     make(map[*Client]bool),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan []byte),
		rooms:       make(map[string]map[*Client]bool),
		RoomMgr:     NewRoomManager(),
		JWTSecret:   jwtSecret,
		clientByUID: make(map[uint]*Client),
	}
}

// Run 运行 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.clientByUID[client.UserID] = client
			h.mu.Unlock()
			log.Printf("[WS] Client connected. UserID=%d Total: %d", client.UserID, len(h.Clients))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				if client.UserID > 0 {
					delete(h.clientByUID, client.UserID)
				}
				close(client.Send)
				h.removeFromRoom(client)

				// 如果客户端有 UserID，处理离开房间逻辑
				if client.UserID > 0 {
					room, stillExists := h.RoomMgr.LeaveRoom(client.UserID)
					if room != nil {
						if stillExists {
							h.broadcastRoomState(room)
						} else {
							// 通知客座玩家房主离开
							if guest := h.getClientByUID(room.GuestID); guest != nil {
								guest.Send <- mustMarshal(map[string]interface{}{
									"type":    "room_closed",
									"message": "房主已离开",
								})
							}
						}
					}
				}
			}
			h.mu.Unlock()
			// 广播房间列表更新
			h.broadcastRoomList()
			log.Printf("[WS] Client disconnected. UserID=%d Total: %d", client.UserID, len(h.Clients))

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

// OnlineCount 在线人数
func (h *Hub) OnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.Clients)
}

// IsUserOnline 判断用户是否在线
func (h *Hub) IsUserOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clientByUID[userID]
	return ok
}

// getClientByUID 通过 userID 获取客户端
func (h *Hub) getClientByUID(uid uint) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clientByUID[uid]
}

// JoinWSRoom 将客户端加入 WS 房间（内部使用）
func (h *Hub) JoinWSRoom(client *Client, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][client] = true
	client.RoomID = roomID
}

// LeaveWSRoom 将客户端移出 WS 房间
func (h *Hub) LeaveWSRoom(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.removeFromRoom(client)
}

// BroadcastToWSRoom 向 WS 房间内所有客户端广播
func (h *Hub) BroadcastToWSRoom(roomID string, message []byte) {
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

// ========= 匹配队列 =========

// JoinMatchQueue 加入匹配队列
func (h *Hub) JoinMatchQueue(c *Client) bool {
	h.matchQueueMu.Lock()
	defer h.matchQueueMu.Unlock()

	// 检查是否已在队列
	for _, client := range h.matchQueue {
		if client == c {
			return false
		}
	}

	// 寻找对手
	if len(h.matchQueue) > 0 {
		opponent := h.matchQueue[0]
		h.matchQueue = h.matchQueue[1:]

		// 创建房间并通知双方
		h.performMatchmaking(c, opponent)
		return true
	}

	h.matchQueue = append(h.matchQueue, c)
	c.Send <- mustMarshal(map[string]interface{}{
		"type":    "match_queued",
		"message": "正在寻找对手...",
	})
	return false
}

// RemoveFromMatchQueue 从匹配队列移除
func (h *Hub) RemoveFromMatchQueue(c *Client) {
	h.matchQueueMu.Lock()
	defer h.matchQueueMu.Unlock()
	for i, client := range h.matchQueue {
		if client == c {
			h.matchQueue = append(h.matchQueue[:i], h.matchQueue[i+1:]...)
			break
		}
	}
}

func (h *Hub) performMatchmaking(p1, p2 *Client) {
	roomID := fmt.Sprintf("match_%d_%d", p1.UserID, p2.UserID)
	host := p1
	guest := p2
	room := h.RoomMgr.CreateRoom(roomID, host.UserID)
	if room == nil {
		p1.Send <- mustMarshal(map[string]interface{}{"type": "error", "message": "创建匹配房间失败"})
		p2.Send <- mustMarshal(map[string]interface{}{"type": "error", "message": "创建匹配房间失败"})
		return
	}
	h.RoomMgr.JoinRoom(roomID, guest.UserID)
	h.JoinWSRoom(host, roomID)
	h.JoinWSRoom(guest, roomID)

	// 自动开始
	h.RoomMgr.SetReady(roomID, host.UserID, true)
	h.RoomMgr.SetReady(roomID, guest.UserID, true)
	h.RoomMgr.StartGame(roomID)

	host.Send <- mustMarshal(map[string]interface{}{
		"type":             "match_found",
		"room_id":          roomID,
		"is_host":          true,
		"opponent_id":      guest.UserID,
		"opponent_name":    guest.Username,
	})
	guest.Send <- mustMarshal(map[string]interface{}{
		"type":             "match_found",
		"room_id":          roomID,
		"is_host":          false,
		"opponent_id":      host.UserID,
		"opponent_name":    host.Username,
	})

	// 通知游戏开始
	host.Send <- mustMarshal(map[string]interface{}{
		"type":             "game_start",
		"room_id":          roomID,
		"opponent_name":    guest.Username,
	})
	guest.Send <- mustMarshal(map[string]interface{}{
		"type":             "game_start",
		"room_id":          roomID,
		"opponent_name":    host.Username,
	})

	log.Printf("[Match] %s(%d) vs %s(%d) - Room: %s", host.Username, host.UserID, guest.Username, guest.UserID, roomID)

	_ = room // suppress unused warning
}

// handleGetRoomList 返回当前所有房间列表
func (c *Client) handleGetRoomList() {
	rooms := c.Hub.RoomMgr.ListRooms()
	c.Send <- mustMarshal(map[string]interface{}{
		"type": "room_list",
		"data": rooms,
	})
}

// broadcastRoomList 向所有已连接客户端广播房间列表
func (h *Hub) broadcastRoomList() {
	rooms := h.RoomMgr.ListRooms()
	msg := mustMarshal(map[string]interface{}{
		"type": "room_list_update",
		"data": rooms,
	})
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.Clients {
		select {
		case client.Send <- msg:
		default:
		}
	}
}

// broadcastRoomState 广播房间状态更新
func (h *Hub) broadcastRoomState(room *GameRoom) {
	if room == nil {
		return
	}
	room.mu.RLock()
	defer room.mu.RUnlock()

	state := map[string]interface{}{
		"type":           "room_update",
		"room_id":        room.ID,
		"host_id":        room.HostID,
		"guest_id":       room.GuestID,
		"host_char_id":   room.HostCharID,
		"guest_char_id":  room.GuestCharID,
		"host_skin_id":   room.HostSkinID,
		"guest_skin_id":  room.GuestSkinID,
		"host_ready":     room.HostReady,
		"guest_ready":    room.GuestReady,
		"status":         string(room.Status),
	}

	data, _ := json.Marshal(state)
	h.BroadcastToWSRoom(room.ID, data)
}

// ===================== 消息处理 =====================

// ReadPump 读取客户端消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("[WS] Read error from user %d: %v", c.UserID, err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("[WS] Invalid message from user %d: %v", c.UserID, err)
			continue
		}

		c.handleMessage(msg.Type, msg.Data, message)
	}
}

func (c *Client) handleMessage(msgType string, data json.RawMessage, raw []byte) {
	switch msgType {
	case "create_room":
		c.handleCreateRoom(data)
	case "join_room":
		c.handleJoinRoom(data)
	case "leave_room":
		c.handleLeaveRoom()
	case "set_ready":
		c.handleSetReady(data)
	case "select_character":
		c.handleSelectCharacter(data)
	case "select_skin":
		c.handleSelectSkin(data)
	case "start_game":
		c.handleStartGame()
	case "frame_input":
		c.handleFrameInput(raw)
	case "battle_over":
		c.handleBattleOver(data)
	case "queue_match":
		c.handleQueueMatch()
	case "cancel_queue":
		c.handleCancelQueue()
	case "invite_friend":
		c.handleFriendInvite(data)
	case "invite_response":
		c.handleInviteResponse(data)
	case "get_room_list":
		c.handleGetRoomList()
	default:
		log.Printf("[WS] Unknown message type: %s from user %d", msgType, c.UserID)
	}
}

// ===================== 消息处理器 =====================

type createRoomData struct {
	RoomID      string `json:"room_id"`
	CharacterID uint   `json:"character_id"`
	SkinID      uint   `json:"skin_id"`
}

func (c *Client) handleCreateRoom(data json.RawMessage) {
	if c.UserID == 0 {
		c.sendError("请先登录")
		return
	}

	var d createRoomData
	json.Unmarshal(data, &d)

	roomID := d.RoomID
	if roomID == "" {
		roomID = fmt.Sprintf("room_%d_%d", c.UserID, time.Now().UnixMilli()%100000)
	}

	existingRoom := c.Hub.RoomMgr.GetUserRoom(c.UserID)
	if existingRoom != nil {
		c.Hub.RoomMgr.LeaveRoom(c.UserID)
		if existingRoom.Status != StatusFinished {
			c.Hub.LeaveWSRoom(c)
		}
	}

	c.Hub.RoomMgr.CreateRoom(roomID, c.UserID)
	c.Hub.JoinWSRoom(c, roomID)

	if d.CharacterID > 0 {
		c.Hub.RoomMgr.SelectCharacter(roomID, c.UserID, d.CharacterID)
	}
	if d.SkinID > 0 {
		c.Hub.RoomMgr.SelectSkin(roomID, c.UserID, d.SkinID)
	}

	c.Send <- mustMarshal(map[string]interface{}{
		"type":    "room_created",
		"room_id": roomID,
	})

	// 广播房间列表给所有客户端
	c.Hub.broadcastRoomList()
}

func (c *Client) handleJoinRoom(data json.RawMessage) {
	if c.UserID == 0 {
		c.sendError("请先登录")
		return
	}

	var d struct{ RoomID string `json:"room_id"` }
	json.Unmarshal(data, &d)

	if d.RoomID == "" {
		c.sendError("缺少房间ID")
		return
	}

	existingRoom := c.Hub.RoomMgr.GetUserRoom(c.UserID)
	if existingRoom != nil {
		c.Hub.RoomMgr.LeaveRoom(c.UserID)
		c.Hub.LeaveWSRoom(c)
	}

	room, err := c.Hub.RoomMgr.JoinRoom(d.RoomID, c.UserID)
	if err != nil {
		c.Send <- err.(*wsError).ToJSON()
		return
	}

	c.Hub.JoinWSRoom(c, d.RoomID)

	// 通知双方
	hostClient := c.Hub.getClientByUID(room.HostID)
	if hostClient != nil {
		hostClient.Send <- mustMarshal(map[string]interface{}{
			"type":     "player_joined",
			"room_id":  d.RoomID,
			"guest_id": c.UserID,
		})
	}

	c.Send <- mustMarshal(map[string]interface{}{
		"type":    "room_joined",
		"room_id": d.RoomID,
		"host_id": room.HostID,
	})

	// 广播房间列表给所有客户端
	c.Hub.broadcastRoomList()
}

func (c *Client) handleLeaveRoom() {
	room, stillExists := c.Hub.RoomMgr.LeaveRoom(c.UserID)
	c.Hub.LeaveWSRoom(c)

	if room != nil {
		if stillExists {
			c.Hub.broadcastRoomState(room)
		} else {
			// 广播房间关闭
			c.Hub.BroadcastToWSRoom(room.ID, mustMarshal(map[string]interface{}{
				"type":    "room_closed",
				"message": "房主已离开房间",
			}))
		}
	}

	c.Send <- mustMarshal(map[string]interface{}{
		"type":    "left_room",
		"message": "已离开房间",
	})

	// 广播房间列表给所有客户端
	c.Hub.broadcastRoomList()
}

func (c *Client) handleSetReady(data json.RawMessage) {
	room := c.Hub.RoomMgr.GetUserRoom(c.UserID)
	if room == nil {
		c.sendError("不在任何房间中")
		return
	}

	var d struct{ Ready bool `json:"ready"` }
	json.Unmarshal(data, &d)

	updated := c.Hub.RoomMgr.SetReady(room.ID, c.UserID, d.Ready)
	if updated != nil {
		c.Hub.broadcastRoomState(updated)
	}
}

func (c *Client) handleSelectCharacter(data json.RawMessage) {
	room := c.Hub.RoomMgr.GetUserRoom(c.UserID)
	if room == nil {
		c.sendError("不在任何房间中")
		return
	}

	var d struct{ CharacterID uint `json:"character_id"` }
	json.Unmarshal(data, &d)

	if d.CharacterID == 0 {
		c.sendError("请选择角色")
		return
	}

	updated := c.Hub.RoomMgr.SelectCharacter(room.ID, c.UserID, d.CharacterID)
	if updated != nil {
		c.Hub.broadcastRoomState(updated)
	}
}

func (c *Client) handleSelectSkin(data json.RawMessage) {
	room := c.Hub.RoomMgr.GetUserRoom(c.UserID)
	if room == nil {
		c.sendError("不在任何房间中")
		return
	}

	var d struct{ SkinID uint `json:"skin_id"` }
	json.Unmarshal(data, &d)

	if d.SkinID == 0 {
		c.sendError("请选择皮肤")
		return
	}

	updated := c.Hub.RoomMgr.SelectSkin(room.ID, c.UserID, d.SkinID)
	if updated != nil {
		c.Hub.broadcastRoomState(updated)
	}
}

func (c *Client) handleStartGame() {
	room := c.Hub.RoomMgr.GetUserRoom(c.UserID)
	if room == nil || room.HostID != c.UserID {
		c.sendError("只有房主可以开始游戏")
		return
	}

	if room.GuestID == 0 {
		c.sendError("等待玩家加入")
		return
	}

	// 3秒倒计时
	for i := 3; i > 0; i-- {
		c.Hub.BroadcastToWSRoom(room.ID, mustMarshal(map[string]interface{}{
			"type":      "game_countdown",
			"countdown": i,
		}))
		time.Sleep(1 * time.Second)
	}

	updated := c.Hub.RoomMgr.StartGame(room.ID)

	// 通知双方游戏开始
	c.Hub.BroadcastToWSRoom(room.ID, mustMarshal(map[string]interface{}{
		"type":           "game_start",
		"room_id":        room.ID,
		"host_id":        room.HostID,
		"guest_id":       room.GuestID,
		"host_char_id":   updated.HostCharID,
		"guest_char_id":  updated.GuestCharID,
		"host_skin_id":   updated.HostSkinID,
		"guest_skin_id":  updated.GuestSkinID,
	}))
}

func (c *Client) handleFrameInput(raw []byte) {
	room := c.Hub.RoomMgr.GetUserRoom(c.UserID)
	if room == nil {
		return
	}

	room.mu.Lock()
	room.FrameCount++
	frame := room.FrameCount
	room.mu.Unlock()

	// 将帧号和消息打包广播给房间其他人
	var rawMsg map[string]interface{}
	json.Unmarshal(raw, &rawMsg)
	rawMsg["frame"] = frame
	rawMsg["from"] = c.UserID

	frameData, _ := json.Marshal(rawMsg)
	c.Hub.BroadcastToWSRoom(room.ID, frameData)
}

type battleOverData struct {
	WinnerID  uint   `json:"winner_id"`
	Result    string `json:"result"`
}

func (c *Client) handleBattleOver(data json.RawMessage) {
	room := c.Hub.RoomMgr.GetUserRoom(c.UserID)
	if room == nil {
		return
	}

	c.Hub.RoomMgr.FinishGame(room.ID)

	var d battleOverData
	json.Unmarshal(data, &d)

	c.Hub.BroadcastToWSRoom(room.ID, mustMarshal(map[string]interface{}{
		"type":      "battle_result",
		"winner_id": d.WinnerID,
		"result":    d.Result,
	}))

	log.Printf("[Game] Room %s battle over. Winner: %d Result: %s", room.ID, d.WinnerID, d.Result)
}

// ========= 匹配 & 邀请 =========

func (c *Client) handleQueueMatch() {
	if c.UserID == 0 {
		c.sendError("请先登录")
		return
	}
	matched := c.Hub.JoinMatchQueue(c)
	if matched {
		log.Printf("[Match] User %d matched", c.UserID)
	}
}

func (c *Client) handleCancelQueue() {
	c.Hub.RemoveFromMatchQueue(c)
	c.Send <- mustMarshal(map[string]interface{}{
		"type":    "match_cancelled",
		"message": "已取消匹配",
	})
}

type inviteData struct {
	FriendID uint `json:"friend_id"`
}

func (c *Client) handleFriendInvite(data json.RawMessage) {
	if c.UserID == 0 {
		c.sendError("请先登录")
		return
	}

	var d inviteData
	json.Unmarshal(data, &d)
	if d.FriendID == 0 {
		c.sendError("请指定好友ID")
		return
	}

	target := c.Hub.getClientByUID(d.FriendID)
	if target == nil {
		c.Send <- mustMarshal(map[string]interface{}{
			"type":    "error",
			"message": "好友不在线",
		})
		return
	}

	// 发送邀请
	target.Send <- mustMarshal(map[string]interface{}{
		"type":        "friend_invite",
		"from_id":     c.UserID,
		"from_name":   c.Username,
	})
	c.Send <- mustMarshal(map[string]interface{}{
		"type":    "invite_sent",
		"message": "邀请已发送",
	})
}

type inviteResponseData struct {
	FromID  uint `json:"from_id"`
	Accept  bool `json:"accept"`
}

func (c *Client) handleInviteResponse(data json.RawMessage) {
	if c.UserID == 0 {
		c.sendError("请先登录")
		return
	}

	var d inviteResponseData
	json.Unmarshal(data, &d)

	inviter := c.Hub.getClientByUID(d.FromID)
	if inviter == nil {
		c.Send <- mustMarshal(map[string]interface{}{
			"type":    "error",
			"message": "对方已离线",
		})
		return
	}

	if !d.Accept {
		inviter.Send <- mustMarshal(map[string]interface{}{
			"type":      "invite_result",
			"accepted":  false,
			"from_name": c.Username,
		})
		return
	}

	// 接受：创建房间
	roomID := fmt.Sprintf("room_%d_%d", inviter.UserID, time.Now().UnixMilli()%100000)
	room := c.Hub.RoomMgr.CreateRoom(roomID, inviter.UserID)
	if room == nil {
		c.sendError("创建房间失败")
		return
	}
	c.Hub.RoomMgr.JoinRoom(roomID, c.UserID)
	c.Hub.JoinWSRoom(inviter, roomID)
	c.Hub.JoinWSRoom(c, roomID)

	// 双方 ready 并自动开始
	c.Hub.RoomMgr.SetReady(roomID, inviter.UserID, true)
	c.Hub.RoomMgr.SetReady(roomID, c.UserID, true)
	c.Hub.RoomMgr.StartGame(roomID)

	inviter.Send <- mustMarshal(map[string]interface{}{
		"type":          "invite_result",
		"accepted":      true,
		"room_id":       roomID,
		"opponent_name": c.Username,
		"opponent_id":   c.UserID,
	})
	c.Send <- mustMarshal(map[string]interface{}{
		"type":          "game_start",
		"room_id":       roomID,
		"opponent_name": inviter.Username,
		"opponent_id":   inviter.UserID,
	})

	_ = room
	log.Printf("[Invite] %s(%d) vs %s(%d) - Room: %s", inviter.Username, inviter.UserID, c.Username, c.UserID, roomID)
}

// ===================== 工具函数 =====================

func (c *Client) sendError(msg string) {
	c.Send <- mustMarshal(map[string]interface{}{
		"type":    "error",
		"message": msg,
	})
}

func parseJWT(tokenStr string, secret string) (uint, string) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, ""
	}

	uid := uint(claims["user_id"].(float64))
	username := claims["username"].(string)
	return uid, username
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
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
	tokenStr := r.URL.Query().Get("token")

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

	// JWT 认证（可选，允许未登录浏览但限制操作）
	if tokenStr != "" {
		uid, username := parseJWT(tokenStr, hub.JWTSecret)
		if uid == 0 {
			log.Printf("[WS] Invalid JWT")
			client.Send <- mustMarshal(map[string]string{
				"type":    "error",
				"message": "登录已过期，请重新登录",
			})
			go client.WritePump() // 发送错误消息后断开
			time.Sleep(100 * time.Millisecond)
			conn.Close()
			return
		}
		client.UserID = uid
		client.Username = username
	}

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
