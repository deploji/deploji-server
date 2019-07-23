package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type DeploymentLog struct {
	gorm.Model
	Deployment   Deployment
	DeploymentID uint `gorm:"index:deployment"`
	Order        uint
	Message      string `gorm:"type:text"`
}

func GetDeploymentLogs(deploymentID uint) []*DeploymentLog {
	var logs []*DeploymentLog
	err := GetDB().Where(fmt.Sprintf("deployment_id = %d", deploymentID)).Order("order").Find(&logs).Error
	if err != nil {
		return nil
	}
	return logs
}

func SaveDeploymentLog(deploymentLog *DeploymentLog) error {
	err := GetDB().Create(deploymentLog).Error
	if err != nil {
		return err
	}
	return nil
}
