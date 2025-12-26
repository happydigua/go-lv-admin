package model

import (
	"gorm.io/gorm"
)

type LvMenu struct {
	gorm.Model
	ParentId   uint     `json:"parentId" gorm:"comment:父菜单ID"`
	Title      string   `json:"title" gorm:"comment:菜单名"`
	Path       string   `json:"path" gorm:"comment:路由地址"`
	Name       string   `json:"name" gorm:"comment:路由名称"`
	Component  string   `json:"component" gorm:"comment:组件路径"`
	Icon       string   `json:"icon" gorm:"comment:菜单图标"`
	Sort       int      `json:"sort" gorm:"default:0;comment:排序"`
	Type       int      `json:"type" gorm:"default:1;comment:类型 1目录 2菜单 3按钮"`
	Permission string   `json:"permission" gorm:"comment:权限标识"`
	Hidden     bool     `json:"hidden" gorm:"default:false;comment:是否隐藏"`
	KeepAlive  bool     `json:"keepAlive" gorm:"default:true;comment:是否缓存"`
	Children   []LvMenu `json:"children" gorm:"-"`
}

func (LvMenu) TableName() string {
	return "lv_menus"
}
