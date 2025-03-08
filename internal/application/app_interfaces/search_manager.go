package app_interfaces

import (
	"image"
	"zevsbot/internal/domain/entities"
)

type SearchManager interface {
	Search(chatId int, phrase string) ([]entities.Item, error)
	SearchImageText(img *image.Image) (string, error)
}
