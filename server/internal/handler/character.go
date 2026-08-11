package handler

import (
	"net/http"

	"fightgame-server/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CharacterHandler 角色/技能处理器
type CharacterHandler struct {
	db *gorm.DB
}

// NewCharacterHandler 创建角色处理器
func NewCharacterHandler(db *gorm.DB) *CharacterHandler {
	return &CharacterHandler{db: db}
}

// ListCharacters 获取所有角色列表
func (h *CharacterHandler) ListCharacters(c *gin.Context) {
	var characters []model.Character
	h.db.Order("id ASC").Find(&characters)
	c.JSON(http.StatusOK, characters)
}

// GetCharacter 获取单个角色详情（含技能列表 + 特殊规则）
func (h *CharacterHandler) GetCharacter(c *gin.Context) {
	id := c.Param("id")

	var character model.Character
	if err := h.db.First(&character, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	var skills []model.Skill
	h.db.Where("character_id = ?", id).Find(&skills)

	var rules []model.HeroSpecialRule
	h.db.Where("character_id = ?", id).Order("priority_order ASC").Find(&rules)

	c.JSON(http.StatusOK, gin.H{
		"character": character,
		"skills":    skills,
		"rules":     rules,
	})
}

// ListSkills 获取某角色的所有技能
func (h *CharacterHandler) ListSkills(c *gin.Context) {
	charID := c.Param("id")
	var skills []model.Skill
	h.db.Where("character_id = ?", charID).Find(&skills)
	c.JSON(http.StatusOK, skills)
}
