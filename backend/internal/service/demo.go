package service

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
)

type DemoService struct{}

// GetDemoList 获取演示数据列表
func (s *DemoService) GetDemoList(page, pageSize int, filters map[string]interface{}, sortField, sortOrder string) ([]model.LvDemo, int64, error) {
	var list []model.LvDemo
	var total int64
	db := global.LV_DB.Model(&model.LvDemo{})

	// 过滤
	if name, ok := filters["name"]; ok && name != "" {
		db = db.Where("name LIKE ?", "%"+name.(string)+"%")
	}
	if category, ok := filters["category"]; ok && category != "" {
		db = db.Where("category = ?", category)
	}
	if status, ok := filters["status"]; ok && status != "" {
		db = db.Where("status = ?", status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	order := "id DESC"
	if sortField != "" {
		if sortOrder == "ascend" {
			order = sortField + " ASC"
		} else if sortOrder == "descend" {
			order = sortField + " DESC"
		}
	}
	db = db.Order(order)

	// 分页
	if pageSize > 0 {
		offset := (page - 1) * pageSize
		db = db.Offset(offset).Limit(pageSize)
	}

	if err := db.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// CreateDemo 创建
func (s *DemoService) CreateDemo(demo *model.LvDemo) error {
	return global.LV_DB.Create(demo).Error
}

// UpdateDemo 更新
func (s *DemoService) UpdateDemo(demo *model.LvDemo) error {
	return global.LV_DB.Model(demo).Updates(demo).Error
}

// DeleteDemo 删除
func (s *DemoService) DeleteDemo(id uint) error {
	return global.LV_DB.Delete(&model.LvDemo{}, id).Error
}
