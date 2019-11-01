package models

type RelatedNotification struct {
	NotificationChannel   NotificationChannel
	NotificationChannelID uint `gorm:"primary_key"`
	SuccessEnabled        bool
	FailEnabled           bool
}

func GetRelatedNotifications(job Job) (error, []RelatedNotification) {
	templateNotifications := GetTemplateNotifications(job.TemplateID)
	applicationNotifications := GetApplicationNotifications(job.ApplicationID)
	projectNotifications := GetProjectNotifications(job.ProjectID)
	relatedNotifications := make([]RelatedNotification, 0)
	for _, n := range *templateNotifications {
		relatedNotifications = append(relatedNotifications, n.RelatedNotification)
	}
	for _, n := range *applicationNotifications {
		relatedNotifications = append(relatedNotifications, n.RelatedNotification)
	}
	for _, n := range *projectNotifications {
		relatedNotifications = append(relatedNotifications, n.RelatedNotification)
	}
	return nil, relatedNotifications
}
