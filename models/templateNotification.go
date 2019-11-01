package models

type TemplateNotification struct {
	TemplateID            uint `gorm:"primary_key"`
	RelatedNotification
}

func GetTemplateNotifications(id uint) *[]TemplateNotification {
	var notificationChannels []NotificationChannel
	if err := GetDB().Find(&notificationChannels).Error;
		err != nil {
		return nil
	}
	var notifications []TemplateNotification
	if err := GetDB().
		Order("notification_channel_id asc").
		Preload("NotificationChannel").
		Where("template_id=?", id).
		Find(&notifications).Error;
		err != nil {
		return nil
	}
	notificationsMap := make(map[uint]TemplateNotification)
	for _, v := range notificationChannels {
		notificationsMap[v.ID] = TemplateNotification{
			TemplateID:            id,
			RelatedNotification: RelatedNotification{
				NotificationChannel:   v,
				NotificationChannelID: v.ID,
				SuccessEnabled:        false,
				FailEnabled:           false,
			},
		}
	}
	for _, v := range notifications {
		notificationsMap[v.NotificationChannelID] = v
	}
	notifications = make([]TemplateNotification, 0)
	for _, v := range notificationsMap {
		notifications = append(notifications, v)
	}

	return &notifications
}

func SaveTemplateNotification(notification *TemplateNotification) error {
	if err := GetDB().Save(notification).Error; err != nil {
		return err
	}
	if err := GetDB().Preload("NotificationChannel").Find(notification).Error; err != nil {
		return err
	}
	return nil
}
