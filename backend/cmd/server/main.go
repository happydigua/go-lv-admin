package main

import (
	"fmt"
	"go-lv-vue-admin/internal/core"
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. Initialize Configuration
	global.LV_VP = core.Viper()

	// 2. Initialize Logger
	global.LV_LOG = core.Zap()
	zap.ReplaceGlobals(global.LV_LOG)

	// 3. Initialize Database
	global.LV_DB = core.Gorm()
	if global.LV_DB != nil {
		// Initialize tables
		core.RegisterTables()
		// Initialize Casbin
		global.LV_ENFORCER = core.InitCasbin()
		db, _ := global.LV_DB.DB()
		defer db.Close()
	} else {
		global.LV_LOG.Error("Database connection failed")
		// For now, don't exit to allow running without DB
	}

	// 4. Initialize Router
	gin.SetMode(global.LV_CONFIG.Server.Mode)
	r := gin.Default()
	router.InitRouter(r)

	addr := fmt.Sprintf(":%d", global.LV_CONFIG.Server.Port)
	global.LV_LOG.Info("Server exiting", zap.String("addr", addr))

	if err := r.Run(addr); err != nil {
		global.LV_LOG.Error(err.Error())
	}
}
