package v1

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"go-lv-vue-admin/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemMenuApi struct{}

var systemMenuService = service.SystemMenuService{}

// GetMenuList
// @Summary 获取菜单列表
// @Router /system/menu/list [get]
func (s *SystemMenuApi) GetMenuList(c *gin.Context) {
	menus, err := systemMenuService.GetMenuList()
	if err != nil {
		global.LV_LOG.Error("获取菜单列表失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "获取菜单列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": menus,
		"msg":  "success",
	})
}

// CreateMenu
// @Summary 创建菜单
// @Router /system/menu [post]
func (s *SystemMenuApi) CreateMenu(c *gin.Context) {
	var menu model.LvMenu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	if err := systemMenuService.CreateMenu(&menu); err != nil {
		global.LV_LOG.Error("创建菜单失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "创建成功"})
}

// UpdateMenu
// @Summary 更新菜单
// @Router /system/menu/:id [put]
func (s *SystemMenuApi) UpdateMenu(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var menu model.LvMenu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}
	menu.ID = uint(id)

	if err := systemMenuService.UpdateMenu(&menu); err != nil {
		global.LV_LOG.Error("更新菜单失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "更新菜单失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "更新成功"})
}

// DeleteMenu
// @Summary 删除菜单
// @Router /system/menu/:id [delete]
func (s *SystemMenuApi) DeleteMenu(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := systemMenuService.DeleteMenu(uint(id)); err != nil {
		global.LV_LOG.Error("删除菜单失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "删除成功"})
}
