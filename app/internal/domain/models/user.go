package models

import (
	"task-manager/internal/auth"
	"time"
)

type User struct {
	ID        string        `gorm:"primaryKey;size:40"`
	Name      string        `gorm:"size:50;not null"`
	Email     string        `gorm:"size:255;unique;not null"`
	Password  string        `gorm:"size:255;not null" json:"-"`
	Role      auth.UserRole `gorm:"size:50;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
