package services

import "zevsbot/internal/domain/entities"

type ZevsApi interface {
	Auth(email string, pass []byte) (*entities.User, error)
	Logout(email string) error
	SearchRemote(token, target string) ([]entities.Item, error)
}
