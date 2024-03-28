package models

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey;Column:id;type:varchar(45)" json:"id"`
	Name      string    `gorm:"Column:name;type:varchar(255)" json:"name"`
	CreatedAt time.Time `gorm:"Column:created_at;type:timestamp;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"Column:updated_at;type:timestamp;not null" json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}
