package model

import (
	"gorm.io/gorm"
)

// LvSetting 系统设置表
type LvSetting struct {
	gorm.Model
	Key         string `json:"key" gorm:"uniqueIndex;comment:设置键"`
	Value       string `json:"value" gorm:"type:text;comment:设置值"`
	Name        string `json:"name" gorm:"comment:显示名称"`
	Description string `json:"description" gorm:"comment:描述"`
}

func (LvSetting) TableName() string {
	return "lv_settings"
}

// 默认设置（仅基础信息）
var DefaultSettings = []LvSetting{
	{Key: "site_name", Value: "Go Lv Admin", Name: "系统名称", Description: "显示在标题栏和登录页"},
	{Key: "site_logo", Value: "", Name: "系统Logo", Description: "Logo图片URL"},
	{Key: "site_footer", Value: "© 2024 Go Lv Admin", Name: "底部版权", Description: "页面底部显示的版权信息"},
}
