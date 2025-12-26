package model

import "gorm.io/gorm"

// LvOperationLog 操作日志
type LvOperationLog struct {
	gorm.Model
	UserId    uint   `json:"userId" gorm:"comment:用户ID"`
	Username  string `json:"username" gorm:"size:64;comment:用户名"`
	Ip        string `json:"ip" gorm:"size:64;comment:IP地址"`
	Method    string `json:"method" gorm:"size:16;comment:请求方式"`
	Path      string `json:"path" gorm:"size:256;comment:请求路径"`
	Status    int    `json:"status" gorm:"comment:状态码"`
	Latency   int64  `json:"latency" gorm:"comment:耗时(ms)"`
	UserAgent string `json:"userAgent" gorm:"size:512;comment:User-Agent"`
	Body      string `json:"body" gorm:"type:text;comment:请求参数"`
	Response  string `json:"response" gorm:"type:text;comment:响应内容"`
	Module    string `json:"module" gorm:"size:64;comment:操作模块"`
	Action    string `json:"action" gorm:"size:64;comment:操作类型"`
}

func (LvOperationLog) TableName() string {
	return "lv_operation_logs"
}
