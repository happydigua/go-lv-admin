package v1

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/service"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GeneratorApi struct{}

var generatorService = service.GeneratorService{}

// GetTables 获取数据库表列表
func (a *GeneratorApi) GetTables(c *gin.Context) {
	tables, err := generatorService.GetTables()
	if err != nil {
		global.LV_LOG.Error("获取表列表失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "获取失败"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": tables, "msg": "success"})
}

// GetTableColumns 获取表的列信息
func (a *GeneratorApi) GetTableColumns(c *gin.Context) {
	tableName := c.Query("tableName")
	if tableName == "" {
		c.JSON(400, gin.H{"code": 7, "msg": "表名不能为空"})
		return
	}

	columns, err := generatorService.GetTableColumns(tableName)
	if err != nil {
		global.LV_LOG.Error("获取列信息失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "获取失败"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": columns, "msg": "success"})
}

// GenerateCode 生成代码并写入文件
func (a *GeneratorApi) GenerateCode(c *gin.Context) {
	var req service.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	// 自动检测表是否有 deleted_at 字段
	req.HasDeletedAt = generatorService.HasDeletedAtColumn(req.TableName)

	// 获取项目根路径（假设在 backend 目录运行）
	backendPath, _ := filepath.Abs(".")
	frontendPath, _ := filepath.Abs("../frontend")

	result, err := generatorService.WriteGeneratedFiles(req, backendPath, frontendPath)
	if err != nil {
		global.LV_LOG.Error("生成代码失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "生成失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "data": result, "msg": "success"})
}

// PreviewCode 预览生成的代码
func (a *GeneratorApi) PreviewCode(c *gin.Context) {
	var config service.GenerateConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	// 自动检测表是否有 deleted_at 字段
	config.HasDeletedAt = generatorService.HasDeletedAtColumn(config.TableName)

	codes, err := generatorService.GenerateCode(config)
	if err != nil {
		global.LV_LOG.Error("预览代码失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "预览失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "data": codes, "msg": "success"})
}
