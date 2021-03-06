package models

import (
	"github.com/jinzhu/gorm"
)

type Template struct {
	gorm.Model
	Permissions          Permissions
	Name                 string `gorm:"type:text"`
	Project              Project
	ProjectID            uint
	Inventory            Inventory
	InventoryID          uint
	SshKey               SshKey
	SshKeyID             uint
	VaultKey             SshKey
	VaultKeyID           uint
	Playbook             string `gorm:"type:text"`
	PromptSshKey         bool
	PromptVaultKey       bool
	PromptPlaybook       bool
	PromptInventory      bool
	PromptProject        bool
	PromptExtraVariables bool
	ExtraVariables       string `gorm:"type:text"`
	Survey               Survey
}

func GetTemplates() []*Template {
	templates := make([]*Template, 0)
	err := GetDB().
		Preload("Project").
		Preload("Inventory").
		Preload("SshKey").
		Preload("VaultKey").
		Order("name asc").
		Find(&templates).Error
	if err != nil {
		return nil
	}
	return templates
}

func GetTemplate(id uint) *Template {
	var template Template
	err := GetDB().
		Preload("Project").
		Preload("Inventory").
		Preload("Survey").
		Preload("SshKey").
		Preload("VaultKey").
		First(&template, id).Error
	if err != nil {
		return nil
	}
	return &template
}

func SaveTemplate(template *Template) error {
	if GetDB().NewRecord(template) {
		err := GetDB().Create(template).Error
		if err != nil {
			return err
		}
	} else {
		err := GetDB().Omit("created_at").Save(template).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteTemplate(template *Template) error {
	err := GetDB().Delete(template).Error
	if err != nil {
		return err
	}
	return nil
}
