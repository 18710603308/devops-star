package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseUint 将字符串转换为 uint，解析失败返回 0
func parseUint(s string) uint {
	id, _ := strconv.ParseUint(s, 10, 64)
	return uint(id)
}

// respondJSON 统一 JSON 响应（简化 helper）
func respondJSON(ctx *gin.Context, status int, data interface{}) {
	ctx.JSON(status, data)
}
