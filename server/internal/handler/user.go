package handler

import (
	"net/http"
	"time"

	"fightgame-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db    *gorm.DB
	rdb   interface{} // *redis.Client, optional
	jwt   []byte
}

func NewUserHandler(db *gorm.DB, rdb interface{}, jwtSecret string) *UserHandler {
	return &UserHandler{db: db, rdb: rdb, jwt: []byte(jwtSecret)}
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	var existing map[string]interface{}
	if err := h.db.Table("users").Where("username = ?", req.Username).Take(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "用户名已存在"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注册失败"})
		return
	}

	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	user := map[string]interface{}{
		"username":      req.Username,
		"password_hash": string(hashed),
		"nickname":      nickname,
		"gold":          200,
		"role":          "user",
		"status":        1,
		"created_at":    time.Now(),
		"updated_at":    time.Now(),
	}

	if err := h.db.Table("users").Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注册失败"})
		return
	}

	// 创建段位积分记录
	userID := user["id"]
	h.db.Table("rank_scores").Create(map[string]interface{}{
		"user_id":    userID,
		"score":      1000,
		"win_count":  0,
		"lose_count": 0,
		"updated_at": time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "注册成功"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var user struct {
		ID       uint
		Username string
		Nickname string
		Password string
		IsAdmin  bool
	}

	if err := h.db.Table("users").
		Where("username = ?", req.Username).
		Select("id, username, nickname, password, is_admin").
		First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "账号或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "账号或密码错误"})
		return
	}

	role := "user"
	if user.IsAdmin {
		role = "admin"
	}

	token, err := h.generateToken(user.ID, user.Username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"token":    token,
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"role":     role,
		},
	})
}

// Logout 退出登录，将 token 加入 Redis 黑名单
func (h *UserHandler) Logout(c *gin.Context) {
	tokenStr := middleware.GetToken(c)
	userID := middleware.GetUserID(c)

	if h.rdb != nil {
		// 解析 token 获取过期时间
		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return h.jwt, nil
		})
		if token != nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				exp := int64(claims["exp"].(float64))
				ttl := time.Duration(exp-time.Now().Unix()) * time.Second
				if ttl > 0 {
					// 将 token 加入 Redis 黑名单
					rdb := h.rdb.(interface {
						Set(key string, value interface{}, expiration time.Duration) error
					})
					rdb.Set("blacklist:"+tokenStr, userID, ttl)
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "退出成功"})
}

// RefreshToken 刷新 token
func (h *UserHandler) RefreshToken(c *gin.Context) {
	userID := middleware.GetUserID(c)
	username := middleware.GetUsername(c)
	role := middleware.GetRole(c)

	token, err := h.generateToken(userID, username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "刷新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"token": token}})
}

func (h *UserHandler) generateToken(userID uint, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.jwt)
}
