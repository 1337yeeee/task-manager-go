package models

import (
	"task-manager/internal/auth"
	"time"
)

type User struct {
	ID        string        `gorm:"primaryKey;size:40" json:"id"`
	Name      string        `gorm:"size:50;not null" json:"name"`
	Email     string        `gorm:"size:255;unique;not null" json:"email"`
	Password  string        `gorm:"size:255;not null" json:"-"`
	Role      auth.UserRole `gorm:"size:50;not null" json:"role"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
