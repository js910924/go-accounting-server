package models

import "time"

type Data struct {
	LogId       int
	UserId      int
	ActionType  int
	DetailType  int
	Money       int
	Description string
	CreateTime  time.Time
}
