package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

type JobLog struct {
	gorm.Model	`json:"-"`
	Job   Job `json:"-"`
	JobID uint `json:"-";gorm:"index:job"`
	Order        uint `;json:"-"`
	Message      string `gorm:"type:text"`
}

func GetJobLogs(jobID uint) []*JobLog {
	var logs []*JobLog
	err := GetDB().Where(fmt.Sprintf("job_id = %d", jobID)).Order("id").Find(&logs).Error
	if err != nil {
		return nil
	}
	return logs
}

func SaveJobLog(jobLog *JobLog) {
	log.Println(jobLog.Message)
	err := GetDB().Create(jobLog).Error
	if err != nil {
		log.Printf("Error saving job log: %s", err)
	}
}
