package entities

import (
	"time"
)

type User struct {
	Id          int
	Email       string
	TgId        int
	IsActivated bool
	Token       string
	TokenExp    time.Time
}
