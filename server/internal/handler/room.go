package handler

import (
	"fmt"
	"net/http"
	"time"

	"fightgame-server/internal/middleware"
	"fightgame-server/internal/websocket"

	"github.com/gin-gonic/gin"
)

// RoomHandler 房间 REST API 处理器
type RoomHandler struct {
	hub *websocket.Hub
}

// NewRoomHandler 创建房间处理器
func NewRoomHandler(hub *websocket.Hub) *RoomHandler {
	return &RoomHandler{hub: hub}
}

// ListRooms 列出所有等待中的房间
func (h *RoomHandler) ListRooms(c *gin.Context) {
	rooms := h.hub.RoomMgr.ListRooms()
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": rooms})
}

// CreateRoom 创建房间
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 自动生成房间ID
	roomID := fmt.Sprintf("room_%d_%d", userID, time.Now().UnixNano())

	room := h.hub.RoomMgr.CreateRoom(roomID, userID)
	if room == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建房间失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": room.ToMap(), "message": "房间创建成功"})
}

// JoinRoom 加入房间
func (h *RoomHandler) JoinRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)
	roomID := c.Param("id")

	room, err := h.hub.RoomMgr.JoinRoom(roomID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": room.ToMap(), "message": "加入成功"})
}

// LeaveRoom 退出房间
func (h *RoomHandler) LeaveRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)

	_, ok := h.hub.RoomMgr.LeaveRoom(userID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "当前不在任何房间"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已退出房间"})
}

// GetRoom 获取单个房间信息
func (h *RoomHandler) GetRoom(c *gin.Context) {
	roomID := c.Param("id")
	room := h.hub.RoomMgr.GetRoom(roomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "房间不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": room.ToMap()})
}

// OnlineCount 在线人数
func (h *RoomHandler) OnlineCount(c *gin.Context) {
	count := h.hub.OnlineCount()
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"online": count}})
}
