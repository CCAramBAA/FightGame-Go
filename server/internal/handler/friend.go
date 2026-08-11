package handler

import (
	"net/http"
	"strconv"

	"fightgame-server/internal/middleware"
	"fightgame-server/internal/websocket"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FriendHandler struct {
	db  *gorm.DB
	hub *websocket.Hub
}

func NewFriendHandler(db *gorm.DB, hub *websocket.Hub) *FriendHandler {
	return &FriendHandler{db: db, hub: hub}
}

// AddFriend 添加好友
func (h *FriendHandler) AddFriend(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req struct {
		FriendID uint `json:"friend_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if userID == req.FriendID {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不能添加自己为好友"})
		return
	}

	// 检查目标用户是否存在
	var targetExists int64
	h.db.Table("users").Where("id = ? AND status = 1", req.FriendID).Count(&targetExists)
	if targetExists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	// 检查是否已存在
	var count int64
	h.db.Table("friend_relations").Where(
		"(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		userID, req.FriendID, req.FriendID, userID,
	).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "已经是好友"})
		return
	}

	// 双向插入
	fr1 := map[string]interface{}{"user_id": userID, "friend_id": req.FriendID}
	fr2 := map[string]interface{}{"user_id": req.FriendID, "friend_id": userID}
	if err := h.db.Table("friend_relations").Create(&fr1).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "添加失败"})
		return
	}
	h.db.Table("friend_relations").Create(&fr2)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "添加成功"})
}

// RemoveFriend 删除好友
func (h *FriendHandler) RemoveFriend(c *gin.Context) {
	userID := middleware.GetUserID(c)
	friendID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if friendID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	h.db.Table("friend_relations").Where(
		"(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		userID, friendID, friendID, userID,
	).Delete(nil)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ListFriends 获取好友列表（含真实在线状态）
func (h *FriendHandler) ListFriends(c *gin.Context) {
	userID := middleware.GetUserID(c)

	type FriendInfo struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Online   bool   `json:"online"`
	}

	var friends []FriendInfo
	h.db.Table("friend_relations fr").
		Select("u.id, u.username, u.nickname").
		Joins("JOIN users u ON u.id = fr.friend_id").
		Where("fr.user_id = ? AND u.status = 1", userID).
		Scan(&friends)

	// 通过 Hub 查询真实在线状态
	for i := range friends {
		if h.hub != nil {
			friends[i].Online = h.hub.IsUserOnline(friends[i].ID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": friends})
}

// SendInvite 发送对战邀请
func (h *FriendHandler) SendInvite(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req struct {
		FriendID uint `json:"friend_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 检查好友关系
	var count int64
	h.db.Table("friend_relations").
		Where("user_id = ? AND friend_id = ?", userID, req.FriendID).
		Count(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不是好友"})
		return
	}

	// 检查对方是否在线
	if h.hub != nil && !h.hub.IsUserOnline(req.FriendID) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "好友不在线"})
		return
	}

	// 检查对方是否在战斗中
	if h.hub != nil && h.hub.RoomMgr.GetUserRoom(req.FriendID) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "好友正在对战中"})
		return
	}

	// TODO: 通过 WebSocket 向对方发送邀请弹窗消息

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "邀请已发送"})
}
