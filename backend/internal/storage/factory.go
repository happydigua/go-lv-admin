package storage

import (
	"fmt"

	"go-lv-vue-admin/internal/global"
)

// driver 全局存储驱动实例
var driver StorageDriver

// InitStorage 初始化存储驱动
func InitStorage() error {
	cfg := global.LV_CONFIG.Storage

	switch cfg.Driver {
	case "local", "":
		driver = NewLocalDriver(cfg.Local)
		global.LV_LOG.Info("使用本地存储驱动")
	case "oss":
		d, err := NewOSSDriver(cfg.OSS)
		if err != nil {
			return fmt.Errorf("初始化OSS驱动失败: %w", err)
		}
		driver = d
		global.LV_LOG.Info("使用阿里云OSS存储驱动")
	case "cos":
		d, err := NewCOSDriver(cfg.COS)
		if err != nil {
			return fmt.Errorf("初始化COS驱动失败: %w", err)
		}
		driver = d
		global.LV_LOG.Info("使用腾讯云COS存储驱动")
	case "r2":
		d, err := NewR2Driver(cfg.R2)
		if err != nil {
			return fmt.Errorf("初始化R2驱动失败: %w", err)
		}
		driver = d
		global.LV_LOG.Info("使用Cloudflare R2存储驱动")
	default:
		return fmt.Errorf("不支持的存储驱动: %s", cfg.Driver)
	}

	return nil
}

// GetDriver 获取存储驱动
func GetDriver() StorageDriver {
	if driver == nil {
		// 默认使用本地存储
		driver = NewLocalDriver(global.LV_CONFIG.Storage.Local)
	}
	return driver
}
