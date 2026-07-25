package model

import "time"

type Comment struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	AuthorID  string    `gorm:"type:uuid;not null;index"`
	PostID    string    `gorm:"type:uuid;not null;index"`
	Message   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`

	Author User `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Post   Post `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Comment) TableName() string {
	return "comments"
}
