package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

type DeploymentLog struct {
	gorm.Model	`json:"-"`
	Deployment   Deployment `json:"-"`
	DeploymentID uint `json:"-";gorm:"index:deployment"`
	Order        uint `;json:"-"`
	Message      string `gorm:"type:text"`
}

func GetDeploymentLogs(deploymentID uint) []*DeploymentLog {
	var logs []*DeploymentLog
	err := GetDB().Where(fmt.Sprintf("deployment_id = %d", deploymentID)).Order("id").Find(&logs).Error
	if err != nil {
		return nil
	}
	return logs
}

func SaveDeploymentLog(deploymentLog *DeploymentLog) {
	log.Println(deploymentLog.Message)
	err := GetDB().Create(deploymentLog).Error
	if err != nil {
		log.Printf("Error saving deployment log: %s", err)
	}
}
