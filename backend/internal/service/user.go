package service

import (
	"errors"
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"go-lv-vue-admin/pkg/utils"

	"gorm.io/gorm"
)

type UserService struct{}

func (s *UserService) Login(u *model.LvUser) (userInter *model.LvUser, err error) {
	if global.LV_DB == nil {
		return nil, errors.New("db not initialized")
	}

	var user model.LvUser
	err = global.LV_DB.Where("username = ?", u.Username).Preload("Role").First(&user).Error
	if err == nil {
		// 使用 bcrypt 验证密码
		if !utils.CheckPassword(u.Password, user.Password) {
			return nil, errors.New("password incorrect")
		}
		return &user, nil
	}
	return nil, err
}

func (s *UserService) CreateToken(user model.LvUser) (string, int64, error) {
	j := utils.NewJWT()
	claims := j.CreateClaims(utils.BaseClaims{
		UserId:   user.ID,
		Username: user.Username,
		RoleId:   user.RoleId,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		return "", 0, err
	}
	// Calculate actual expiration timestamp
	exp, _ := utils.ParseDuration(global.LV_CONFIG.JWT.ExpiresTime)
	expiresAt := int64(exp.Seconds())
	return token, expiresAt, nil
}

// Register (Optional MVP)
func (s *UserService) Register(u model.LvUser) (userInter model.LvUser, err error) {
	// Check if user exists
	var user model.LvUser
	if !errors.Is(global.LV_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) {
		return userInter, errors.New("username already exists")
	}
	// 使用 bcrypt 加密密码
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return userInter, err
	}
	u.Password = hashedPassword
	err = global.LV_DB.Create(&u).Error
	return u, err
}
