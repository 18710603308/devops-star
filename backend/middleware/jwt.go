package middleware

import (
	"net/http"
	"strings"

	"devops-star/backend/config"
	"devops-star/backend/services"

	"github.com/gin-gonic/gin"
)

// JWT 中间件（验证 Token）
func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从 Header 获取 Authorization
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			ctx.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		// 验证 Token
		userID, username, err := services.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "认证失败：" + err.Error()})
			ctx.Abort()
			return
		}

		// 将用户信息存入上下文
		ctx.Set("user_id", userID)
		ctx.Set("username", username)

		ctx.Next()
	}
}

// 获取当前用户 ID（从上下文）
func GetUserID(ctx *gin.Context) uint {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}

// 获取当前用户名（从上下文）
func GetUsername(ctx *gin.Context) string {
	username, exists := ctx.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}
