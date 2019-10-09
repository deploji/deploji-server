package models

import (
	"github.com/jinzhu/gorm"
)

type Application struct {
	gorm.Model
	Permissions        Permissions
	Name               string `gorm:"type:text"`
	AnsibleName        string `gorm:"type:text"`
	Project            Project
	ProjectID          uint
	Repository         Repository
	RepositoryID       uint
	RepositoryArtifact string
	Inventories        []ApplicationInventory
	AnsiblePlaybook    string `gorm:"type:text"`
	RepositoryGroup    string `gorm:"type:text"`
}

func GetApplications() ([]*Application, error) {
	applications := make([]*Application, 0)
	err := GetDB().
		Preload("Project").
		Preload("Repository").
		Preload("Inventories.Inventory").
		Find(&applications).Error
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func GetApplication(id uint) *Application {
	var application Application
	err := GetDB().
		Preload("Project").
		Preload("Repository").
		Preload("Inventories.Inventory").
		Preload("Inventories.Application").
		Preload("Inventories.Key").
		First(&application, id).Error
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
		for _, inventory := range application.Inventories {
			err := GetDB().Save(&inventory).Error
			if err != nil {
				return err
			}
		}
		if err := GetDB().
			Set("gorm:association_autocreate", false).
			Omit("created_at").
			Model(&application).
			Update(application).
			Error; err != nil {
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
