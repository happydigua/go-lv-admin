package v1

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"go-lv-vue-admin/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PermissionApi struct{}

var permissionService = service.PermissionService{}

// GetRoleMenus
// @Summary 获取角色的菜单权限
// @Router /system/role/:id/menus [get]
func (p *PermissionApi) GetRoleMenus(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": "参数错误"})
		return
	}

	menuIds, err := permissionService.GetRoleMenus(uint(roleId))
	if err != nil {
		global.LV_LOG.Error("获取角色菜单失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "获取失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "data": menuIds, "msg": "success"})
}

// SetRoleMenus
// @Summary 设置角色的菜单权限
// @Router /system/role/:id/menus [put]
func (p *PermissionApi) SetRoleMenus(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": "参数错误"})
		return
	}

	var req struct {
		MenuIds []uint `json:"menuIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	if err := permissionService.SetRoleMenus(uint(roleId), req.MenuIds); err != nil {
		global.LV_LOG.Error("设置角色菜单失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "设置失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "设置成功"})
}

// GetUserPermissions
// @Summary 获取当前用户的按钮权限
// @Router /user/permissions [get]
func (p *PermissionApi) GetUserPermissions(c *gin.Context) {
	roleId, exists := c.Get("roleId")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "msg": "未登录"})
		return
	}

	permissions, err := permissionService.GetUserPermissions(roleId.(uint))
	if err != nil {
		global.LV_LOG.Error("获取用户权限失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "获取失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "data": permissions, "msg": "success"})
}

// GetUserMenus
// @Summary 获取当前用户可访问的菜单
// @Router /user/menus [get]
func (p *PermissionApi) GetUserMenus(c *gin.Context) {
	roleId, exists := c.Get("roleId")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "msg": "未登录"})
		return
	}

	menus, err := permissionService.GetUserMenus(roleId.(uint))
	if err != nil {
		global.LV_LOG.Error("获取用户菜单失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "获取失败"})
		return
	}

	// 构建菜单树
	menuTree := buildMenuTree(menus, 0)
	c.JSON(200, gin.H{"code": 0, "data": menuTree, "msg": "success"})
}

// buildMenuTree 递归构建菜单树
func buildMenuTree(menus []model.LvMenu, parentId uint) []model.LvMenu {
	var tree []model.LvMenu
	for _, menu := range menus {
		if menu.ParentId == parentId {
			children := buildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				menu.Children = children
			}
			tree = append(tree, menu)
		}
	}
	return tree
}
