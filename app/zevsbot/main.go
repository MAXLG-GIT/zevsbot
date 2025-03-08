package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	_ "image/jpeg"
	"os"
	"zevsbot/internal/application/app_interfaces"
	"zevsbot/internal/application/usecases"
	"zevsbot/internal/domain/domain_interfaces"
	"zevsbot/internal/infrastructure/barcode"
	"zevsbot/internal/infrastructure/config"
	"zevsbot/internal/infrastructure/img_processing"
	"zevsbot/internal/infrastructure/logger"
	"zevsbot/internal/infrastructure/messages"
	"zevsbot/internal/infrastructure/middleware"
	"zevsbot/internal/infrastructure/repo"
	"zevsbot/internal/infrastructure/tgbot"
	"zevsbot/internal/infrastructure/zevs_api"
)

//TODO
// тесты авторизации |не верно указан логин, |пароль, |успешный с получением токена, пустое сообщение

func main() {

	fmt.Println("starting bot")
	config.ReadEnvConfig(".env")
	logger.InitLogger(os.Getenv("LOG_PATH"))
	messages.Init("messages.yml")

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	badgerRepo, err := repo.InitBadgerRepo(os.Getenv("ZEVS_BADGER_REPO_DIR"))
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	defer func(repo domain_interfaces.RepoService) {
		err = repo.Close()
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
	}(badgerRepo)

	imgRedactor := img_processing.Init()
	barcodeReader := barcode.Init(imgRedactor)
	zevsApi := zevs_api.Init(ctx)

	useCases := app_interfaces.UseCases{
		Um: usecases.InitUserManager(badgerRepo, zevsApi),
		Sm: usecases.InitSearchManager(barcodeReader, badgerRepo, zevsApi),
	}

	var mWare *middleware.MWare
	tgBot := tgbot.Init(useCases, mWare)

	err = tgBot.Run(ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

}
