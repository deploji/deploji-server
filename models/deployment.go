package models

import (
	"fmt"
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/jinzhu/gorm"
	"github.com/sotomskir/mastermind-server/utils"
	"time"
)

type Deployment struct {
	gorm.Model
	StartedAt     time.Time
	FinishedAt    time.Time
	Application   Application `gorm:"association_autoupdate:false"`
	ApplicationID uint
	Version       string    `gorm:"type:text"`
	Inventory     Inventory `gorm:"association_autoupdate:false"`
	InventoryID   uint
	Status        Status
}

func GetDeployments(page utils.Page, filters []utils.Filter) ([]*Deployment, *pagination.Paginator) {
	var deployments []*Deployment
	db := GetDB().Preload("Application").Preload("Inventory")
	for _, filter := range filters {
		db = db.Where(fmt.Sprintf("%s=?", filter.Key), filter.Value)
	}
	paginator := pagination.Paging(&pagination.Param{
		DB:      db,
		Page:    page.Page,
		Limit:   page.Limit,
		OrderBy: page.OrderBy,
	}, &deployments)
	return deployments, paginator
}

func GetLatestDeployments() []*Deployment {
	var deployments []*Deployment
	sql := `
select * from deployments where id in (
	select max(id)
	from deployments
	where status = 2
	group by application_id, inventory_id
)`
	err := GetDB().Raw(sql).Scan(&deployments).Error
	if err != nil {
		return nil
	}
	return deployments
}

func GetDeployment(id uint) *Deployment {
	var deployment Deployment
	err := GetDB().Preload("Application.Project").Preload("Inventory").First(&deployment, id).Error
	if err != nil {
		return nil
	}
	return &deployment
}

func SaveDeployment(deployment *Deployment) error {
	err := GetDB().Create(deployment).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateDeploymentStatus(deployment *Deployment, updates map[string]interface{}) error {
	err := GetDB().Model(deployment).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}
