package usecases

import (
	"image"
	"strconv"
	"zevsbot/internal/application/app_interfaces"
	"zevsbot/internal/application/services"
	"zevsbot/internal/domain/domain_interfaces"
	"zevsbot/internal/domain/entities"
	"zevsbot/internal/infrastructure/messages"
	"zevsbot/internal/zevs_errors"
)

type searchManager struct {
	barcodeReader domain_interfaces.Barcode
	zevsApi       services.ZevsApi
	repo          domain_interfaces.RepoService
}

func InitSearchManager(barcodeReader domain_interfaces.Barcode,
	repo domain_interfaces.RepoService,
	api services.ZevsApi) app_interfaces.SearchManager {

	return &searchManager{barcodeReader: barcodeReader, zevsApi: api, repo: repo}
}

func (sm searchManager) Search(chatId int, query string) ([]entities.Item, error) {
	res := make([]entities.Item, 0)
	token, err := sm.repo.Get(strconv.Itoa(chatId))
	if err != nil {
		if token == "" {
			return nil, &zevs_errors.PublicError{Text: messages.GetMessage("unauthorized")}
		}
		return nil, err
	}

	res, err = sm.zevsApi.SearchRemote(token, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (sm searchManager) SearchImageText(img *image.Image) (string, error) {

	imageText, err := sm.barcodeReader.ReadImage(img)
	if err != nil {
		return "", err
	}

	return imageText, nil
}
