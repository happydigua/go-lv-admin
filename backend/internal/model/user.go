package model

import (
	"gorm.io/gorm"
)

type LvUser struct {
	gorm.Model
	Username string `json:"username" gorm:"index;comment:用户登录名"`
	Password string `json:"-"  gorm:"comment:用户登录密码"`
	Nickname string `json:"nickname" gorm:"default:系统用户;comment:用户昵称"`
	Avatar   string `json:"avatar" gorm:"default:https://via.placeholder.com/200;comment:用户头像"`
	Email    string `json:"email" gorm:"comment:用户邮箱"`
	Phone    string `json:"phone" gorm:"comment:用户手机号"`
	Status   int    `json:"status" gorm:"default:1;comment:用户状态 1正常 2冻结"`
	RoleId   uint   `json:"role_id" gorm:"comment:用户角色ID"`
	Role     LvRole `json:"Role" gorm:"foreignKey:RoleId"`
}

func (LvUser) TableName() string {
	return "lv_users"
}
