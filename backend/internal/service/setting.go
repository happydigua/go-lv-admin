package service

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
)

type SettingService struct{}

// GetAllSettings 获取所有设置
func (s *SettingService) GetAllSettings() (map[string]interface{}, error) {
	var settings []model.LvSetting
	if err := global.LV_DB.Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}
	return result, nil
}

// GetSetting 获取单个设置
func (s *SettingService) GetSetting(key string) (string, error) {
	var setting model.LvSetting
	if err := global.LV_DB.Where("`key` = ?", key).First(&setting).Error; err != nil {
		return "", err
	}
	return setting.Value, nil
}

// UpdateSetting 更新设置
func (s *SettingService) UpdateSetting(key, value string) error {
	return global.LV_DB.Model(&model.LvSetting{}).Where("`key` = ?", key).Update("value", value).Error
}

// BatchUpdateSettings 批量更新设置
func (s *SettingService) BatchUpdateSettings(settings map[string]string) error {
	for key, value := range settings {
		if err := s.UpdateSetting(key, value); err != nil {
			return err
		}
	}
	return nil
}

// GetPublicSettings 获取公开设置（无需登录即可获取）
func (s *SettingService) GetPublicSettings() (map[string]interface{}, error) {
	publicKeys := []string{"site_name", "site_logo", "site_footer"}
	var settings []model.LvSetting
	if err := global.LV_DB.Where("`key` IN ?", publicKeys).Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}
	return result, nil
}
