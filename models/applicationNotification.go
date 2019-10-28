package models

type ApplicationNotification struct {
	ApplicationID         uint `gorm:"primary_key"`
	NotificationChannel   NotificationChannel
	NotificationChannelID uint `gorm:"primary_key"`
	SuccessEnabled        bool
	FailEnabled           bool
}

func GetApplicationNotifications(id uint) *[]ApplicationNotification {
	var notificationChannels []NotificationChannel
	if err := GetDB().Find(&notificationChannels).Error; err != nil {
		return nil
	}
	var notifications []ApplicationNotification
	if err := GetDB().
		Preload("NotificationChannel").
		Where("application_id=?", id).
		Find(&notifications).Error; err != nil {
		return nil
	}
	notificationsMap := make(map[uint]ApplicationNotification)
	for _, v := range notificationChannels {
		notificationsMap[v.ID] = ApplicationNotification{
			ApplicationID:         id,
			NotificationChannel:   v,
			NotificationChannelID: v.ID,
			SuccessEnabled:        false,
			FailEnabled:           false,
		}
	}
	for _, v := range notifications {
		notificationsMap[v.NotificationChannelID] = v
	}
	notifications = make([]ApplicationNotification, 0)
	for _, v := range notificationsMap {
		notifications = append(notifications, v)
	}

	return &notifications
}

func SaveApplicationNotification(notification *ApplicationNotification) error {
	if err := GetDB().Save(notification).Error; err != nil {
		return err
	}
	if err := GetDB().Preload("NotificationChannel").Find(notification).Error; err != nil {
		return err
	}
	return nil
}

func UpdateApplicationNotification(notification *ApplicationNotification) error {
	if err := GetDB().Model(notification).Updates(*notification).Error; err != nil {
		return err
	}
	if err := GetDB().Preload("NotificationChannel").Find(notification).Error; err != nil {
		return err
	}
	return nil
}
