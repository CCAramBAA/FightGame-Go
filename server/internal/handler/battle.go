package handler

import (
	"net/http"
	"strconv"

	"fightgame-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BattleHandler struct {
	db *gorm.DB
}

func NewBattleHandler(db *gorm.DB) *BattleHandler {
	return &BattleHandler{db: db}
}

// RecordBattle 记录对局结果（PVP/PVE）
func (h *BattleHandler) RecordBattle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req struct {
		Mode      string `json:"mode" binding:"required"` // "pvp" / "pve"
		Result    string `json:"result"`                    // "win" / "lose" / "draw"
		GoldEarned int   `json:"gold_earned"`
		RankChange  int   `json:"rank_change"`
		RoomID      string `json:"room_id"`
		CharacterID uint   `json:"character_id"`
		OpponentID  uint   `json:"opponent_id"`
		FrameData   string `json:"frame_data"` // 对局帧数据JSON
		StageID     uint   `json:"stage_id"`   // PVE关卡ID
		Stars       int    `json:"stars"`      // PVE星级
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// PVP 发金币
	if req.Mode == "pvp" && req.Result == "win" && req.GoldEarned > 0 {
		h.db.Table("users").Where("id = ?", userID).
			Update("gold", gorm.Expr("gold + ?", req.GoldEarned))

		h.db.Table("gold_transactions").Create(map[string]interface{}{
			"user_id": userID,
			"amount":  req.GoldEarned,
			"type":    "pvp_win",
			"item_id": req.CharacterID,
		})

		// 更新胜场
		h.db.Table("rank_scores").Where("user_id = ?", userID).
			Update("win_count", gorm.Expr("win_count + 1"))
	}

	if req.Mode == "pvp" && req.Result == "lose" {
		h.db.Table("rank_scores").Where("user_id = ?", userID).
			Update("lose_count", gorm.Expr("lose_count + 1"))
	}

	// 更新段位积分
	if req.RankChange != 0 {
		h.db.Table("rank_scores").Where("user_id = ?", userID).
			Update("score", gorm.Expr("score + ?", req.RankChange))
	}

	// 保存对局记录
	record := map[string]interface{}{
		"user_id":      userID,
		"mode":         req.Mode,
		"result":       req.Result,
		"gold_earned":  req.GoldEarned,
		"rank_change":  req.RankChange,
		"room_id":      req.RoomID,
		"character_id": req.CharacterID,
		"opponent_id":  req.OpponentID,
		"frame_data":   req.FrameData,
		"stage_id":     req.StageID,
		"stars":        req.Stars,
	}
	if err := h.db.Table("battle_records").Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "记录成功"})
}

// GetBattleRecords 获取对局记录
func (h *BattleHandler) GetBattleRecords(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var records []map[string]interface{}
	h.db.Table("battle_records").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(20).
		Find(&records)
	if records == nil {
		records = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": records})
}

// GetReplay 获取对局回放
func (h *BattleHandler) GetReplay(c *gin.Context) {
	id := c.Param("id")
	recordID, _ := strconv.ParseUint(id, 10, 64)

	var record map[string]interface{}
	if err := h.db.Table("battle_records").Where("id = ?", recordID).Take(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "回放不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": record})
}
