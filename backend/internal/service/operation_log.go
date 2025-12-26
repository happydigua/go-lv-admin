package service

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
)

type OperationLogService struct{}

// GetOperationLogList 获取操作日志列表
func (s *OperationLogService) GetOperationLogList(page, pageSize int, username, module, action string) ([]model.LvOperationLog, int64, error) {
	var logs []model.LvOperationLog
	var total int64

	db := global.LV_DB.Model(&model.LvOperationLog{})

	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if module != "" {
		db = db.Where("module = ?", module)
	}
	if action != "" {
		db = db.Where("action = ?", action)
	}

	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

// DeleteOperationLogs 批量删除操作日志
func (s *OperationLogService) DeleteOperationLogs(ids []uint) error {
	return global.LV_DB.Delete(&model.LvOperationLog{}, ids).Error
}

// ClearOperationLogs 清空操作日志
func (s *OperationLogService) ClearOperationLogs() error {
	return global.LV_DB.Where("1=1").Delete(&model.LvOperationLog{}).Error
}
