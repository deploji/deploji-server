package models

import (
	"github.com/jinzhu/gorm"
)

type Deployment struct {
	gorm.Model
	Application   Application
	ApplicationID uint
	Version       string `gorm:"type:text"`
	Inventory     Inventory
	InventoryID   uint
}

func GetDeployments() []*Deployment {
	var deployments []*Deployment
	err := GetDB().Preload("Application").Preload("Inventory").Order("id desc").Find(&deployments).Error
	if err != nil {
		return nil
	}
	return deployments
}

func GetDeployment(id uint) *Deployment {
	var deployment Deployment
	err := GetDB().Preload("Application").Preload("Inventory").First(&deployment, id).Error
	if err != nil {
		return nil
	}
	return &deployment
}

func SaveDeployment(deployment *Deployment) error {
	err := GetDB().Create(deployment).Error
	if err != nil {
		return err
	}
	return nil
}
