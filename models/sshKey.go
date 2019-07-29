package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type SshKey struct {
	gorm.Model
	Title string `gorm:"type:text"`
	Key   string `gorm:"type:text"`
}

func GetSshKeys() []*SshKey {
	keys := make([]*SshKey, 0)
	err := GetDB().Find(&keys).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return keys
}

func GetSshKey(id uint64) *SshKey {
	var key SshKey
	err := GetDB().First(&key, id).Error
	if err != nil {
		return nil
	}
	return &key
}


func SaveSshKey(key *SshKey) error {
	if GetDB().NewRecord(key) {
		err := GetDB().Create(key).Error
		if err != nil {
			return err
		}
	} else {
		err := GetDB().Omit("created_at").Save(key).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteSshKey(key *SshKey) error {
	err := GetDB().Delete(key).Error
	if err != nil {
		return err
	}
	return nil
}
