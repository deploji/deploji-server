package models

import (
	"github.com/jinzhu/gorm"
)

type NotificationChannel struct {
	gorm.Model
	Name        string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	Type        string `gorm:"type:text"`
	Recipients  string `gorm:"type:text"`
	WebhookURL  string `gorm:"type:text"`
	UserID      uint
}

func GetNotificationChannels() []*NotificationChannel {
	notificationChannels := make([]*NotificationChannel, 0)
	err := GetDB().Where("user_id = 0").Order("name asc").Find(&notificationChannels).Error
	if err != nil {
		return notificationChannels
	}
	return notificationChannels
}

func GetNotificationChannel(id uint) *NotificationChannel {
	var notificationChannel NotificationChannel
	err := GetDB().First(&notificationChannel, id).Error
	if err != nil {
		return nil
	}
	return &notificationChannel
}

func SaveNotificationChannel(notificationChannel *NotificationChannel) error {
	if GetDB().NewRecord(notificationChannel) {
		err := GetDB().Create(notificationChannel).Error
		if err != nil {
			return err
		}
	} else {
		err := GetDB().Omit("created_at").Save(notificationChannel).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteNotificationChannel(notificationChannel *NotificationChannel) error {
	err := GetDB().Delete(notificationChannel).Error
	if err != nil {
		return err
	}
	return nil
}
