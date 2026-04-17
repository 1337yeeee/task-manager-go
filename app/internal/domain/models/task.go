package models

import "time"

type Task struct {
	ID          string     `gorm:"primaryKey;size:40" json:"id"`
	ProjectID   string     `gorm:"size:40;not null" json:"project_id"`
	Name        string     `gorm:"size:50;not null" json:"name"`
	Content     string     `json:"content"`
	ExecutiveID string     `gorm:"size:40;not null" json:"executive_id"`
	AuditorID   string     `gorm:"size:40;" json:"auditor_id"`
	CreatedBy   string     `gorm:"size:40;" json:"created_by"`
	UpdatedBy   string     `gorm:"size:40;" json:"updated_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Status      TaskStatus `gorm:"size:50;not null" json:"status"`
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

func (s TaskStatus) String() string {
	return string(s)
}

func ParseTaskStatus(s string) *TaskStatus {
	status := TaskStatus(s)
	if !status.IsValid() {
		return nil
	}
	return &status
}
