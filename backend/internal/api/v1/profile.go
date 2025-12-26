package v1

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProfileApi struct{}

// GetProfile 获取当前用户信息
// @Router /profile [get]
func (p *ProfileApi) GetProfile(c *gin.Context) {
	// 从 context 获取用户 ID（实际应从 JWT 中解析）
	userIdStr := c.GetHeader("X-User-Id")
	if userIdStr == "" {
		userIdStr = "1" // 临时默认值，后续通过 JWT 中间件获取
	}
	userId, _ := strconv.Atoi(userIdStr)

	var user model.LvUser
	if err := global.LV_DB.Preload("Role").First(&user, userId).Error; err != nil {
		c.JSON(500, gin.H{"code": 7, "msg": "获取用户信息失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": user,
		"msg":  "success",
	})
}

// UpdateProfile 更新个人资料
// @Router /profile [put]
func (p *ProfileApi) UpdateProfile(c *gin.Context) {
	userIdStr := c.GetHeader("X-User-Id")
	if userIdStr == "" {
		userIdStr = "1"
	}
	userId, _ := strconv.Atoi(userIdStr)

	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	if err := global.LV_DB.Model(&model.LvUser{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"nickname": req.Nickname,
		"email":    req.Email,
		"phone":    req.Phone,
		"avatar":   req.Avatar,
	}).Error; err != nil {
		global.LV_LOG.Error("更新个人资料失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "更新失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "更新成功"})
}

// ChangePassword 修改密码
// @Router /profile/password [put]
func (p *ProfileApi) ChangePassword(c *gin.Context) {
	userIdStr := c.GetHeader("X-User-Id")
	if userIdStr == "" {
		userIdStr = "1"
	}
	userId, _ := strconv.Atoi(userIdStr)

	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": "请填写完整的密码信息"})
		return
	}

	// 验证旧密码
	var user model.LvUser
	if err := global.LV_DB.First(&user, userId).Error; err != nil {
		c.JSON(500, gin.H{"code": 7, "msg": "用户不存在"})
		return
	}

	// TODO: 使用 bcrypt 比较密码
	if user.Password != req.OldPassword {
		c.JSON(400, gin.H{"code": 7, "msg": "原密码错误"})
		return
	}

	// 更新密码
	if err := global.LV_DB.Model(&model.LvUser{}).Where("id = ?", userId).Update("password", req.NewPassword).Error; err != nil {
		global.LV_LOG.Error("修改密码失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "修改密码失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "密码修改成功"})
}
