package v1

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SettingApi struct{}

var settingService = service.SettingService{}

// GetSettings 获取所有设置
func (s *SettingApi) GetSettings(c *gin.Context) {
	settings, err := settingService.GetAllSettings()
	if err != nil {
		global.LV_LOG.Error("获取设置失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "获取设置失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "data": settings, "msg": "success"})
}

// GetPublicSettings 获取公开设置（无需登录）
func (s *SettingApi) GetPublicSettings(c *gin.Context) {
	settings, err := settingService.GetPublicSettings()
	if err != nil {
		c.JSON(500, gin.H{"code": 7, "msg": "获取设置失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "data": settings, "msg": "success"})
}

// UpdateSettings 批量更新设置
func (s *SettingApi) UpdateSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": "参数错误"})
		return
	}

	if err := settingService.BatchUpdateSettings(req); err != nil {
		global.LV_LOG.Error("更新设置失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "更新设置失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "保存成功"})
}
