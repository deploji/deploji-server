package models

import (
	"fmt"
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/jinzhu/gorm"
	"github.com/sotomskir/mastermind-server/utils"
	"time"
)

type Job struct {
	gorm.Model
	StartedAt   time.Time
	FinishedAt  time.Time
	Project     Application `gorm:"association_autoupdate:false"`
	ProjectID   uint
	Inventory   Inventory `gorm:"association_autoupdate:false"`
	InventoryID uint
	Status      Status
}

func GetJobs(page utils.Page, filters []utils.Filter) ([]*Job, *pagination.Paginator) {
	var jobs []*Job
	db := GetDB().Preload("Application").Preload("Inventory")
	for _, filter := range filters {
		db = db.Where(fmt.Sprintf("%s=?", filter.Key), filter.Value)
	}
	paginator := pagination.Paging(&pagination.Param{
		DB:      db,
		Page:    page.Page,
		Limit:   page.Limit,
		OrderBy: page.OrderBy,
	}, &jobs)
	return jobs, paginator
}

func GetJob(id uint) *Job {
	var job Job
	err := GetDB().First(&job, id).Error
	if err != nil {
		return nil
	}
	return &job
}

func SaveJob(job *Job) error {
	err := GetDB().Create(job).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateJobStatus(job *Job, updates map[string]interface{}) error {
	err := GetDB().Model(job).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}
