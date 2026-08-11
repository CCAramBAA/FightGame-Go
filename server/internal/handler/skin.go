package handler

import (
	"net/http"
	"strconv"

	"fightgame-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SkinHandler struct {
	db *gorm.DB
}

func NewSkinHandler(db *gorm.DB) *SkinHandler {
	return &SkinHandler{db: db}
}

// ListSkins 获取所有皮肤
func (h *SkinHandler) ListSkins(c *gin.Context) {
	var skins []map[string]interface{}
	h.db.Table("skins").Order("id ASC").Find(&skins)
	if skins == nil {
		skins = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": skins})
}

// CreateSkin 创建皮肤（管理后台）
func (h *SkinHandler) CreateSkin(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := h.db.Table("skins").Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功"})
}

// UpdateSkin 更新皮肤（管理后台）
func (h *SkinHandler) UpdateSkin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := h.db.Table("skins").Where("id = ?", id).Updates(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

// DeleteSkin 删除皮肤（管理后台）
func (h *SkinHandler) DeleteSkin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Table("skins").Where("id = ?", id).Delete(nil)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// GetUserSkins 获取用户拥有的皮肤
func (h *SkinHandler) GetUserSkins(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		// 使用当前登录用户的 ID
		authenticatedID := middleware.GetUserID(c)
		if authenticatedID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少用户ID"})
			return
		}
		userID = strconv.FormatUint(uint64(authenticatedID), 10)
	}
	var skins []map[string]interface{}
	h.db.Table("skins s").
		Select("s.*, us.id as owned").
		Joins("LEFT JOIN user_skins us ON us.skin_id = s.id AND us.user_id = ?", userID).
		Order("s.id ASC").
		Scan(&skins)
	if skins == nil {
		skins = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": skins})
}
