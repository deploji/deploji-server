package models

import (
	"github.com/jinzhu/gorm"
)

type Setting struct {
	gorm.Model
	SettingGroup   SettingGroup
	SettingGroupID uint
	Key            string `gorm:"type:text"`
	Value          string `gorm:"type:text"`
	BoolValue      bool
	ValueType      string `gorm:"type:text"`
	Label          string `gorm:"type:text"`
	Description    string `gorm:"type:text"`
}

func GetSettingValue(groupName string, key string, defaultValue string) string {
	var group SettingGroup
	var setting Setting
	if err := GetDB().
		Where(SettingGroup{Name: groupName}).
		Find(&group).Error; err != nil {
			return defaultValue
	}
	if err := GetDB().
		Where(Setting{Key: key, SettingGroupID: group.ID}).
		Find(&setting).Error; err != nil {
		return defaultValue
	}
	return setting.Value
}

func GetSettingBoolValue(groupName string, key string, defaultValue bool) bool {
	var group SettingGroup
	var setting Setting
	if err := GetDB().
		Where(SettingGroup{Name: groupName}).
		Find(&group).Error; err != nil {
		return defaultValue
	}
	if err := GetDB().
		Where(Setting{Key: key, SettingGroupID: group.ID}).
		Find(&setting).Error; err != nil {
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
