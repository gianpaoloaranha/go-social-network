package model

type UserFollow struct {
	FollowerID string `gorm:"primaryKey;type:uuid"`
	FolloweeID string `gorm:"primaryKey;type:uuid"`

	Follower User `gorm:"foreignKey:FollowerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Followee User `gorm:"foreignKey:FolloweeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (UserFollow) TableName() string {
	return "user_follows"
}
