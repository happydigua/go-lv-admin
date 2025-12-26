package middleware

import (
	"bytes"
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 自定义 ResponseWriter 用于捕获响应
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// OperationLog 操作日志中间件
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要记录的路径
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/health") {
			c.Next()
			return
		}

		// 记录开始时间
		start := time.Now()

		// 获取请求体
		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// 包装 ResponseWriter 以捕获响应
		blw := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start).Milliseconds()

		// 获取用户信息
		userId, _ := c.Get("userId")
		username, _ := c.Get("username")

		// 解析模块和操作类型
		module, action := parseModuleAction(c.Request.Method, path)

		// 限制 body 和 response 长度
		bodyStr := string(body)
		if len(bodyStr) > 2000 {
			bodyStr = bodyStr[:2000] + "..."
		}
		respStr := blw.body.String()
		if len(respStr) > 2000 {
			respStr = respStr[:2000] + "..."
		}

		// 创建日志记录
		log := model.LvOperationLog{
			UserId:    toUint(userId),
			Username:  toString(username),
			Ip:        c.ClientIP(),
			Method:    c.Request.Method,
			Path:      path,
			Status:    c.Writer.Status(),
			Latency:   latency,
			UserAgent: c.Request.UserAgent(),
			Body:      bodyStr,
			Response:  respStr,
			Module:    module,
			Action:    action,
		}

		// 异步写入日志
		go func() {
			if err := global.LV_DB.Create(&log).Error; err != nil {
				global.LV_LOG.Error("操作日志写入失败: " + err.Error())
			}
		}()
	}
}

// 解析模块和操作类型
func parseModuleAction(method, path string) (module, action string) {
	// 解析模块
	if strings.Contains(path, "/system/user") {
		module = "用户管理"
	} else if strings.Contains(path, "/system/role") {
		module = "角色管理"
	} else if strings.Contains(path, "/system/menu") {
		module = "菜单管理"
	} else if strings.Contains(path, "/dashboard") {
		module = "仪表盘"
	} else if strings.Contains(path, "/profile") {
		module = "个人中心"
	} else if strings.Contains(path, "/base/login") {
		module = "登录"
	} else {
		module = "其他"
	}

	// 解析操作类型
	switch method {
	case "GET":
		action = "查询"
	case "POST":
		if strings.Contains(path, "/login") {
			action = "登录"
		} else {
			action = "新增"
		}
	case "PUT":
		if strings.Contains(path, "reset-password") {
			action = "重置密码"
		} else if strings.Contains(path, "password") {
			action = "修改密码"
		} else {
			action = "修改"
		}
	case "DELETE":
		action = "删除"
	default:
		action = method
	}

	return
}

func toUint(v interface{}) uint {
	if v == nil {
		return 0
	}
	if u, ok := v.(uint); ok {
		return u
	}
	return 0
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
