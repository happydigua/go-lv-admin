package v1

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"go-lv-vue-admin/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemUserApi struct{}

var systemUserService = service.SystemUserService{}

// GetUserList
// @Summary 获取用户列表
// @Router /system/user/list [get]
func (s *SystemUserApi) GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	username := c.Query("username")
	phone := c.Query("phone")

	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		statusVal, _ := strconv.Atoi(statusStr)
		status = &statusVal
	}

	users, total, err := systemUserService.GetUserList(page, pageSize, username, phone, status)
	if err != nil {
		global.LV_LOG.Error("获取用户列表失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "获取用户列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"list":     users,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
		"msg": "success",
	})
}

// CreateUser
// @Summary 创建用户
// @Router /system/user [post]
func (s *SystemUserApi) CreateUser(c *gin.Context) {
	var user model.LvUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	if err := systemUserService.CreateUser(&user); err != nil {
		global.LV_LOG.Error("创建用户失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "创建成功"})
}

// UpdateUser
// @Summary 更新用户
// @Router /system/user/:id [put]
func (s *SystemUserApi) UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user model.LvUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}
	user.ID = uint(id)

	if err := systemUserService.UpdateUser(&user); err != nil {
		global.LV_LOG.Error("更新用户失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "更新用户失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "更新成功"})
}

// DeleteUser
// @Summary 删除用户
// @Router /system/user/:id [delete]
func (s *SystemUserApi) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := systemUserService.DeleteUser(uint(id)); err != nil {
		global.LV_LOG.Error("删除用户失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "删除成功"})
}

// ResetPassword
// @Summary 重置用户密码
// @Router /system/user/:id/reset-password [put]
func (s *SystemUserApi) ResetPassword(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	newPassword := req.Password
	if newPassword == "" {
		newPassword = "123456" // 默认密码
	}

	if err := systemUserService.ResetPassword(uint(id), newPassword); err != nil {
		global.LV_LOG.Error("重置密码失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "密码已重置为: " + newPassword})
}

// GetRoleOptions
// @Summary 获取角色选项
// @Router /system/user/role-options [get]
func (s *SystemUserApi) GetRoleOptions(c *gin.Context) {
	roles, err := systemUserService.GetRoleList()
	if err != nil {
		c.JSON(500, gin.H{"code": 7, "msg": "获取角色列表失败"})
		return
	}

	options := make([]gin.H, len(roles))
	for i, role := range roles {
		options[i] = gin.H{"label": role.Name, "value": role.ID}
	}

	c.JSON(200, gin.H{"code": 0, "data": options, "msg": "success"})
}
