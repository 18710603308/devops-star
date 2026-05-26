package services

import (
	"devops-star/backend/models"
	"errors"

	"gorm.io/gorm"
)

// RBACService RBAC 权限服务
type RBACService struct {
	DB *gorm.DB
}

func NewRBACService(db *gorm.DB) *RBACService {
	return &RBACService{DB: db}
}

// ========== 角色管理 ==========

// GetRoles 获取所有角色
func (s *RBACService) GetRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := s.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRole 获取单个角色
func (s *RBACService) GetRole(roleID uint) (*models.Role, error) {
	var role models.Role
	if err := s.DB.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// CreateRole 创建角色
func (s *RBACService) CreateRole(name, displayName, description string, permissionIDs []uint) (*models.Role, error) {
	role := &models.Role{
		Name:        name,
		DisplayName: displayName,
		Description: description,
	}

	if err := s.DB.Create(role).Error; err != nil {
		return nil, err
	}

	// 分配权限
	if len(permissionIDs) > 0 {
		var perms []models.Permission
		s.DB.Where("id IN ?", permissionIDs).Find(&perms)
		s.DB.Model(role).Association("Permissions").Replace(perms)
	}

	// 重新加载权限
	s.DB.Preload("Permissions").First(role, role.ID)

	return role, nil
}

// UpdateRole 更新角色
func (s *RBACService) UpdateRole(roleID uint, displayName, description string, permissionIDs []uint) (*models.Role, error) {
	var role models.Role
	if err := s.DB.First(&role, roleID).Error; err != nil {
		return nil, err
	}

	if displayName != "" {
		role.DisplayName = displayName
	}
	if description != "" {
		role.Description = description
	}
	s.DB.Save(&role)

	// 更新权限
	if permissionIDs != nil {
		var perms []models.Permission
		s.DB.Where("id IN ?", permissionIDs).Find(&perms)
		s.DB.Model(&role).Association("Permissions").Replace(perms)
	}

	s.DB.Preload("Permissions").First(&role, roleID)
	return &role, nil
}

// DeleteRole 删除角色
func (s *RBACService) DeleteRole(roleID uint) error {
	return s.DB.Delete(&models.Role{}, roleID).Error
}

// ========== 权限管理 ==========

// GetPermissions 获取所有权限
func (s *RBACService) GetPermissions() ([]models.Permission, error) {
	var perms []models.Permission
	if err := s.DB.Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// ========== 用户角色管理 ==========

// GetUserRoles 获取用户的所有角色
func (s *RBACService) GetUserRoles(userID uint) ([]models.Role, error) {
	return models.GetUserRoles(s.DB, userID)
}

// AssignRoleToUser 为用户分配角色
func (s *RBACService) AssignRoleToUser(userID, roleID uint) error {
	// 检查是否已分配
	var count int64
	s.DB.Model(&models.UserRole{}).Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count)
	if count > 0 {
		return errors.New("用户已有该角色")
	}

	userRole := &models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return s.DB.Create(userRole).Error
}

// RemoveRoleFromUser 移除用户角色
func (s *RBACService) RemoveRoleFromUser(userID, roleID uint) error {
	return s.DB.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}

// ========== 权限检查 ==========

// CheckPermission 检查用户是否有某个权限
func (s *RBACService) CheckPermission(userID uint, permissionName string) (bool, error) {
	return models.HasUserPermission(s.DB, userID, permissionName)
}

// CheckAnyPermission 检查用户是否有任意一个权限
func (s *RBACService) CheckAnyPermission(userID uint, permissionNames []string) (bool, error) {
	for _, permName := range permissionNames {
		has, err := s.CheckPermission(userID, permName)
		if err != nil {
			return false, err
		}
		if has {
			return true, nil
		}
	}
	return false, nil
}

// GetUserPermissions 获取用户的所有权限
func (s *RBACService) GetUserPermissions(userID uint) ([]models.Permission, error) {
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	permMap := make(map[uint]models.Permission)
	for _, role := range roles {
		for _, perm := range role.Permissions {
			permMap[perm.ID] = perm
		}
	}

	perms := make([]models.Permission, 0, len(permMap))
	for _, perm := range permMap {
		perms = append(perms, perm)
	}

	return perms, nil
}

// ========== 初始化 ==========

// Initialize 初始化 RBAC（创建默认角色和权限）
func (s *RBACService) Initialize() error {
	return models.InitializeRBAC(s.DB)
}
