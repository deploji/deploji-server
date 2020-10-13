package models

import (
	"github.com/jinzhu/gorm"
)

type PushSubscription struct {
	gorm.Model
	Endpoint string `gorm:"type:text"`
	Sub      string `gorm:"type:text"`
	UserID   uint
}

func FindByEndpoint(endpoint string) *PushSubscription {
	var pushSubscription PushSubscription
	err := GetDB().First(&pushSubscription, "endpoint = ?", endpoint).Error
	if err != nil {
		return &PushSubscription{}
	}
	return &pushSubscription
}

func FindByUserID(id uint) []*PushSubscription {
	var pushSubscriptions []*PushSubscription
	err := GetDB().Where("user_id = ?", id).Find(&pushSubscriptions).Error
	if err != nil {
		return nil
	}
	return pushSubscriptions
}

func SavePushSubscription(endpoint string, sub string, userID uint) error {
	record := FindByEndpoint(endpoint)
	record.UserID = userID
	record.Sub = sub
	record.Endpoint = endpoint
	if GetDB().NewRecord(record) {
		err := GetDB().Create(record).Error
		if err != nil {
			return err
		}
	} else {
		err := GetDB().Omit("created_at").Save(record).Error
		if err != nil {
			return err
		}
	}

	return nil
}
