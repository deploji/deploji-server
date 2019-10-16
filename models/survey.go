package models

import (
	"github.com/jinzhu/gorm"
)

type Survey struct {
	gorm.Model
	Enabled      bool
	Inputs []SurveyInput
	TemplateID   uint
}

func GetSurveyByTemplateID(id uint) *Survey {
	var survey Survey
	err := GetDB().
		Preload("Inputs").
		First(&survey, "template_id=?", id).Error
	if err != nil {
		return nil
	}
	return &survey
}

func SaveSurvey(survey *Survey) error {
	if GetDB().NewRecord(survey) {
		err := GetDB().Create(survey).Error
		if err != nil {
			return err
		}
	} else {
		if err := GetDB().
			Set("gorm:association_autocreate", false).
			Omit("created_at").
			Model(&survey).
			Update(survey).
			Error; err != nil {
			return err
		}
	}
	return nil
}

func DeleteSurvey(survey *Survey) error {
	err := GetDB().Delete(survey).Error
	if err != nil {
		return err
	}
	return nil
}
