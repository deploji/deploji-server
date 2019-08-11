package dto

import (
	"encoding/json"
	"github.com/sotomskir/mastermind-server/models"
	"log"
)

type JobType string

const (
	Deployment JobType = "Deployment"
	Job        JobType = "Job"
	SCMPull    JobType = "SCMPull"
)

type JobMessage struct {
	Type JobType
	ID   uint
}

type Message []byte

type StatusMessage struct {
	Type   JobType
	ID     uint
	Status models.Status
}

func NewStatusMessage(jobType JobType, id uint, status models.Status) Message {
	message := StatusMessage{
		Type:   jobType,
		ID:     id,
		Status: status,
	}
	return MarshallMessage(message)
}

func MarshallMessage(data interface{}) []byte {
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling message: %s", err)
	}
	return body
}
