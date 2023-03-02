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

type Activity struct {
	User      string    `json:"user" bson:"user" required:"true"`
	UserName  string    `json:"userName" bson:"userName"`
	ChannelId string    `json:"channelId" bson:"channelId" required:"true"`
	Activity  string    `json:"activity" bson:"activity" required:"true" enum:"attend,react,receive,play"`
	Reward    int       `json:"reward" bson:"reward" required:"true" enum:"-10,5,10,50"`
	MessageId string    `json:"messageId" bson:"messageId"`
	Emoji     string    `json:"emoji" bson:"emoji"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
