package models

import (
	"github.com/jinzhu/gorm"
)

type PushSubscription struct {
	gorm.Model
	Endpoint       string `gorm:"type:text"`
	ExpirationTime uint
	P256dh         string `gorm:"type:text"`
	Auth           string `gorm:"type:text"`
	UserID         uint
}

func FindByEndpoint(endpoint string) *PushSubscription {
	var pushSubscription PushSubscription
	err := GetDB().First(&pushSubscription, "endpoint = ?", endpoint).Error
	if err != nil {
		return &PushSubscription{}
	}
	return &pushSubscription
}

func SavePushSubscription(endpoint string, expirationTime uint, auth string, p256dh string, userID uint) error {
	record := FindByEndpoint(endpoint)
	record.UserID = userID
	record.ExpirationTime = expirationTime
	record.Endpoint = endpoint
	record.Auth = auth
	record.P256dh = p256dh
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
