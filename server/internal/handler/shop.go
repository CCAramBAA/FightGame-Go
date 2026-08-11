package handler

import (
	"net/http"

	"fightgame-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopHandler struct {
	db *gorm.DB
}

func NewShopHandler(db *gorm.DB) *ShopHandler {
	return &ShopHandler{db: db}
}

// ListShopItems 获取商城商品列表（含名称描述）
func (h *ShopHandler) ListShopItems(c *gin.Context) {
	type ShopItemDisplay struct {
		ID          uint   `json:"id"`
		ItemType    string `json:"item_type"`
		ItemID      uint   `json:"item_id"`
		Price       int64  `json:"price"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var items []ShopItemDisplay
	// 合并角色数据
	rows, err := h.db.Raw(`
		SELECT si.id, si.item_type, si.item_id, si.price,
			COALESCE(c.name, s.name, '') as name,
			COALESCE(c.description, s.name, '') as description
		FROM shop_items si
		LEFT JOIN characters c ON si.item_type = 'character' AND si.item_id = c.id
		LEFT JOIN skins s ON si.item_type = 'skin' AND si.item_id = s.id
		ORDER BY si.id ASC
	`).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item ShopItemDisplay
		if err := h.db.ScanRows(rows, &item); err != nil {
			continue
		}
		items = append(items, item)
	}
	if items == nil {
		items = []ShopItemDisplay{}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": items})
}

// PurchaseItem 购买商品
func (h *ShopHandler) PurchaseItem(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req struct {
		ItemID uint `json:"item_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 查询商品
	var item map[string]interface{}
	if err := h.db.Table("shop_items").Where("id = ?", req.ItemID).Take(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "商品不存在"})
		return
	}

	price := item["price"].(int64)
	itemType := item["item_type"].(string)
	itemID := uint(item["item_id"].(int64))

	// 检查用户金币
	var user struct {
		ID   uint
		Gold int
	}
	if err := h.db.Table("users").Select("id, gold").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "用户不存在"})
		return
	}
	if user.Gold < int(price) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "金币不足"})
		return
	}

	// 检查是否已拥有
	if itemType == "character" {
		var count int64
		h.db.Table("user_characters").Where("user_id = ? AND character_id = ?", userID, itemID).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "已拥有该角色"})
			return
		}
	} else if itemType == "skin" {
		var count int64
		h.db.Table("user_skins").Where("user_id = ? AND skin_id = ?", userID, itemID).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "已拥有该皮肤"})
			return
		}
	}

	// 扣金币
	if err := h.db.Table("users").Where("id = ? AND gold >= ?", userID, price).
		Update("gold", gorm.Expr("gold - ?", price)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "购买失败"})
		return
	}

	// 记录金币流水
	goldFlow := map[string]interface{}{
		"user_id":    userID,
		"amount":     -price,
		"type":       "shop_purchase",
		"item_id":    itemID,
		"item_type":  itemType,
		"balance":    user.Gold - int(price),
	}
	h.db.Table("gold_transactions").Create(&goldFlow)

	// 发放物品
	if itemType == "character" {
		h.db.Table("user_characters").Create(map[string]interface{}{
			"user_id": userID, "character_id": itemID,
		})
	} else if itemType == "skin" {
		h.db.Table("user_skins").Create(map[string]interface{}{
			"user_id": userID, "skin_id": itemID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "购买成功", "data": gin.H{"gold": user.Gold - int(price)}})
}

// GetGoldTransactions 获取金币流水
func (h *ShopHandler) GetGoldTransactions(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var transactions []map[string]interface{}
	h.db.Table("gold_transactions").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(50).
		Find(&transactions)
	if transactions == nil {
		transactions = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": transactions})
}
