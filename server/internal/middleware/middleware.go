package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "请先登录"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "无效的token格式"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "token已过期或无效"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "无效的token"})
			c.Abort()
			return
		}

		c.Set("token", tokenStr)
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("username", claims["username"].(string))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role == nil || (role != "admin" && role != "super_admin") {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// GetToken 从 Context 获取 token 字符串
func GetToken(c *gin.Context) string {
	token, _ := c.Get("token")
	return token.(string)
}

// GetUserID 从 Context 获取用户ID
func GetUserID(c *gin.Context) uint {
	id, _ := c.Get("user_id")
	return id.(uint)
}

// GetUsername 从 Context 获取用户名
func GetUsername(c *gin.Context) string {
	username, _ := c.Get("username")
	return username.(string)
}

// GetRole 从 Context 获取角色
func GetRole(c *gin.Context) string {
	role, _ := c.Get("role")
	return role.(string)
}
