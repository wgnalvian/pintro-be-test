package entity

import "time"

type Transaction struct {
	ID       string    `bson:"_id" json:"id"`
	UserId   string    `bson:"user_id" json:"user_id"`
	Amount   int       `bson:"amount" json:"amount"`
	Desc     string    `bson:"desc" json:"desc"`
	From     string    `bson:"from" json:"from"`
	Type     int64     `bson:"type" json:"type"`
	CreateAt time.Time `bson:"create_at" json:"create_at"`
}
