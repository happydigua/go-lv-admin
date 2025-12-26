package core

import (
	"fmt"
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// InitDemoData 初始化演示数据
func InitDemoData(db *gorm.DB) {
	var count int64
	db.Model(&model.LvDemo{}).Count(&count)
	if count > 0 {
		return
	}

	categories := []string{"电子产品", "家居用品", "办公文具", "日用百货", "服装服饰"}

	var demos []model.LvDemo
	for i := 1; i <= 55; i++ {
		demos = append(demos, model.LvDemo{
			Code:        fmt.Sprintf("DEMO%05d", i),
			Name:        fmt.Sprintf("测试商品 %d", i),
			Category:    categories[rand.Intn(len(categories))],
			Status:      rand.Intn(2) + 1, // 1 or 2
			Amount:      float64(rand.Intn(100000)) / 100,
			Description: fmt.Sprintf("这是一个测试商品的详细描述信息 %d ...", i),
		})
	}

	if err := db.CreateInBatches(demos, 100).Error; err != nil {
		global.LV_LOG.Error("init demo data failed")
	} else {
		global.LV_LOG.Info("init demo data success")
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
