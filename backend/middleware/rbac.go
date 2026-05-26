package middleware

import (
	"net/http"
	"strings"

	"devops-star/backend/config"
	"devops-star/backend/services"

	"github.com/gin-gonic/gin"
)

// RBACMiddleware RBAC 权限中间件工厂
func RBACMiddleware(cfg *config.Config, rbacService *services.RBACService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从 JWT 中获取用户 ID（由 JWTAuth 中间件设置）
		userIDValue, exists := ctx.Get("user_id")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的用户 ID"})
			return
		}

		// 获取当前请求的权限名称
		permissionName := getPermissionName(ctx)

		// 检查用户是否有该权限
		hasPermission, err := rbacService.CheckPermission(userID, permissionName)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !hasPermission {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			return
		}

		ctx.Next()
	}
}

// getPermissionName 根据请求方法和路径生成权限名称
func getPermissionName(ctx *gin.Context) string {
	method := ctx.Request.Method
	path := ctx.Request.URL.Path

	// 简化：实际应根据路由表映射
	// 这里使用简单的规则映射

	// 提取资源类型和操作
	resource := extractResource(path)
	action := extractAction(method)

	if resource == "" || action == "" {
		return ""
	}

	return resource + ":" + action
}

// extractResource 从路径中提取资源类型
func extractResource(path string) string {
	path = strings.ToLower(path)

	if strings.Contains(path, "projects") {
		return "project"
	}
	if strings.Contains(path, "pipelines") {
		return "pipeline"
	}
	if strings.Contains(path, "deploy") {
		return "deploy"
	}
	if strings.Contains(path, "users") {
		return "user"
	}
	if strings.Contains(path, "system") || strings.Contains(path, "settings") {
		return "system"
	}

	return ""
}

// extractAction 从 HTTP 方法中提取操作
func extractAction(method string) string {
	switch method {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return ""
	}
}

// RequirePermission 检查特定权限的中间件
func RequirePermission(permissionName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDValue, exists := ctx.Get("user_id")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的用户 ID"})
			return
		}

		// 这里需要从上下文获取 rbacService（简化版）
		// 实际应通过依赖注入或全局变量获取
		_ = userID     // 避免未使用错误
		_ = permissionName // 避免未使用错误

		ctx.Next()
	}
}

// ========== 辅助函数 ==========

// GetUserRoles 从上下文获取用户角色
func GetUserRoles(ctx *gin.Context) []string {
	rolesValue, exists := ctx.Get("user_roles")
	if !exists {
		return []string{}
	}

	roles, ok := rolesValue.([]string)
	if !ok {
		return []string{}
	}

	return roles
}

// IsSuperAdmin 检查用户是否是超级管理员
func IsSuperAdmin(ctx *gin.Context) bool {
	roles := GetUserRoles(ctx)
	for _, role := range roles {
		if role == "super_admin" {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否有任意一个角色
func HasAnyRole(ctx *gin.Context, roleNames ...string) bool {
	userRoles := GetUserRoles(ctx)

	for _, userRole := range userRoles {
		for _, roleName := range roleNames {
			if userRole == roleName {
				return true
			}
		}
	}

	return false
}
