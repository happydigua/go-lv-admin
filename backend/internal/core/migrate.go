package core

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"go-lv-vue-admin/pkg/utils"
	"os"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RegisterTables() {
	db := global.LV_DB
	err := db.AutoMigrate(
		&model.LvUser{},
		&model.LvRole{},
		&model.LvMenu{},
		&model.LvOperationLog{},
		&model.LvSetting{},
	)
	if err != nil {
		global.LV_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.LV_LOG.Info("register table success")
	InitData(db)
	InitSettings(db)
}

func InitData(db *gorm.DB) {
	var adminRole model.LvRole
	// Check if admin role exists
	if err := db.Where("keyword = ?", "admin").First(&adminRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create Admin Role
			adminRole = model.LvRole{
				Name:    "管理员",
				Keyword: "admin",
				Desc:    "系统管理员",
				Status:  1,
				Sort:    1,
			}
			if err := db.Create(&adminRole).Error; err != nil {
				global.LV_LOG.Error("init admin role failed", zap.Error(err))
				return
			}
			global.LV_LOG.Info("init admin role success")
		} else {
			global.LV_LOG.Error("query admin role failed", zap.Error(err))
			return
		}
	}

	var adminUser model.LvUser
	// Check if admin user exists
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create Admin User with bcrypt hashed password
			hashedPassword, _ := utils.HashPassword("password")
			adminUser = model.LvUser{
				Username: "admin",
				Password: hashedPassword,
				Nickname: "超级管理员",
				RoleId:   adminRole.ID,
				Status:   1,
			}
			if err := db.Create(&adminUser).Error; err != nil {
				global.LV_LOG.Error("init admin user failed", zap.Error(err))
				return
			}
			global.LV_LOG.Info("init admin user success")
		}
	}

	// Initialize default menus
	var menuCount int64
	db.Model(&model.LvMenu{}).Count(&menuCount)
	if menuCount == 0 {
		initMenus(db)
	}
}

func initMenus(db *gorm.DB) {
	// 仪表盘
	dashboard := model.LvMenu{
		ParentId:  0,
		Title:     "仪表盘",
		Path:      "/dashboard",
		Name:      "Dashboard",
		Component: "views/dashboard/index",
		Icon:      "HomeOutline",
		Sort:      1,
		Type:      2,
	}
	db.Create(&dashboard)

	// 系统管理
	system := model.LvMenu{
		ParentId:  0,
		Title:     "系统管理",
		Path:      "/system",
		Name:      "System",
		Component: "",
		Icon:      "SettingsOutline",
		Sort:      2,
		Type:      1,
	}
	db.Create(&system)

	// 用户管理
	db.Create(&model.LvMenu{
		ParentId:  system.ID,
		Title:     "用户管理",
		Path:      "/system/user",
		Name:      "SystemUser",
		Component: "views/system/user/index",
		Icon:      "PersonOutline",
		Sort:      1,
		Type:      2,
	})

	// 角色管理
	db.Create(&model.LvMenu{
		ParentId:  system.ID,
		Title:     "角色管理",
		Path:      "/system/role",
		Name:      "SystemRole",
		Component: "views/system/role/index",
		Icon:      "PeopleOutline",
		Sort:      2,
		Type:      2,
	})

	// 菜单管理
	db.Create(&model.LvMenu{
		ParentId:  system.ID,
		Title:     "菜单管理",
		Path:      "/system/menu",
		Name:      "SystemMenu",
		Component: "views/system/menu/index",
		Icon:      "MenuOutline",
		Sort:      3,
		Type:      2,
	})

	global.LV_LOG.Info("init menus success")
}

// InitSettings 初始化默认设置
func InitSettings(db *gorm.DB) {
	for _, setting := range model.DefaultSettings {
		var count int64
		db.Model(&model.LvSetting{}).Where("`key` = ?", setting.Key).Count(&count)
		if count == 0 {
			db.Create(&setting)
		}
	}
	global.LV_LOG.Info("init settings success")
}
