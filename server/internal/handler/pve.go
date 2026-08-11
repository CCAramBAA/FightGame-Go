package handler

import (
	"net/http"
	"strconv"

	"fightgame-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PVEHandler struct {
	db *gorm.DB
}

func NewPVEHandler(db *gorm.DB) *PVEHandler {
	return &PVEHandler{db: db}
}

// ListStages 获取关卡列表
func (h *PVEHandler) ListStages(c *gin.Context) {
	var stages []map[string]interface{}
	h.db.Table("pve_stages").Order("id ASC").Find(&stages)
	if stages == nil {
		stages = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": stages})
}

// GetStage 获取单个关卡
func (h *PVEHandler) GetStage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var stage map[string]interface{}
	if err := h.db.Table("pve_stages").Where("id = ?", id).Take(&stage).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "关卡不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": stage})
}

// SaveProgress 保存PVE进度
func (h *PVEHandler) SaveProgress(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req struct {
		StageID uint `json:"stage_id" binding:"required"`
		Stars   int  `json:"stars" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// Upsert progress
	var existingID uint
	h.db.Table("pve_progresses").
		Where("user_id = ? AND stage_id = ?", userID, req.StageID).
		Select("id").Scan(&existingID)

	if existingID > 0 {
		h.db.Table("pve_progresses").Where("id = ?", existingID).Updates(map[string]interface{}{
			"stars":    req.Stars,
			"cleared":  true,
		})
	} else {
		h.db.Table("pve_progresses").Create(map[string]interface{}{
			"user_id":  userID,
			"stage_id": req.StageID,
			"stars":    req.Stars,
			"cleared":  true,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "进度已保存"})
}

// GetProgress 获取用户PVE进度
func (h *PVEHandler) GetProgress(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var progress []map[string]interface{}
	h.db.Table("pve_progresses").Where("user_id = ?", userID).Find(&progress)
	if progress == nil {
		progress = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": progress})
}

// Admin: CreateStage 创建关卡
func (h *PVEHandler) CreateStage(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := h.db.Table("pve_stages").Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功"})
}

// Admin: UpdateStage 更新关卡
func (h *PVEHandler) UpdateStage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := h.db.Table("pve_stages").Where("id = ?", id).Updates(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

// Admin: DeleteStage 删除关卡
func (h *PVEHandler) DeleteStage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Table("pve_stages").Where("id = ?", id).Delete(nil)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}
