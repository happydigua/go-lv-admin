package v1

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"go-lv-vue-admin/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DemoApi struct{}

var demoService = service.DemoService{}

// GetDemoList 获取列表
func (a *DemoApi) GetDemoList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	sortField := c.Query("sortField")
	sortOrder := c.Query("sortOrder")

	list, total, err := demoService.GetDemoList(page, pageSize, filters, sortField, sortOrder)
	if err != nil {
		global.LV_LOG.Error("获取列表失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "获取列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  list,
			"total": total,
		},
		"msg": "success",
	})
}

// CreateDemo 创建
func (a *DemoApi) CreateDemo(c *gin.Context) {
	var demo model.LvDemo
	if err := c.ShouldBindJSON(&demo); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": "参数错误"})
		return
	}

	if err := demoService.CreateDemo(&demo); err != nil {
		global.LV_LOG.Error("创建失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "创建失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "创建成功"})
}

// UpdateDemo 更新
func (a *DemoApi) UpdateDemo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var demo model.LvDemo
	if err := c.ShouldBindJSON(&demo); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": "参数错误"})
		return
	}
	demo.ID = uint(id)

	if err := demoService.UpdateDemo(&demo); err != nil {
		global.LV_LOG.Error("更新失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "更新失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "更新成功"})
}

// DeleteDemo 删除
func (a *DemoApi) DeleteDemo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := demoService.DeleteDemo(uint(id)); err != nil {
		global.LV_LOG.Error("删除失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "删除失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "删除成功"})
}
