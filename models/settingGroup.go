package models

import (
	"github.com/jinzhu/gorm"
)

type SettingGroup struct {
	gorm.Model
	Name     string `gorm:"type:text"`
	Settings []Setting
}

func GetSettingGroups() []*SettingGroup {
	settingGroups := make([]*SettingGroup, 0)
	err := GetDB().Preload("Settings").Find(&settingGroups).Error
	if err != nil {
		return nil
	}
	return settingGroups
}
