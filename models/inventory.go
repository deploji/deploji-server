package models

import (
	"github.com/jinzhu/gorm"
)

type Inventory struct {
	gorm.Model
	Permissions
	Name                   string `gorm:"type:text"`
	Project                Project
	ProjectID              uint
	SourceFile             string
	ApplicationInventories []ApplicationInventory
}

func GetInventories() []*Inventory {
	inventories := make([]*Inventory, 0)
	err := GetDB().
		Preload("Project").
		Preload("ApplicationInventories.Application").
		Find(&inventories).Error
	if err != nil {
		return nil
	}
	return inventories
}

func GetInventory(id uint) *Inventory {
	var inventory Inventory
	err := GetDB().Preload("Project").First(&inventory, id).Error
	if err != nil {
		return nil
	}
	return &inventory
}

func SaveInventory(inventory *Inventory) error {
	if GetDB().NewRecord(inventory) {
		err := GetDB().Create(inventory).Error
		if err != nil {
			return err
		}
	} else {
		err := GetDB().Omit("created_at").Save(inventory).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteInventory(inventory *Inventory) error {
	err := GetDB().Delete(inventory).Error
	if err != nil {
		return err
	}
	return nil
}

func GetInventoriesByApplicationId(id uint) *[]Inventory {
	var inventories []Inventory
	var appInventories []ApplicationInventory
	err := GetDB().
		Preload("Inventory").
		Where(&ApplicationInventory{IsActive: true, ApplicationID: id}).
		Find(&appInventories).Error
	if err != nil {
		return nil
	}
	for _, v := range appInventories {
		inventories = append(inventories, v.Inventory)
	}
	return &inventories
}
