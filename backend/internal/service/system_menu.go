package service

import (
	"errors"
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
)

type SystemMenuService struct{}

// GetMenuList 获取菜单列表（树形）
func (s *SystemMenuService) GetMenuList() ([]model.LvMenu, error) {
	var menus []model.LvMenu
	err := global.LV_DB.Order("sort ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus, 0), nil
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []model.LvMenu, parentId uint) []model.LvMenu {
	var tree []model.LvMenu
	for _, menu := range menus {
		if menu.ParentId == parentId {
			children := buildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				menu.Children = children
			}
			tree = append(tree, menu)
		}
	}
	return tree
}

// CreateMenu 创建菜单
func (s *SystemMenuService) CreateMenu(menu *model.LvMenu) error {
	return global.LV_DB.Create(menu).Error
}

// UpdateMenu 更新菜单
func (s *SystemMenuService) UpdateMenu(menu *model.LvMenu) error {
	return global.LV_DB.Model(&model.LvMenu{}).Where("id = ?", menu.ID).Updates(map[string]interface{}{
		"parent_id":  menu.ParentId,
		"title":      menu.Title,
		"path":       menu.Path,
		"name":       menu.Name,
		"component":  menu.Component,
		"icon":       menu.Icon,
		"sort":       menu.Sort,
		"type":       menu.Type,
		"permission": menu.Permission,
		"hidden":     menu.Hidden,
		"keep_alive": menu.KeepAlive,
	}).Error
}

// DeleteMenu 删除菜单
func (s *SystemMenuService) DeleteMenu(id uint) error {
	// 检查是否有子菜单
	var count int64
	global.LV_DB.Model(&model.LvMenu{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("请先删除子菜单")
	}
	return global.LV_DB.Delete(&model.LvMenu{}, id).Error
}
