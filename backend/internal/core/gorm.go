package core

import (
	"go-lv-vue-admin/internal/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Gorm() *gorm.DB {
	m := global.LV_CONFIG.Database
	if m.Source == "" {
		return nil
	}
	dsn := m.Source
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig()); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		return db
	}
}

func gormConfig() *gorm.Config {
	config := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	logMode := global.LV_CONFIG.Database.LogMode
	switch logMode {
	case "silent":
		config.Logger = logger.Default.LogMode(logger.Silent)
	case "error":
		config.Logger = logger.Default.LogMode(logger.Error)
	case "warn":
		config.Logger = logger.Default.LogMode(logger.Warn)
	case "info":
		config.Logger = logger.Default.LogMode(logger.Info)
	default:
		config.Logger = logger.Default.LogMode(logger.Info)
	}
	return config
}
