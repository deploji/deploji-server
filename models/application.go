package models

import (
	"fmt"
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
	Inventories        []ApplicationInventory
	AnsiblePlaybook    string `gorm:"type:text"`
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
		//if err := GetDB().
		//	Table("application_inventories").
		//	Where("application_id=?", application.ID).
		//	UpdateColumn("deleted_at", nil).Error; err != nil {
		//	return err
		//}

		var inventoryIds []uint
		for _, inventory := range application.Inventories {
			fmt.Printf("%#v", inventory.IsActive)
			if err := GetDB().
				Model(&inventory).
				Updates(map[string]interface{}{
					"IsActive":        inventory.IsActive,
					"InventoryID":     inventory.InventoryID,
					"ApplicationUrls": inventory.ApplicationUrls,
					"KeyID":           inventory.KeyID,
				}).
				Error; err != nil {
				return err
			}
			inventoryIds = append(inventoryIds, inventory.InventoryID)
		}

		if err := GetDB().
			Omit("created_at").
			Model(&application).
			Update(application).
			Error; err != nil {
			return err
		}

		//if err := GetDB().
		//	Table("application_inventories").
		//	Where("application_id=?", application.ID).
		//	UpdateColumn("deleted_at", time.Now()).Error; err != nil {
		//	return err
		//}
		//if err := GetDB().
		//	Table("application_inventories").
		//	Where("application_id=?", application.ID).
		//	Where("inventory_id IN (?)", inventoryIds).
		//	UpdateColumn("deleted_at", nil).Error; err != nil {
		//	return err
		//}
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
