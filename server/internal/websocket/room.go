package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// RoomStatus 房间状态
type RoomStatus string

const (
	StatusWaiting   RoomStatus = "waiting"
	StatusSelecting RoomStatus = "selecting"
	StatusReady     RoomStatus = "ready"
	StatusPlaying   RoomStatus = "playing"
	StatusFinished  RoomStatus = "finished"
)

// GameRoom 游戏房间
type GameRoom struct {
	ID          string     `json:"id"`
	HostID      uint       `json:"host_id"`
	GuestID     uint       `json:"guest_id"`
	HostCharID  uint       `json:"host_char_id"`
	GuestCharID uint       `json:"guest_char_id"`
	HostSkinID  uint       `json:"host_skin_id"`
	GuestSkinID uint       `json:"guest_skin_id"`
	HostReady   bool       `json:"host_ready"`
	GuestReady  bool       `json:"guest_ready"`
	Status      RoomStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	FrameCount  int        `json:"frame_count"`
	mu          sync.RWMutex
}

// RoomManager 房间管理器
type RoomManager struct {
	rooms    map[string]*GameRoom
	userRoom map[uint]string // userID -> roomID
	mu       sync.RWMutex
}

// NewRoomManager 创建房间管理器
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms:    make(map[string]*GameRoom),
		userRoom: make(map[uint]string),
	}
}

// CreateRoom 创建房间
func (rm *RoomManager) CreateRoom(roomID string, hostID uint) *GameRoom {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room := &GameRoom{
		ID:        roomID,
		HostID:    hostID,
		Status:    StatusWaiting,
		CreatedAt: time.Now(),
	}
	rm.rooms[roomID] = room
	rm.userRoom[hostID] = roomID

	log.Printf("[Room] Created room %s by user %d", roomID, hostID)
	return room
}

// JoinRoom 加入房间
func (rm *RoomManager) JoinRoom(roomID string, guestID uint) (*GameRoom, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, ok := rm.rooms[roomID]
	if !ok {
		return nil, ErrRoomNotFound
	}
	room.mu.Lock()
	defer room.mu.Unlock()

	if room.Status != StatusWaiting {
		return nil, ErrRoomNotAvailable
	}
	if room.GuestID != 0 {
		return nil, ErrRoomFull
	}
	if room.HostID == guestID {
		return nil, ErrSamePlayer
	}

	room.GuestID = guestID
	rm.userRoom[guestID] = roomID
	log.Printf("[Room] User %d joined room %s", guestID, roomID)
	return room, nil
}

// LeaveRoom 离开房间
func (rm *RoomManager) LeaveRoom(userID uint) (*GameRoom, bool) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	roomID, ok := rm.userRoom[userID]
	if !ok {
		return nil, false
	}

	room, exists := rm.rooms[roomID]
	if !exists {
		delete(rm.userRoom, userID)
		return nil, false
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	delete(rm.userRoom, userID)

	isHost := room.HostID == userID

	if isHost {
		// 房主离开，关闭房间
		delete(rm.rooms, roomID)
		if room.GuestID != 0 {
			delete(rm.userRoom, room.GuestID)
		}
		room.Status = StatusFinished
		room.GuestID = 0
	} else {
		// 客座玩家离开
		room.GuestID = 0
		room.GuestReady = false
		room.Status = StatusWaiting
	}

	log.Printf("[Room] User %d left room %s (isHost=%v)", userID, roomID, isHost)
	return room, !isHost // true = room still exists
}

// SetReady 设置准备状态
func (rm *RoomManager) SetReady(roomID string, userID uint, ready bool) *GameRoom {
	rm.mu.RLock()
	room := rm.rooms[roomID]
	rm.mu.RUnlock()
	if room == nil {
		return nil
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	if userID == room.HostID {
		room.HostReady = ready
	} else if userID == room.GuestID {
		room.GuestReady = ready
	}

	if room.HostReady && room.GuestReady && room.HostCharID > 0 && room.GuestCharID > 0 {
		room.Status = StatusReady
	}

	return room
}

// SelectCharacter 选择角色
func (rm *RoomManager) SelectCharacter(roomID string, userID uint, charID uint) *GameRoom {
	rm.mu.RLock()
	room := rm.rooms[roomID]
	rm.mu.RUnlock()
	if room == nil {
		return nil
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	if userID == room.HostID {
		room.HostCharID = charID
	} else if userID == room.GuestID {
		room.GuestCharID = charID
	}

	return room
}

// SelectSkin 选择皮肤
func (rm *RoomManager) SelectSkin(roomID string, userID uint, skinID uint) *GameRoom {
	rm.mu.RLock()
	room := rm.rooms[roomID]
	rm.mu.RUnlock()
	if room == nil {
		return nil
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	if userID == room.HostID {
		room.HostSkinID = skinID
	} else if userID == room.GuestID {
		room.GuestSkinID = skinID
	}

	return room
}

// StartGame 开始游戏
func (rm *RoomManager) StartGame(roomID string) *GameRoom {
	rm.mu.RLock()
	room := rm.rooms[roomID]
	rm.mu.RUnlock()
	if room == nil {
		return nil
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	room.Status = StatusPlaying
	room.FrameCount = 0
	return room
}

// FinishGame 结束游戏
func (rm *RoomManager) FinishGame(roomID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, ok := rm.rooms[roomID]
	if !ok {
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	room.Status = StatusFinished
}

// GetUserRoom 获取用户所在房间
func (rm *RoomManager) GetUserRoom(userID uint) *GameRoom {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	roomID, ok := rm.userRoom[userID]
	if !ok {
		return nil
	}
	return rm.rooms[roomID]
}

// GetRoom 获取房间
func (rm *RoomManager) GetRoom(roomID string) *GameRoom {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.rooms[roomID]
}

// ListRooms 列出所有等待中的房间
func (rm *RoomManager) ListRooms() []map[string]interface{} {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	var result []map[string]interface{}
	for _, room := range rm.rooms {
		room.mu.RLock()
		status := room.Status
		result = append(result, map[string]interface{}{
			"id":          room.ID,
			"host_id":     room.HostID,
			"guest_id":    room.GuestID,
			"status":      string(status),
			"player_count": playerCount(room),
			"created_at":  room.CreatedAt.Unix(),
		})
		room.mu.RUnlock()
	}
	return result
}

func playerCount(r *GameRoom) int {
	if r.GuestID != 0 {
		return 2
	}
	return 1
}

// ToMap 将房间信息导出为 map（线程安全）
func (r *GameRoom) ToMap() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return map[string]interface{}{
		"id":             r.ID,
		"host_id":        r.HostID,
		"guest_id":       r.GuestID,
		"host_char_id":   r.HostCharID,
		"guest_char_id":  r.GuestCharID,
		"host_skin_id":   r.HostSkinID,
		"guest_skin_id":  r.GuestSkinID,
		"host_ready":     r.HostReady,
		"guest_ready":    r.GuestReady,
		"status":         string(r.Status),
		"created_at":     r.CreatedAt.Unix(),
	}
}

// RemoveIdleRooms 移除闲置超过指定时长的房间
func (rm *RoomManager) RemoveIdleRooms(maxIdle time.Duration) int {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	cutoff := time.Now().Add(-maxIdle)
	removed := 0
	for id, room := range rm.rooms {
		room.mu.RLock()
		isPlaying := room.Status == StatusPlaying
		created := room.CreatedAt
		hostID := room.HostID
		guestID := room.GuestID
		room.mu.RUnlock()

		if !isPlaying && created.Before(cutoff) {
			delete(rm.rooms, id)
			delete(rm.userRoom, hostID)
			if guestID != 0 {
				delete(rm.userRoom, guestID)
			}
			removed++
		}
	}
	return removed
}

// 错误定义
var (
	ErrRoomNotFound      = &wsError{"房间不存在"}
	ErrRoomNotAvailable  = &wsError{"房间不可加入"}
	ErrRoomFull          = &wsError{"房间已满"}
	ErrSamePlayer        = &wsError{"不能加入自己的房间"}
	ErrNotInRoom         = &wsError{"不在任何房间中"}
	ErrNotAuthorized     = &wsError{"未授权，请先登录"}
)

type wsError struct {
	message string
}

func (e *wsError) Error() string { return e.message }
func (e *wsError) ToJSON() []byte {
	b, _ := json.Marshal(map[string]string{"type": "error", "message": e.message})
	return b
}
