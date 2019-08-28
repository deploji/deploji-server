package models

import (
	"github.com/jinzhu/gorm"
)

type Repository struct {
	gorm.Model
	Name      string `gorm:"type:text"`
	Type      string `gorm:"type:text"`
	Url       string `gorm:"type:text"`
	Username  string `gorm:"type:text"`
	Password  string `gorm:"type:text"`
	NexusName string `gorm:"type:text"`
}

func GetRepositories() []*Repository {
	var repositories []*Repository
	err := GetDB().Find(&repositories).Error
	if err != nil {
		return nil
	}
	return repositories
}

func GetRepository(id uint) *Repository {
	var repository Repository
	err := GetDB().First(&repository, id).Error
	if err != nil {
		return nil
	}
	return &repository
}

func SaveRepository(repository *Repository) error {
	if GetDB().NewRecord(repository) {
		err := GetDB().Create(repository).Error
		if err != nil {
			return err
		}
	} else {
		err := GetDB().Omit("created_at").Save(repository).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteRepository(repository *Repository) error {
	err := GetDB().Delete(repository).Error
	if err != nil {
		return err
	}
	return nil
}
