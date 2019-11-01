package models

type ProjectNotification struct {
	ProjectID             uint `gorm:"primary_key"`
	NotificationChannel   NotificationChannel
	NotificationChannelID uint `gorm:"primary_key"`
	SuccessEnabled        bool
	FailEnabled           bool
}

func GetProjectNotifications(id uint) *[]ProjectNotification {
	var notificationChannels []NotificationChannel
	if err := GetDB().Find(&notificationChannels).Error;
		err != nil {
		return nil
	}
	var notifications []ProjectNotification
	if err := GetDB().
		Order("notification_channel_id asc").
		Preload("NotificationChannel").
		Where("project_id=?", id).
		Find(&notifications).Error;
		err != nil {
		return nil
	}
	notificationsMap := make(map[uint]ProjectNotification)
	for _, v := range notificationChannels {
		notificationsMap[v.ID] = ProjectNotification{
			ProjectID:             id,
			NotificationChannel:   v,
			NotificationChannelID: v.ID,
			SuccessEnabled:        false,
			FailEnabled:           false,
		}
	}
	for _, v := range notifications {
		notificationsMap[v.NotificationChannelID] = v
	}
	notifications = make([]ProjectNotification, 0)
	for _, v := range notificationsMap {
		notifications = append(notifications, v)
	}

	return &notifications
}

func SaveProjectNotification(notification *ProjectNotification) error {
	if err := GetDB().Save(notification).Error; err != nil {
		return err
	}
	if err := GetDB().Preload("NotificationChannel").Find(notification).Error; err != nil {
		return err
	}
	return nil
}
