package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Project struct {
	gorm.Model
	Name       string `gorm:"type:text"`
	RepoUrl    string `gorm:"type:text"`
	RepoBranch string `gorm:"type:text"`
	RepoUser   string `gorm:"type:text"`
	SshKeyID   uint
	SshKey     SshKey
}

func GetProjects() []*Project {
	projects := make([]*Project, 0)
	err := GetDB().Preload("SshKey").Find(&projects).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return projects
}

func GetProject(id uint64) *Project {
	var project Project
	err := GetDB().Preload("SshKey").First(&project, id).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &project
}

func SaveProject(project *Project) error {
	if GetDB().NewRecord(project) {
		err := GetDB().Create(project).Error
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else {
		err := GetDB().Omit("created_at").Save(project).Error
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
