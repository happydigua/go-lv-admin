package v1

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OperationLogApi struct{}

var operationLogService = service.OperationLogService{}

// GetOperationLogList
// @Summary 获取操作日志列表
// @Router /system/log/list [get]
func (o *OperationLogApi) GetOperationLogList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	username := c.Query("username")
	module := c.Query("module")
	action := c.Query("action")

	logs, total, err := operationLogService.GetOperationLogList(page, pageSize, username, module, action)
	if err != nil {
		global.LV_LOG.Error("获取操作日志列表失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "获取操作日志列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"list":     logs,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
		"msg": "success",
	})
}

// DeleteOperationLogs
// @Summary 批量删除操作日志
// @Router /system/log [delete]
func (o *OperationLogApi) DeleteOperationLogs(c *gin.Context) {
	var req struct {
		Ids []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	if err := operationLogService.DeleteOperationLogs(req.Ids); err != nil {
		global.LV_LOG.Error("删除操作日志失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "删除失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "删除成功"})
}

// ClearOperationLogs
// @Summary 清空操作日志
// @Router /system/log/clear [delete]
func (o *OperationLogApi) ClearOperationLogs(c *gin.Context) {
	if err := operationLogService.ClearOperationLogs(); err != nil {
		global.LV_LOG.Error("清空操作日志失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "清空失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "清空成功"})
}
