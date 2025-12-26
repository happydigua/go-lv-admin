package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-lv-vue-admin/internal/config"
)

var (
	LV_DB       *gorm.DB
	LV_CONFIG   config.Config
	LV_VP       *viper.Viper
	LV_LOG      *zap.Logger
	LV_ENFORCER *casbin.Enforcer
)
