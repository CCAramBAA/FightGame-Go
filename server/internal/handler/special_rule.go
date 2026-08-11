package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SpecialRuleHandler struct {
	db *gorm.DB
}

func NewSpecialRuleHandler(db *gorm.DB) *SpecialRuleHandler {
	return &SpecialRuleHandler{db: db}
}

// ListRules 获取所有英雄特殊交互规则
func (h *SpecialRuleHandler) ListRules(c *gin.Context) {
	var rules []map[string]interface{}
	h.db.Table("hero_special_rules").Order("priority ASC").Find(&rules)
	if rules == nil {
		rules = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": rules})
}

// GetRule 获取单条规则
func (h *SpecialRuleHandler) GetRule(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var rule map[string]interface{}
	if err := h.db.Table("hero_special_rules").Where("id = ?", id).Take(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "规则不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": rule})
}

// CreateRule 创建规则
func (h *SpecialRuleHandler) CreateRule(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := h.db.Table("hero_special_rules").Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功"})
}

// UpdateRule 更新规则
func (h *SpecialRuleHandler) UpdateRule(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := h.db.Table("hero_special_rules").Where("id = ?", id).Updates(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

// DeleteRule 删除规则
func (h *SpecialRuleHandler) DeleteRule(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Table("hero_special_rules").Where("id = ?", id).Delete(nil)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ReorderRules 重排序规则
func (h *SpecialRuleHandler) ReorderRules(c *gin.Context) {
	var req []struct {
		ID       uint `json:"id"`
		Priority int  `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	for _, r := range req {
		h.db.Table("hero_special_rules").Where("id = ?", r.ID).Update("priority", r.Priority)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "排序更新成功"})
}
