package models

import "time"

type Project struct {
	ID        string `gorm:"primaryKey;size:40"`
	Name      string `gorm:"size:50;not null"`
	Desc      string
	CreatedBy string `gorm:"size:40;not null"`
	UpdatedBy string `gorm:"size:40;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
