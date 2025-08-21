package models

import (
	"time"

	"gorm.io/gorm"
)

type Phone struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Number    string         `gorm:"not null" json:"number"`
	ContactID uint           `gorm:"not null" json:"contactId"`
	Contact   Contact        `gorm:"foreignKey:ContactID" json:"-"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
