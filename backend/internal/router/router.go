package router

import (
	v1 "go-lv-vue-admin/internal/api/v1"
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// 静态文件服务（上传的文件）
	r.Static("/uploads", "./uploads")

	// Public Group (无需认证)
	publicGroup := r.Group("")
	{
		// Health Check
		publicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})

		// Public Settings
		settingApi := v1.SettingApi{}
		publicGroup.GET("/settings/public", settingApi.GetPublicSettings)
	}

	// Base Router (Login, Register - 无需认证)
	baseApi := v1.UserApi{}
	baseGroup := r.Group("base")
	{
		baseGroup.POST("login", baseApi.Login)
	}

	// =========== 以下路由需要 JWT 认证 ===========
	privateGroup := r.Group("")
	privateGroup.Use(middleware.JWTAuth())
	privateGroup.Use(middleware.OperationLog()) // 操作日志中间件
	{
		// Settings (Private)
		settingApi := v1.SettingApi{}
		privateGroup.GET("/settings", settingApi.GetSettings)
		privateGroup.PUT("/settings", settingApi.UpdateSettings)

		// Dashboard Router
		dashboardApi := v1.DashboardApi{}
		privateGroup.GET("/dashboard/stats", dashboardApi.GetStats)
		privateGroup.GET("/dashboard/charts", dashboardApi.GetCharts)

		// Upload Router
		uploadApi := v1.UploadApi{}
		privateGroup.POST("/upload/image", uploadApi.UploadImage)
		privateGroup.POST("/upload/file", uploadApi.UploadFile)
		privateGroup.DELETE("/upload/file", uploadApi.DeleteFile)

		// System User Router
		systemUserApi := v1.SystemUserApi{}
		systemUserGroup := privateGroup.Group("system/user")
		{
			systemUserGroup.GET("list", systemUserApi.GetUserList)
			systemUserGroup.GET("role-options", systemUserApi.GetRoleOptions)
			systemUserGroup.POST("", systemUserApi.CreateUser)
			systemUserGroup.PUT(":id", systemUserApi.UpdateUser)
			systemUserGroup.DELETE(":id", systemUserApi.DeleteUser)
			systemUserGroup.PUT(":id/reset-password", systemUserApi.ResetPassword)
		}

		// System Role Router
		systemRoleApi := v1.SystemRoleApi{}
		permissionApi := v1.PermissionApi{}
		systemRoleGroup := privateGroup.Group("system/role")
		{
			systemRoleGroup.GET("list", systemRoleApi.GetRoleList)
			systemRoleGroup.POST("", systemRoleApi.CreateRole)
			systemRoleGroup.PUT(":id", systemRoleApi.UpdateRole)
			systemRoleGroup.DELETE(":id", systemRoleApi.DeleteRole)
			systemRoleGroup.GET(":id/menus", permissionApi.GetRoleMenus)
			systemRoleGroup.PUT(":id/menus", permissionApi.SetRoleMenus)
		}

		// System Menu Router
		systemMenuApi := v1.SystemMenuApi{}
		systemMenuGroup := privateGroup.Group("system/menu")
		{
			systemMenuGroup.GET("list", systemMenuApi.GetMenuList)
			systemMenuGroup.POST("", systemMenuApi.CreateMenu)
			systemMenuGroup.PUT(":id", systemMenuApi.UpdateMenu)
			systemMenuGroup.DELETE(":id", systemMenuApi.DeleteMenu)
		}

		// Operation Log Router
		operationLogApi := v1.OperationLogApi{}
		operationLogGroup := privateGroup.Group("system/log")
		{
			operationLogGroup.GET("list", operationLogApi.GetOperationLogList)
			operationLogGroup.DELETE("", operationLogApi.DeleteOperationLogs)
			operationLogGroup.DELETE("clear", operationLogApi.ClearOperationLogs)
		}

		// Profile Router
		profileApi := v1.ProfileApi{}
		profileGroup := privateGroup.Group("profile")
		{
			profileGroup.GET("", profileApi.GetProfile)
			profileGroup.PUT("", profileApi.UpdateProfile)
			profileGroup.PUT("password", profileApi.ChangePassword)
		}

		// User Permission Router (获取当前登录用户的权限信息)
		userGroup := privateGroup.Group("user")
		{
			userGroup.GET("permissions", permissionApi.GetUserPermissions)
			userGroup.GET("menus", permissionApi.GetUserMenus)
		}

		// Generator Router (代码生成器)
		generatorApi := v1.GeneratorApi{}
		generatorGroup := privateGroup.Group("generator")
		{
			generatorGroup.GET("tables", generatorApi.GetTables)
			generatorGroup.GET("columns", generatorApi.GetTableColumns)
			generatorGroup.POST("preview", generatorApi.PreviewCode)
			generatorGroup.POST("generate", generatorApi.GenerateCode)
		}
	}

	global.LV_LOG.Info("router register success")
}
