package v1

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardApi struct{}

type DashboardStats struct {
	UserCount  int64 `json:"userCount"`
	RoleCount  int64 `json:"roleCount"`
	MenuCount  int64 `json:"menuCount"`
	TodayVisit int64 `json:"todayVisit"`
}

type ChartData struct {
	Categories []string `json:"categories"`
	Series     []int64  `json:"series"`
}

type DashboardCharts struct {
	VisitTrend  ChartData `json:"visitTrend"`  // 访问趋势
	UserGrowth  ChartData `json:"userGrowth"`  // 用户增长
	ModuleStats []PieItem `json:"moduleStats"` // 模块访问占比
	LatestLogs  []LogItem `json:"latestLogs"`  // 最近操作
}

type PieItem struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type LogItem struct {
	Username  string `json:"username"`
	Action    string `json:"action"`
	Module    string `json:"module"`
	CreatedAt string `json:"createdAt"`
}

// GetStats
// @Tags Dashboard
// @Summary Get dashboard statistics
// @Produce application/json
// @Success 200 {object} response.Response{data=DashboardStats}
// @Router /dashboard/stats [get]
func (d *DashboardApi) GetStats(c *gin.Context) {
	var stats DashboardStats

	// Count users
	global.LV_DB.Model(&model.LvUser{}).Count(&stats.UserCount)

	// Count roles
	global.LV_DB.Model(&model.LvRole{}).Count(&stats.RoleCount)

	// Count menus
	global.LV_DB.Model(&model.LvMenu{}).Count(&stats.MenuCount)

	// Today visit from operation logs
	today := time.Now().Format("2006-01-02")
	global.LV_DB.Model(&model.LvOperationLog{}).
		Where("DATE(created_at) = ?", today).
		Count(&stats.TodayVisit)

	c.JSON(200, gin.H{
		"code": 0,
		"data": stats,
		"msg":  "success",
	})
}

// GetCharts
// @Tags Dashboard
// @Summary Get dashboard chart data
// @Produce application/json
// @Success 200 {object} response.Response{data=DashboardCharts}
// @Router /dashboard/charts [get]
func (d *DashboardApi) GetCharts(c *gin.Context) {
	var charts DashboardCharts

	// 1. 最近7天访问趋势
	charts.VisitTrend = getVisitTrend()

	// 2. 用户增长趋势（最近7天）
	charts.UserGrowth = getUserGrowth()

	// 3. 模块访问占比
	charts.ModuleStats = getModuleStats()

	// 4. 最近操作日志
	charts.LatestLogs = getLatestLogs()

	c.JSON(200, gin.H{
		"code": 0,
		"data": charts,
		"msg":  "success",
	})
}

// 获取最近7天访问趋势
func getVisitTrend() ChartData {
	var result ChartData
	now := time.Now()

	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		result.Categories = append(result.Categories, date.Format("01/02"))

		var count int64
		global.LV_DB.Model(&model.LvOperationLog{}).
			Where("DATE(created_at) = ?", dateStr).
			Count(&count)
		result.Series = append(result.Series, count)
	}

	return result
}

// 获取用户增长趋势
func getUserGrowth() ChartData {
	var result ChartData
	now := time.Now()

	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		result.Categories = append(result.Categories, date.Format("01/02"))

		var count int64
		global.LV_DB.Model(&model.LvUser{}).
			Where("DATE(created_at) <= ?", dateStr).
			Count(&count)
		result.Series = append(result.Series, count)
	}

	return result
}

// 获取模块访问统计
func getModuleStats() []PieItem {
	var results []struct {
		Module string
		Count  int64
	}

	global.LV_DB.Model(&model.LvOperationLog{}).
		Select("module, COUNT(*) as count").
		Group("module").
		Order("count DESC").
		Limit(6).
		Scan(&results)

	var items []PieItem
	for _, r := range results {
		if r.Module != "" {
			items = append(items, PieItem{Name: r.Module, Value: r.Count})
		}
	}

	// 如果没有数据，返回示例数据
	if len(items) == 0 {
		items = []PieItem{
			{Name: "系统管理", Value: 45},
			{Name: "用户管理", Value: 30},
			{Name: "角色管理", Value: 15},
			{Name: "其他", Value: 10},
		}
	}

	return items
}

// 获取最近操作日志
func getLatestLogs() []LogItem {
	var logs []model.LvOperationLog
	global.LV_DB.Order("created_at DESC").Limit(10).Find(&logs)

	var items []LogItem
	for _, log := range logs {
		items = append(items, LogItem{
			Username:  log.Username,
			Action:    log.Action,
			Module:    log.Module,
			CreatedAt: log.CreatedAt.Format("15:04:05"),
		})
	}

	return items
}
