package model

type User struct {
	ID    string    `gorm:"primaryKey;type:uuid"`
	Name  string    `gorm:"not null"`
	Email string    `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`

	Posts     []Post    `gorm:"foreignKey:AuthorID"`
	Comments  []Comment `gorm:"foreignKey:AuthorID"`
	Following []User    `gorm:"many2many:user_follows;joinForeignKey:FollowerID;joinReferences:FolloweeID"`
	Followers []User    `gorm:"many2many:user_follows;joinForeignKey:FolloweeID;joinReferences:FollowerID"`
}

func (User) TableName() string {
	return "users"
}
