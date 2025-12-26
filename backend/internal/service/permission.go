package service

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
)

type PermissionService struct{}

// GetRoleMenus 获取角色已分配的菜单ID列表
func (s *PermissionService) GetRoleMenus(roleId uint) ([]uint, error) {
	var role model.LvRole
	err := global.LV_DB.Preload("Menus").First(&role, roleId).Error
	if err != nil {
		return nil, err
	}

	var menuIds []uint
	for _, menu := range role.Menus {
		menuIds = append(menuIds, menu.ID)
	}
	return menuIds, nil
}

// SetRoleMenus 设置角色的菜单权限
func (s *PermissionService) SetRoleMenus(roleId uint, menuIds []uint) error {
	var role model.LvRole
	if err := global.LV_DB.First(&role, roleId).Error; err != nil {
		return err
	}

	// 获取菜单列表
	var menus []model.LvMenu
	if len(menuIds) > 0 {
		global.LV_DB.Find(&menus, menuIds)
	}

	// 更新角色的菜单关联
	err := global.LV_DB.Model(&role).Association("Menus").Replace(menus)
	if err != nil {
		return err
	}

	// 更新 Casbin 策略
	s.updateCasbinPolicy(role.Keyword, menus)

	return nil
}

// updateCasbinPolicy 更新 Casbin 权限策略
func (s *PermissionService) updateCasbinPolicy(roleKeyword string, menus []model.LvMenu) {
	if global.LV_ENFORCER == nil {
		return
	}

	// 删除旧策略
	global.LV_ENFORCER.DeletePermissionsForUser(roleKeyword)

	// 添加新策略
	for _, menu := range menus {
		if menu.Path != "" {
			// 添加路由权限
			global.LV_ENFORCER.AddPolicy(roleKeyword, menu.Path, "GET")
			global.LV_ENFORCER.AddPolicy(roleKeyword, menu.Path, "POST")
			global.LV_ENFORCER.AddPolicy(roleKeyword, menu.Path, "PUT")
			global.LV_ENFORCER.AddPolicy(roleKeyword, menu.Path, "DELETE")
		}
		if menu.Permission != "" {
			// 添加按钮权限
			global.LV_ENFORCER.AddPolicy(roleKeyword, menu.Permission, "btn")
		}
	}

	// 保存策略
	global.LV_ENFORCER.SavePolicy()
}

// GetUserPermissions 获取用户的按钮权限列表
func (s *PermissionService) GetUserPermissions(roleId uint) ([]string, error) {
	var role model.LvRole
	err := global.LV_DB.Preload("Menus").First(&role, roleId).Error
	if err != nil {
		return nil, err
	}

	// admin 角色拥有所有权限
	if role.Keyword == "admin" {
		return []string{"*"}, nil
	}

	var permissions []string
	for _, menu := range role.Menus {
		if menu.Permission != "" {
			permissions = append(permissions, menu.Permission)
		}
	}
	return permissions, nil
}

// GetUserMenus 获取用户可访问的菜单（根据角色）
func (s *PermissionService) GetUserMenus(roleId uint) ([]model.LvMenu, error) {
	var role model.LvRole
	err := global.LV_DB.Preload("Menus").First(&role, roleId).Error
	if err != nil {
		return nil, err
	}

	// admin 角色返回所有菜单
	if role.Keyword == "admin" {
		var allMenus []model.LvMenu
		global.LV_DB.Order("sort ASC").Find(&allMenus)
		return allMenus, nil
	}

	return role.Menus, nil
}
