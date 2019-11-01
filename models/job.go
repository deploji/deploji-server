package models

import (
	"fmt"
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/deploji/deploji-server/utils"
	"github.com/jinzhu/gorm"
	"time"
)

type JobType string

const (
	TypeDeployment JobType = "Deployment"
	TypeJob        JobType = "Job"
	TypeSCMPull    JobType = "SCMPull"
)

type Job struct {
	gorm.Model
	Type           JobType
	StartedAt      time.Time
	FinishedAt     time.Time
	Application    Application
	ApplicationID  uint
	Playbook       string `gorm:"type:text"`
	Project        Project
	ProjectID      uint
	Inventory      Inventory
	InventoryID    uint
	Template       Template
	TemplateID     uint
	Key            SshKey
	KeyID          uint
	User           User
	UserID         uint
	Status         Status
	Version        string `gorm:"type:text"`
	ExtraVariables string `gorm:"type:text"`
}

func GetLatestDeployments() []*Job {
	var jobs []*Job
	sql := `
select * from jobs where id in (
	select max(id)
	from jobs
	where status = 2 and type = 'Deployment'
	group by application_id, inventory_id
)`
	err := GetDB().Raw(sql).Scan(&jobs).Error
	if err != nil {
		return nil
	}
	return jobs
}

func GetLatestSCMPulls() []*Job {
	var jobs []*Job
	sql := `
select * from jobs where id in (
	select max(id)
	from jobs
	where type = 'SCMPull'
	group by project_id
)`
	err := GetDB().Raw(sql).Scan(&jobs).Error
	if err != nil {
		return nil
	}
	return jobs
}

func GetJobs(page utils.Page, filters []utils.Filter) ([]*Job, *pagination.Paginator) {
	var jobs []*Job
	db := GetDB().
		Order("id desc").
		Preload("Application").
		Preload("Inventory").
		Preload("User").
		Preload("Project")
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
	err := GetDB().
		Preload("Application.Project").
		Preload("Project").
		Preload("User").
		Preload("Inventory").
		Preload("Key").
		First(&job, id).Error
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
