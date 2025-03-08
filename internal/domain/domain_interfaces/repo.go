package domain_interfaces

import "time"

type RepoService interface {
	Save(id, token string, expire time.Time) error
	Get(id string) (string, error)
	Delete(id string) error
	Close() error
}
