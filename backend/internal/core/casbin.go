package core

import (
	"go-lv-vue-admin/internal/global"
	"os"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
)

func InitCasbin() *casbin.Enforcer {
	// 获取当前工作目录
	wd, _ := os.Getwd()
	modelPath := filepath.Join(wd, "config", "rbac_model.conf")

	// 使用 GORM 适配器，将策略存储到数据库
	adapter, err := gormadapter.NewAdapterByDB(global.LV_DB)
	if err != nil {
		global.LV_LOG.Error("Casbin adapter init failed", zap.Error(err))
		return nil
	}

	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		global.LV_LOG.Error("Casbin enforcer init failed", zap.Error(err))
		return nil
	}

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		global.LV_LOG.Error("Casbin load policy failed", zap.Error(err))
		return nil
	}

	global.LV_LOG.Info("Casbin init success")
	return enforcer
}
