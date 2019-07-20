package models

import (
	"github.com/jinzhu/gorm"
)

type Application struct {
	gorm.Model
	Name               string `gorm:"type:text"`
	AnsibleName        string `gorm:"type:text"`
	Project            Project
	ProjectID          uint
	Repository         Repository
	RepositoryID       uint
	RepositoryArtifact string
	AnsiblePlaybook    string `gorm:"type:text"`
}

func GetApplications() ([]*Application, error) {
	applications := make([]*Application, 0)
	err := GetDB().Preload("Project").Preload("Repository").Find(&applications).Error
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func GetApplication(id uint) *Application {
	var application Application
	err := GetDB().Preload("Project").Preload("Repository").First(&application, id).Error
	if err != nil {
		return nil
	}
	return &application
}

func SaveApplication(application *Application) error {
	if GetDB().NewRecord(application) {
		err := GetDB().Create(application).Error
		if err != nil {
			return err
		}
	} else {
		err := GetDB().Omit("created_at").Save(application).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteApplication(application *Application) error {
	err := GetDB().Delete(application).Error
	if err != nil {
		return err
	}
	return nil
}
