package handler

import (
	"net/http"

	"fightgame-server/internal/middleware"
	"fightgame-server/internal/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ProfileHandler 用户个人资料处理器
type ProfileHandler struct {
	db *gorm.DB
}

// NewProfileHandler 创建个人资料处理器
func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=255"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// GetProfile 获取用户个人信息
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	type ProfileResp struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Gold     int    `json:"gold"`
		Role     string `json:"role"`
	}

	var profile ProfileResp
	if err := h.db.Table("users").
		Select("id, username, nickname, avatar, gold, role").
		Where("id = ?", userID).
		First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": profile})
}

// GetRankScore 获取用户段位积分
func (h *ProfileHandler) GetRankScore(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var score struct {
		Score     int `json:"score"`
		WinCount  int `json:"win_count"`
		LoseCount int `json:"lose_count"`
	}

	h.db.Table("users").
		Select("rank_score as score, win_count, lose_count").
		Where("id = ?", userID).
		First(&score)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": score})
}

// UpdateProfile 更新昵称和头像
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无更新内容"})
		return
	}

	if err := h.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

// ChangePassword 修改密码
func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var currentHash string
	if err := h.db.Table("users").Where("id = ?", userID).Select("password_hash").Scan(&currentHash).Error; err != nil || currentHash == "" {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "原密码错误"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	if err := h.db.Table("users").Where("id = ?", userID).Update("password_hash", string(hashed)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "修改失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "密码修改成功"})
}

// GetMyCharacters 获取我的角色详情（含技能信息）
func (h *ProfileHandler) GetMyCharacters(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var userChars []model.UserCharacter
	h.db.Where("user_id = ?", userID).Find(&userChars)

	type CharacterDetail struct {
		ID          uint          `json:"id"`
		Name        string        `json:"name"`
		Description string        `json:"description"`
		HP          int           `json:"hp"`
		Energy      int           `json:"energy"`
		Speed       int           `json:"speed"`
		Attack      int           `json:"attack"`
		Defense     int           `json:"defense"`
		AvatarPath  string        `json:"avatar_path"`
		Skills      []model.Skill `json:"skills"`
	}

	var result []CharacterDetail
	for _, uc := range userChars {
		var char model.Character
		if err := h.db.First(&char, uc.CharacterID).Error; err != nil {
			continue
		}
		var skills []model.Skill
		h.db.Where("character_id = ?", uc.CharacterID).Find(&skills)

		result = append(result, CharacterDetail{
			ID:          char.ID,
			Name:        char.Name,
			Description: char.Description,
			HP:          char.HP,
			Energy:      char.Energy,
			Speed:       char.Speed,
			Attack:      char.Attack,
			Defense:     char.Defense,
			AvatarPath:  char.AvatarPath,
			Skills:      skills,
		})
	}

	if result == nil {
		result = []CharacterDetail{}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
}

// GetMySkins 获取我的皮肤列表
func (h *ProfileHandler) GetMySkins(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var userSkins []model.UserSkin
	h.db.Where("user_id = ?", userID).Find(&userSkins)

	type SkinDetail struct {
		ID          uint   `json:"id"`
		CharacterID uint   `json:"character_id"`
		Name        string `json:"name"`
		PreviewPath string `json:"preview_path"`
	}

	var result []SkinDetail
	for _, us := range userSkins {
		var skin model.Skin
		if err := h.db.First(&skin, us.SkinID).Error; err != nil {
			continue
		}
		result = append(result, SkinDetail{
			ID:          skin.ID,
			CharacterID: skin.CharacterID,
			Name:        skin.Name,
			PreviewPath: skin.PreviewPath,
		})
	}

	if result == nil {
		result = []SkinDetail{}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
}
