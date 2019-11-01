package models

import (
	"github.com/jinzhu/gorm"
)

type Team struct {
	gorm.Model
	Name string `gorm:"type:text"`
}

func GetTeams() []*Team {
	teams := make([]*Team, 0)
	err := GetDB().Order("name asc").Find(&teams).Error
	if err != nil {
		return nil
	}
	return teams
}

func GetTeam(id uint) *Team {
	var team Team
	err := GetDB().First(&team, id).Error
	if err != nil {
		return nil
	}
	return &team
}

func SaveTeam(team *Team) error {
	if GetDB().NewRecord(team) {
		err := GetDB().Create(team).Error
		if err != nil {
			return err
		}
	} else {
		err := GetDB().Omit("created_at").Save(team).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteTeam(team *Team) error {
	err := GetDB().Delete(team).Error
	if err != nil {
		return err
	}
	return nil
}
