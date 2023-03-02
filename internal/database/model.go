package database

import "time"

type User struct {
	ID         string    `bson:"_id" json:"id"`
	UserName   string    `bson:"userName" json:"userName"`
	Points     int       `bson:"points" json:"points"`
	JoinedDate time.Time `bson:"joinedDate" json:"joinedDate"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time `bson:"updatedAt" json:"updatedAt"`
}
