package models

import (
	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null;uniqueIndex" json:"name"`
	DisplayName string         `gorm:"size:100" json:"display_name"`
	Description string         `gorm:"type:text" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Permission 权限模型
type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null;uniqueIndex" json:"name"`
	DisplayName string    `gorm:"size:100" json:"display_name"`
	Description string    `gorm:"type:text" json:"description"`
	Resource    string    `gorm:"size:50;index" json:"resource"`     // 资源：project, pipeline, deploy, user, system
	Action      string    `gorm:"size:50" json:"action"`         // 操作：create, read, update, delete, *
	Roles       []Role   `gorm:"many2many:role_permissions;" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID   uint      `gorm:"not null;uniqueIndex:idx_user_role" json:"user_id"`
	RoleID   uint      `gorm:"not null;uniqueIndex:idx_user_role" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ========== RBAC 预定义角色 ==========

// SuperAdmin 超级管理员（所有权限）
const RoleSuperAdmin = "super_admin"

// Admin 管理员（大部分权限）
const RoleAdmin = "admin"

// Developer 开发者（读写自己的项目）
const RoleDeveloper = "developer"

// Viewer 观察者（只读）
const RoleViewer = "viewer"

// ========== RBAC 预定义权限 ==========

// 项目管理权限
const (
	PermProjectCreate = "project:create"
	PermProjectRead   = "project:read"
	PermProjectUpdate = "project:update"
	PermProjectDelete = "project:delete"
	PermProjectAll    = "project:*"
)

// 流水线管理权限
const (
	PermPipelineCreate = "pipeline:create"
	PermPipelineRead   = "pipeline:read"
	PermPipelineUpdate = "pipeline:update"
	PermPipelineDelete = "pipeline:delete"
	PermPipelineTrigger = "pipeline:trigger"
	PermPipelineAll    = "pipeline:*"
)

// 部署管理权限
const (
	PermDeployCreate = "deploy:create"
	PermDeployRead   = "deploy:read"
	PermDeployUpdate = "deploy:update"
	PermDeployDelete = "deploy:delete"
	PermDeployExecute = "deploy:execute"
	PermDeployAll    = "deploy:*"
)

// 用户管理权限
const (
	PermUserCreate = "user:create"
	PermUserRead   = "user:read"
	PermUserUpdate = "user:update"
	PermUserDelete = "user:delete"
	PermUserAll    = "user:*"
)

// 系统管理权限
const (
	PermSystemSettings = "system:settings"
	PermSystemLogs    = "system:logs"
	PermSystemAll     = "system:*"
)

// ========== 模型方法 ==========

// HasPermission 检查角色是否有某个权限
func (r *Role) HasPermission(permissionName string) bool {
	for _, perm := range r.Permissions {
		if perm.Name == permissionName {
			return true
		}
		// 通配符检查（例如 project:* 匹配 project:create）
		if perm.Name == r.Resource+":*" {
			return true
		}
	}
	return false
}

// GetUserRoles 获取用户的所有角色
func GetUserRoles(db *gorm.DB, userID uint) ([]Role, error) {
	var roles []Role
	err := db.Model(&UserRole{}).Where("user_id = ?", userID).
		Preload("Permissions").
		Find(&roles).Error
	return roles, err
}

// HasUserPermission 检查用户是否有某个权限
func HasUserPermission(db *gorm.DB, userID uint, permissionName string) (bool, error) {
	roles, err := GetUserRoles(db, userID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role.HasPermission(permissionName) {
			return true, nil
		}
	}
	return false, nil
}

// InitializeRBAC 初始化 RBAC（创建默认角色和权限）
func InitializeRBAC(db *gorm.DB) error {
	// 1. 创建默认权限
	defaultPermissions := []Permission{
		// 项目管理
		{Name: PermProjectCreate, DisplayName: "创建项目", Resource: "project", Action: "create"},
		{Name: PermProjectRead, DisplayName: "查看项目", Resource: "project", Action: "read"},
		{Name: PermProjectUpdate, DisplayName: "更新项目", Resource: "project", Action: "update"},
		{Name: PermProjectDelete, DisplayName: "删除项目", Resource: "project", Action: "delete"},
		{Name: PermProjectAll, DisplayName: "项目管理（全部）", Resource: "project", Action: "*"},

		// 流水线管理
		{Name: PermPipelineCreate, DisplayName: "创建流水线", Resource: "pipeline", Action: "create"},
		{Name: PermPipelineRead, DisplayName: "查看流水线", Resource: "pipeline", Action: "read"},
		{Name: PermPipelineUpdate, DisplayName: "更新流水线", Resource: "pipeline", Action: "update"},
		{Name: PermPipelineDelete, DisplayName: "删除流水线", Resource: "pipeline", Action: "delete"},
		{Name: PermPipelineTrigger, DisplayName: "触发流水线", Resource: "pipeline", Action: "trigger"},
		{Name: PermPipelineAll, DisplayName: "流水线管理（全部）", Resource: "pipeline", Action: "*"},

		// 部署管理
		{Name: PermDeployCreate, DisplayName: "创建部署环境", Resource: "deploy", Action: "create"},
		{Name: PermDeployRead, DisplayName: "查看部署", Resource: "deploy", Action: "read"},
		{Name: PermDeployUpdate, DisplayName: "更新部署", Resource: "deploy", Action: "update"},
		{Name: PermDeployDelete, DisplayName: "删除部署", Resource: "deploy", Action: "delete"},
		{Name: PermDeployExecute, DisplayName: "执行部署", Resource: "deploy", Action: "execute"},
		{Name: PermDeployAll, DisplayName: "部署管理（全部）", Resource: "deploy", Action: "*"},

		// 用户管理
		{Name: PermUserCreate, DisplayName: "创建用户", Resource: "user", Action: "create"},
		{Name: PermUserRead, DisplayName: "查看用户", Resource: "user", Action: "read"},
		{Name: PermUserUpdate, DisplayName: "更新用户", Resource: "user", Action: "update"},
		{Name: PermUserDelete, DisplayName: "删除用户", Resource: "user", Action: "delete"},
		{Name: PermUserAll, DisplayName: "用户管理（全部）", Resource: "user", Action: "*"},

		// 系统管理
		{Name: PermSystemSettings, DisplayName: "系统设置", Resource: "system", Action: "settings"},
		{Name: PermSystemLogs, DisplayName: "查看日志", Resource: "system", Action: "logs"},
		{Name: PermSystemAll, DisplayName: "系统管理（全部）", Resource: "system", Action: "*"},
	}

	for _, perm := range defaultPermissions {
		db.Where("name = ?", perm.Name).FirstOrCreate(&perm)
	}

	// 2. 创建默认角色
	superAdminRole := Role{
		Name:        RoleSuperAdmin,
		DisplayName: "超级管理员",
		Description: "拥有所有权限",
	}
	adminRole := Role{
		Name:        RoleAdmin,
		DisplayName: "管理员",
		Description: "拥有大部分管理权限",
	}
	developerRole := Role{
		Name:        RoleDeveloper,
		DisplayName: "开发者",
		Description: "可以创建和管理自己的项目",
	}
	viewerRole := Role{
		Name:        RoleViewer,
		DisplayName: "观察者",
		Description: "只能查看，不能修改",
	}

	// 创建角色
	db.Where("name = ?", RoleSuperAdmin).FirstOrCreate(&superAdminRole)
	db.Where("name = ?", RoleAdmin).FirstOrCreate(&adminRole)
	db.Where("name = ?", RoleDeveloper).FirstOrCreate(&developerRole)
	db.Where("name = ?", RoleViewer).FirstOrCreate(&viewerRole)

	// 3. 为超级管理员分配所有权限
	var allPermissions []Permission
	db.Find(&allPermissions)
	db.Model(&superAdminRole).Association("Permissions").Replace(allPermissions)

	// 4. 为管理员分配大部分权限（除了系统管理）
	var adminPermissions []Permission
	db.Where("name NOT LIKE ? AND name NOT LIKE ?", "system:%", "user:%").Find(&adminPermissions)
	db.Model(&adminRole).Association("Permissions").Replace(adminPermissions)

	// 5. 为开发者分配项目、流水线、部署权限
	var devPermissions []Permission
	db.Where("name LIKE ? OR name LIKE ? OR name LIKE ?", "project:%", "pipeline:%", "deploy:%").
		Where("name NOT LIKE ?", "%delete").
		Find(&devPermissions)
	db.Model(&developerRole).Association("Permissions").Replace(devPermissions)

	// 6. 为观察者分配只读权限
	var viewerPermissions []Permission
	db.Where("name LIKE ? OR name LIKE ? OR name LIKE ? OR name LIKE ?",
		"project:read", "pipeline:read", "deploy:read", "user:read").
		Find(&viewerPermissions)
	db.Model(&viewerRole).Association("Permissions").Replace(viewerPermissions)

	// 7. 为 admin 用户分配超级管理员角色
	var adminUser User
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err == nil {
		// 检查是否已有角色
		var count int64
		db.Model(&UserRole{}).Where("user_id = ?", adminUser.ID).Count(&count)
		if count == 0 {
			db.Create(&UserRole{UserID: adminUser.ID, RoleID: superAdminRole.ID})
		}
	}

	return nil
}
