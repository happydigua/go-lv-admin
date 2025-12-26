package service

import (
	"errors"
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
)

type SystemRoleService struct{}

// GetRoleList 获取角色列表
func (s *SystemRoleService) GetRoleList() ([]model.LvRole, error) {
	var roles []model.LvRole
	err := global.LV_DB.Order("sort ASC").Find(&roles).Error
	return roles, err
}

// CreateRole 创建角色
func (s *SystemRoleService) CreateRole(role *model.LvRole) error {
	// 检查角色标识是否已存在
	var count int64
	global.LV_DB.Model(&model.LvRole{}).Where("keyword = ?", role.Keyword).Count(&count)
	if count > 0 {
		return errors.New("角色标识已存在")
	}
	return global.LV_DB.Create(role).Error
}

// UpdateRole 更新角色
func (s *SystemRoleService) UpdateRole(role *model.LvRole) error {
	return global.LV_DB.Model(&model.LvRole{}).Where("id = ?", role.ID).Updates(map[string]interface{}{
		"name":   role.Name,
		"desc":   role.Desc,
		"status": role.Status,
		"sort":   role.Sort,
	}).Error
}

// DeleteRole 删除角色
func (s *SystemRoleService) DeleteRole(id uint) error {
	// 不允许删除 ID 为 1 的管理员角色
	if id == 1 {
		return errors.New("不能删除管理员角色")
	}

	// 检查是否有用户使用该角色
	var count int64
	global.LV_DB.Model(&model.LvUser{}).Where("role_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该角色下还有用户，无法删除")
	}

	return global.LV_DB.Delete(&model.LvRole{}, id).Error
}
