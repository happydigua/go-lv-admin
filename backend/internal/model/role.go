package model

import (
	"gorm.io/gorm"
)

type LvRole struct {
	gorm.Model
	Name    string   `json:"name" gorm:"comment:角色名"`
	Keyword string   `json:"keyword" gorm:"unique;comment:角色关键字"`
	Desc    string   `json:"desc" gorm:"comment:角色说明"`
	Status  int      `json:"status" gorm:"default:1;comment:角色状态"`
	Sort    int      `json:"sort" gorm:"default:0;comment:角色排序"`
	Menus   []LvMenu `json:"menus" gorm:"many2many:lv_role_menus;"`
}

func (LvRole) TableName() string {
	return "lv_roles"
}
