package models

import "time"

type Project struct {
	ID          string    `gorm:"primaryKey;size:40" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Description string    `gorm:"size:255;" json:"description"`
	CreatedBy   string    `gorm:"size:40;not null" json:"created_by"`
	UpdatedBy   string    `gorm:"size:40;not null" json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
