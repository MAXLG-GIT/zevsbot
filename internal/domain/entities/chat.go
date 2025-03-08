package entities

import (
	"time"
)

type Chat struct {
	Id       int
	TgId     int
	UserName string
	Token    string
	TokenExp time.Time
}
