package service

import (
	"errors"
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"go-lv-vue-admin/pkg/utils"
)

type SystemUserService struct{}

// GetUserList 获取用户列表
func (s *SystemUserService) GetUserList(page, pageSize int, username, phone string, status *int) ([]model.LvUser, int64, error) {
	var users []model.LvUser
	var total int64

	db := global.LV_DB.Model(&model.LvUser{})

	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if phone != "" {
		db = db.Where("phone LIKE ?", "%"+phone+"%")
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}

	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Preload("Role").Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}

// CreateUser 创建用户
func (s *SystemUserService) CreateUser(user *model.LvUser) error {
	// 检查用户名是否已存在
	var count int64
	global.LV_DB.Model(&model.LvUser{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 使用 bcrypt 加密密码
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return global.LV_DB.Create(user).Error
}

// UpdateUser 更新用户
func (s *SystemUserService) UpdateUser(user *model.LvUser) error {
	return global.LV_DB.Model(&model.LvUser{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"nickname": user.Nickname,
		"email":    user.Email,
		"phone":    user.Phone,
		"role_id":  user.RoleId,
		"status":   user.Status,
	}).Error
}

// DeleteUser 删除用户
func (s *SystemUserService) DeleteUser(id uint) error {
	// 不允许删除 ID 为 1 的超级管理员
	if id == 1 {
		return errors.New("不能删除超级管理员")
	}
	return global.LV_DB.Delete(&model.LvUser{}, id).Error
}

// ResetPassword 重置密码
func (s *SystemUserService) ResetPassword(id uint, newPassword string) error {
	// 不允许通过此接口修改超级管理员密码
	if id == 1 {
		return errors.New("不能通过此接口修改超级管理员密码")
	}
	// 使用 bcrypt 加密密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return global.LV_DB.Model(&model.LvUser{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

// GetRoleList 获取角色列表
func (s *SystemUserService) GetRoleList() ([]model.LvRole, error) {
	var roles []model.LvRole
	err := global.LV_DB.Find(&roles).Error
	return roles, err
}
