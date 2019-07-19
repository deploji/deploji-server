package models

import (
	"fmt"
)

type SshKey struct {
	ID uint
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
		fmt.Println(err)
		return nil
	}
	return &key
}

func SaveSshKey(key *SshKey) *SshKey {
	if GetDB().NewRecord(key) {
		err := GetDB().Create(key).Error
		if err != nil {
			fmt.Println(err)
			return nil
		}
	} else {
		err := GetDB().Omit("created_at").Update(key).Error
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	return key
}

func DeleteSshKey(key *SshKey) error {
	err := GetDB().Delete(key).Error
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
