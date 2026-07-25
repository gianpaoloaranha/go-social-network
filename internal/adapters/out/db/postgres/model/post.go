package model

import "time"

type Post struct {
	ID          string    `gorm:"primaryKey;type:uuid"`
	AuthorID    string    `gorm:"type:uuid;not null;index"`
	Description string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`

	Author   User      `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}

func (Post) TableName() string {
	return "posts"
}
