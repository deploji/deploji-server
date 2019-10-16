package models

import (
	"github.com/jinzhu/gorm"
)

type SurveyInput struct {
	gorm.Model
	Label        string `gorm:"type:text"`
	Hint         string `gorm:"type:text"`
	VariableName string `gorm:"type:text"`
	Type         string `gorm:"type:text"`
	Options      string `gorm:"type:text"`
	SurveyID     uint
}

func GetSurveyInputsByTemplateID(id uint) (*[]SurveyInput, error) {
	var survey Survey
	err := GetDB().
		Preload("Inputs").
		First(&survey, "template_id=?", id).
		Error
	if err != nil {
		return nil, err
	}
	if survey.Inputs == nil {
		survey.Inputs = make([]SurveyInput, 0)
	}
	return &survey.Inputs, nil
}

func GetSurveyInput(id uint) *SurveyInput {
	var surveyInput SurveyInput
	err := GetDB().
		Preload("Survey").
		First(&surveyInput, id).Error
	if err != nil {
		return nil
	}
	return &surveyInput
}

func SaveSurveyInput(surveyInput *SurveyInput) error {
	if GetDB().NewRecord(surveyInput) {
		err := GetDB().Create(surveyInput).Error
		if err != nil {
			return err
		}
	} else {
		if err := GetDB().
			Set("gorm:association_autocreate", false).
			Omit("created_at").
			Model(&surveyInput).
			Update(surveyInput).
			Error; err != nil {
			return err
		}
	}
	return nil
}

func DeleteSurveyInput(surveyInput *SurveyInput) error {
	err := GetDB().Delete(surveyInput).Error
	if err != nil {
		return err
	}
	return nil
}
