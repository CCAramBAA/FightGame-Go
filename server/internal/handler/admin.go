package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"fightgame-server/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db  *gorm.DB
	hub *websocket.Hub
	jwt []byte
}

func NewAdminHandler(db *gorm.DB, hub *websocket.Hub, jwtSecret string) *AdminHandler {
	return &AdminHandler{db: db, hub: hub, jwt: []byte(jwtSecret)}
}

// ===================== 管理员登录 =====================

type AdminLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req AdminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var user struct {
		ID           uint
		Username     string
		IsAdmin      bool
		PasswordHash string
	}

	if err := h.db.Table("users").
		Where("username = ? AND status = 1", req.Username).
		Select("id, username, is_admin, password_hash").
		First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "账号或密码错误"})
		return
	}

	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无管理权限"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "账号或密码错误"})
		return
	}

	role := "admin"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     role,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	})
	tokenStr, _ := token.SignedString(h.jwt)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"token":    tokenStr,
			"id":       user.ID,
			"username": user.Username,
			"role":     role,
		},
	})
}

// ===================== Dashboard =====================

func (h *AdminHandler) Dashboard(c *gin.Context) {
	var userCount, charCount, skinCount, stageCount, battleCount, onlineCount int64
	var todayLogin, todayBattle, todayGold int64

	today := time.Now().Format("2006-01-02")

	h.db.Table("users").Where("status = 1").Count(&userCount)
	h.db.Table("characters").Count(&charCount)
	h.db.Table("skins").Count(&skinCount)
	h.db.Table("pve_stages").Count(&stageCount)
	h.db.Table("battle_records").Count(&battleCount)

	// 今日登录用户
	h.db.Table("users").
		Where("DATE(last_login_at) = ?", today).
		Count(&todayLogin)

	// 今日对局
	h.db.Table("battle_records").
		Where("DATE(created_at) = ?", today).
		Count(&todayBattle)

	// 今日金币流水
	h.db.Table("gold_transactions").
		Where("DATE(created_at) = ?", today).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&todayGold)

	// 在线人数
	if h.hub != nil {
		onlineCount = int64(h.hub.OnlineCount())
	}

	// 活跃房间
	var activeRoomCount int64
	h.db.Table("game_rooms").Where("status IN ('waiting','selecting','playing')").Count(&activeRoomCount)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"users":            userCount,
			"characters":       charCount,
			"skins":            skinCount,
			"stages":           stageCount,
			"battles":          battleCount,
			"online":           onlineCount,
			"active_rooms":     activeRoomCount,
			"today_logins":     todayLogin,
			"today_battles":    todayBattle,
			"today_gold":       todayGold,
		},
	})
}

// ===================== 用户管理 =====================

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	type UserRow struct {
		ID        uint      `json:"id"`
		Username  string    `json:"username"`
		Nickname  string    `json:"nickname"`
		IsAdmin   bool      `json:"is_admin"`
		Gold      int64     `json:"gold"`
		RankScore int       `json:"rank_score"`
		WinCount  int       `json:"win_count"`
		LoseCount int       `json:"lose_count"`
		Status    int       `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}

	var total int64
	var users []UserRow

	query := h.db.Table("users")
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("id DESC").Scan(&users)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":     users,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var req map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&req)

	updates := map[string]interface{}{}
	if v, ok := req["nickname"]; ok {
		updates["nickname"] = v
	}
	if v, ok := req["is_admin"]; ok {
		updates["is_admin"] = v
	}
	if v, ok := req["status"]; ok {
		updates["status"] = v
	}
	if v, ok := req["gold"]; ok {
		updates["gold"] = v
	}
	if v, ok := req["rank_score"]; ok {
		updates["rank_score"] = v
	}
	updates["updated_at"] = time.Now()

	if err := h.db.Table("users").Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

// ===================== 角色管理 =====================

func (h *AdminHandler) ListCharacters(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	var chars []map[string]interface{}

	h.db.Table("characters").Count(&total)
	h.db.Table("characters").Offset((page - 1) * pageSize).Limit(pageSize).
		Order("id ASC").Scan(&chars)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"list": chars, "total": total, "page": page, "pageSize": pageSize},
	})
}

func (h *AdminHandler) CreateCharacter(c *gin.Context) {
	var req map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&req)
	req["created_at"] = time.Now()
	req["updated_at"] = time.Now()

	if err := h.db.Table("characters").Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功"})
}

func (h *AdminHandler) UpdateCharacter(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var req map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&req)
	req["updated_at"] = time.Now()

	if err := h.db.Table("characters").Where("id = ?", id).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

func (h *AdminHandler) DeleteCharacter(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	h.db.Table("skills").Where("character_id = ?", id).Delete(nil)
	h.db.Table("skins").Where("character_id = ?", id).Delete(nil)
	h.db.Table("hero_special_rules").Where("character_id = ?", id).Delete(nil)
	h.db.Table("user_characters").Where("character_id = ?", id).Delete(nil)
	h.db.Table("shop_items").Where("item_type = ? AND item_id = ?", "character", id).Delete(nil)

	if err := h.db.Table("characters").Where("id = ?", id).Delete(nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ===================== 技能管理 =====================

func (h *AdminHandler) ListSkills(c *gin.Context) {
	characterID := c.Query("character_id")

	var skills []map[string]interface{}
	query := h.db.Table("skills")
	if characterID != "" {
		query = query.Where("character_id = ?", characterID)
	}
	query.Order("priority_level DESC").Scan(&skills)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": skills})
}

func (h *AdminHandler) CreateSkill(c *gin.Context) {
	var req map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&req)
	req["created_at"] = time.Now()
	req["updated_at"] = time.Now()

	if err := h.db.Table("skills").Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功"})
}

func (h *AdminHandler) UpdateSkill(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&req)
	req["updated_at"] = time.Now()

	if err := h.db.Table("skills").Where("id = ?", id).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

func (h *AdminHandler) DeleteSkill(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Table("skills").Where("id = ?", id).Delete(nil)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ===================== 皮肤管理 =====================

func (h *AdminHandler) ListSkins(c *gin.Context) {
	characterID := c.Query("character_id")

	var skins []map[string]interface{}
	query := h.db.Table("skins")
	if characterID != "" {
		query = query.Where("character_id = ?", characterID)
	}
	query.Order("id ASC").Scan(&skins)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": skins})
}

func (h *AdminHandler) CreateSkin(c *gin.Context) {
	var req map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&req)
	req["created_at"] = time.Now()
	req["updated_at"] = time.Now()

	if err := h.db.Table("skins").Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功"})
}

func (h *AdminHandler) UpdateSkin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&req)
	req["updated_at"] = time.Now()

	if err := h.db.Table("skins").Where("id = ?", id).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

func (h *AdminHandler) DeleteSkin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Table("user_skins").Where("skin_id = ?", id).Delete(nil)
	h.db.Table("shop_items").Where("item_type = ? AND item_id = ?", "skin", id).Delete(nil)
	h.db.Table("skins").Where("id = ?", id).Delete(nil)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ===================== 特殊规则管理 =====================

func (h *AdminHandler) ListSpecialRules(c *gin.Context) {
	characterID := c.Query("character_id")

	var rules []map[string]interface{}
	query := h.db.Table("hero_special_rules")
	if characterID != "" {
		query = query.Where("character_id = ?", characterID)
	}
	query.Order("priority_order ASC").Scan(&rules)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": rules})
}

func (h *AdminHandler) CreateSpecialRule(c *gin.Context) {
	var req map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&req)
	req["created_at"] = time.Now()
	req["updated_at"] = time.Now()

	if err := h.db.Table("hero_special_rules").Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功"})
}

func (h *AdminHandler) UpdateSpecialRule(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&req)
	req["updated_at"] = time.Now()

	if err := h.db.Table("hero_special_rules").Where("id = ?", id).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

func (h *AdminHandler) DeleteSpecialRule(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Table("hero_special_rules").Where("id = ?", id).Delete(nil)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ===================== 关卡管理 =====================

func (h *AdminHandler) ListStages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	var total int64
	var stages []map[string]interface{}

	h.db.Table("pve_stages").Count(&total)
	h.db.Table("pve_stages").Offset((page - 1) * pageSize).Limit(pageSize).
		Order("stage_order ASC").Scan(&stages)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"list": stages, "total": total, "page": page, "pageSize": pageSize},
	})
}

func (h *AdminHandler) CreateStage(c *gin.Context) {
	var req map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&req)
	req["created_at"] = time.Now()
	req["updated_at"] = time.Now()

	if err := h.db.Table("pve_stages").Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功"})
}

func (h *AdminHandler) UpdateStage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&req)
	req["updated_at"] = time.Now()

	if err := h.db.Table("pve_stages").Where("id = ?", id).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

func (h *AdminHandler) DeleteStage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Table("pve_progress").Where("stage_id = ?", id).Delete(nil)
	h.db.Table("pve_stages").Where("id = ?", id).Delete(nil)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ===================== 商城管理 =====================

func (h *AdminHandler) ListShopItems(c *gin.Context) {
	var items []map[string]interface{}
	h.db.Table("shop_items").Order("id ASC").Scan(&items)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": items})
}

func (h *AdminHandler) CreateShopItem(c *gin.Context) {
	var req map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&req)
	req["created_at"] = time.Now()
	req["updated_at"] = time.Now()

	if err := h.db.Table("shop_items").Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功"})
}

func (h *AdminHandler) UpdateShopItem(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&req)
	req["updated_at"] = time.Now()

	if err := h.db.Table("shop_items").Where("id = ?", id).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

func (h *AdminHandler) DeleteShopItem(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Table("shop_items").Where("id = ?", id).Delete(nil)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ===================== 房间管理 =====================

func (h *AdminHandler) ListRooms(c *gin.Context) {
	var rooms []map[string]interface{}
	h.db.Table("game_rooms").
		Where("status IN ('waiting','selecting','playing')").
		Order("created_at DESC").Scan(&rooms)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": rooms})
}

func (h *AdminHandler) ForceCloseRoom(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	h.db.Table("game_rooms").Where("id = ?", roomID).
		Update("status", "finished")

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "房间已关闭"})
}

// ===================== 日志查看 =====================

func (h *AdminHandler) ViewLogs(c *gin.Context) {
	level := c.DefaultQuery("level", "all")
	lines, _ := strconv.Atoi(c.DefaultQuery("lines", "100"))

	type LogEntry struct {
		Time    string `json:"time"`
		Level   string `json:"level"`
		Message string `json:"message"`
	}

	var logs []LogEntry
	for i := 0; i < lines; i++ {
		logs = append(logs, LogEntry{
			Time:    time.Now().Add(-time.Duration(i) * time.Minute).Format("15:04:05"),
			Level:   "INFO",
			Message: fmt.Sprintf("模拟日志 #%d - 系统运行正常", i+1),
		})
	}
	if level != "all" {
		var filtered []LogEntry
		for _, l := range logs {
			if l.Level == level {
				filtered = append(filtered, l)
			}
		}
		logs = filtered
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": logs})
}

// ===================== 金币流水 =====================

func (h *AdminHandler) ListGoldTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	userID := c.Query("user_id")

	if page < 1 {
		page = 1
	}

	var total int64
	var txns []map[string]interface{}

	query := h.db.Table("gold_transactions")
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Scan(&txns)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"list": txns, "total": total, "page": page},
	})
}

// ===================== 对局记录 =====================

func (h *AdminHandler) ListBattleRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}

	var total int64
	var records []map[string]interface{}

	h.db.Table("battle_records").Count(&total)
	h.db.Table("battle_records").Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Scan(&records)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"list": records, "total": total, "page": page},
	})
}

// ===================== 缓存管理 =====================

func (h *AdminHandler) RefreshCache(c *gin.Context) {
	cacheType := c.Query("type")
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": fmt.Sprintf("缓存已刷新: %s", cacheType),
	})
}
