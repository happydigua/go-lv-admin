package v1

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"go-lv-vue-admin/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemRoleApi struct{}

var systemRoleService = service.SystemRoleService{}

// GetRoleList
// @Summary 获取角色列表
// @Router /system/role/list [get]
func (s *SystemRoleApi) GetRoleList(c *gin.Context) {
	roles, err := systemRoleService.GetRoleList()
	if err != nil {
		global.LV_LOG.Error("获取角色列表失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "获取角色列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": roles,
		"msg":  "success",
	})
}

// CreateRole
// @Summary 创建角色
// @Router /system/role [post]
func (s *SystemRoleApi) CreateRole(c *gin.Context) {
	var role model.LvRole
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	if err := systemRoleService.CreateRole(&role); err != nil {
		global.LV_LOG.Error("创建角色失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "创建成功"})
}

// UpdateRole
// @Summary 更新角色
// @Router /system/role/:id [put]
func (s *SystemRoleApi) UpdateRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var role model.LvRole
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}
	role.ID = uint(id)

	if err := systemRoleService.UpdateRole(&role); err != nil {
		global.LV_LOG.Error("更新角色失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "更新角色失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "更新成功"})
}

// DeleteRole
// @Summary 删除角色
// @Router /system/role/:id [delete]
func (s *SystemRoleApi) DeleteRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := systemRoleService.DeleteRole(uint(id)); err != nil {
		global.LV_LOG.Error("删除角色失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "删除成功"})
}
