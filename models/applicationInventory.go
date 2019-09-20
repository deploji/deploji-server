package models

import (
	"github.com/jinzhu/gorm"
)

type ApplicationInventory struct {
	gorm.Model
	IsActive        bool
	Application     Application
	ApplicationID   uint
	Inventory       Inventory
	InventoryID     uint
	Name            string `gorm:"type:text"`
	ApplicationUrls string `gorm:"type:text"`
	ExtraVariables  string `gorm:"type:text"`
	Key             SshKey
	KeyID           uint
}

func GetApplicationInventories() ([]*ApplicationInventory, error) {
	applicationInventories := make([]*ApplicationInventory, 0)
	err := GetDB().Preload("Application").Preload("Inventory").Find(&applicationInventories).Error
	if err != nil {
		return nil, err
	}
	return applicationInventories, nil
}

func GetApplicationInventory(id uint) *ApplicationInventory {
	var applicationInventory ApplicationInventory
	err := GetDB().Preload("Application").Preload("Inventory").First(&applicationInventory, id).Error
	if err != nil {
		return nil
	}
	return &applicationInventory
}

func SaveApplicationInventory(applicationInventory *ApplicationInventory) error {
	if GetDB().NewRecord(applicationInventory) {
		err := GetDB().Create(applicationInventory).Error
		if err != nil {
			return err
		}
	} else {
		err := GetDB().Omit("created_at").Save(applicationInventory).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteApplicationInventory(applicationInventory *ApplicationInventory) error {
	err := GetDB().Delete(applicationInventory).Error
	if err != nil {
		return err
	}
	return nil
}
