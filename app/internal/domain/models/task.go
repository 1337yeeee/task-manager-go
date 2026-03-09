package models

import "time"

type Task struct {
	ID          string `gorm:"primaryKey;size:40"`
	ProjectID   string `gorm:"size:40;not null"`
	Name        string `gorm:"size:50;not null"`
	Content     string
	ExecutiveID string `gorm:"size:40;not null"`
	AuditorID   string `gorm:"size:40;"`
	CreatedBy   string `gorm:"size:40;"`
	UpdatedBy   string `gorm:"size:40;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      TaskStatus `gorm:"size:50;not null"`
}

type TaskStatus string

const (
	TaskStatusCreated    TaskStatus = "created"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusAudit      TaskStatus = "audit"
	TaskStatusDone       TaskStatus = "done"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusCreated,
		TaskStatusInProgress,
		TaskStatusAudit,
		TaskStatusDone:
		return true
	}
	return false
}

func ParseTaskStatus(s string) *TaskStatus {
	status := TaskStatus(s)
	if !status.IsValid() {
		return nil
	}
	return &status
}
