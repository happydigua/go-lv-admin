package v1

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"go-lv-vue-admin/internal/model/request"
	"go-lv-vue-admin/internal/model/response"
	"go-lv-vue-admin/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserApi struct{}

// Login
// @Tags Base
// @Summary Login
// @accept application/json
// @Produce application/json
// @Param data body request.Login true "Username, Password"
// @Success 200 {object} response.Response{data=response.LoginResponse,msg=string}
// @Router /base/login [post]
func (b *UserApi) Login(c *gin.Context) {
	var l request.Login
	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	// TODO: Verify Captcha

	u := &model.LvUser{Username: l.Username, Password: l.Password}
	userService := service.UserService{}

	user, err := userService.Login(u)
	if err != nil {
		global.LV_LOG.Error("login failed", zap.Error(err))
		c.JSON(400, gin.H{"code": 7, "msg": "用户名或密码错误"})
		return
	}

	if user.Status != 1 {
		c.JSON(400, gin.H{"code": 7, "msg": "用户被冻结"})
		return
	}

	// Generate Token
	token, expiresAt, err := userService.CreateToken(*user)
	if err != nil {
		global.LV_LOG.Error("get token failed", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "获取Token失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": response.LoginResponse{
			User:      *user,
			Token:     token,
			ExpiresAt: expiresAt,
		},
		"msg": "登录成功",
	})
}
