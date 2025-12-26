package model

import (
	"gorm.io/gorm"
)

type LvDemo struct {
	gorm.Model
	Code        string  `json:"code" gorm:"type:varchar(50);uniqueIndex;comment:编号"`
	Name        string  `json:"name" gorm:"type:varchar(100);comment:名称"`
	Category    string  `json:"category" gorm:"type:varchar(50);comment:分类"`
	Status      int     `json:"status" gorm:"default:1;comment:状态 1正常 2停用"`
	Amount      float64 `json:"amount" gorm:"type:decimal(10,2);comment:金额"`
	Description string  `json:"description" gorm:"type:varchar(255);comment:描述"`
}
