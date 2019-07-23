package models

type Status uint

const (
	StatusPending    Status = 0
	StatusProcessing Status = 1
	StatusCompleted  Status = 2
	StatusFailed     Status = 3
)
