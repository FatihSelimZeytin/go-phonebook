package models

import "time"

type Contact struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FirstName string    `gorm:"not null" json:"firstName"`
	Surname   string    `gorm:"not null" json:"surname"`
	Company   string    `json:"company"` // Nullable
	Phones    []Phone   `gorm:"foreignKey:ContactID" json:"phones"`
	UserID    uint      `gorm:"not null" json:"userId"`
	Status    bool      `gorm:"not null;default:true" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
