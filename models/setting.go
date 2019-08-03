package models

import (
	"github.com/jinzhu/gorm"
)

type Setting struct {
	gorm.Model
	SettingGroupID uint
	Key            string `gorm:"type:text"`
	Value          string `gorm:"type:text"`
	BoolValue      bool
	ValueType      string `gorm:"type:text"`
	Label          string `gorm:"type:text"`
	Description    string `gorm:"type:text"`
}

func GetSettingValue(group string, key string, defaultValue string) string {
	var setting Setting
	err := GetDB().Where("group=?", group).Where("key=?", key).Find(&setting).Error
	if err != nil || setting.Value == "" {
		return defaultValue
	}
	return setting.Value
}

func GetSettingBoolValue(group string, key string, defaultValue bool) bool {
	var setting Setting
	err := GetDB().Where("group=?", group).Where("key=?", key).Find(&setting).Error
	if err != nil || setting.Value == "" {
		return defaultValue
	}
	return setting.BoolValue
}

func SaveSettings(settings *[]Setting) error {
	for _, setting := range *settings {
		err := GetDB().Table("settings").Select("Value", "BoolValue").Save(setting).Error
		if err != nil {
			return err
		}
	}
	return nil
}
