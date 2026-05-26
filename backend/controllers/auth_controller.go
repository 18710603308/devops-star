package controllers

import (
	"net/http"

	"devops-star/backend/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email" binding:"required,email"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "注册成功", "user": user})
}

// 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := c.authService.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// 获取当前用户信息
func (c *AuthController) GetMe(ctx *gin.Context) {
	// 从 JWT 中获取用户 ID（简化版）
	userID := uint(1) // 实际应从 JWT 中解析

	user, err := c.authService.GetUserInfo(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
