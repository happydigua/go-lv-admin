package middleware

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"code": 401, "msg": "未登录或非法访问"})
			c.Abort()
			return
		}

		// Bearer token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(401, gin.H{"code": 401, "msg": "请求头中 Authorization 格式有误"})
			c.Abort()
			return
		}

		// 解析 token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			global.LV_LOG.Error("token 解析失败: " + err.Error())
			c.JSON(401, gin.H{"code": 401, "msg": "Token 已过期或无效"})
			c.Abort()
			return
		}

		// 将用户信息存入 context
		c.Set("userId", claims.UserId)
		c.Set("username", claims.Username)
		c.Set("roleId", claims.RoleId)
		c.Set("claims", claims)

		c.Next()
	}
}
