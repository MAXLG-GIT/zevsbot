package middleware

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/rs/zerolog/log"
	"zevsbot/internal/application/app_interfaces"
)

type MWare struct {
	useCases app_interfaces.UseCases
}

//func Init(useCases app_interfaces.UseCases) *MWare {
//	return &MWare{useCases: useCases}
//}

func (mw *MWare) CheckUserAuth(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {

		auth, err := mw.useCases.Um.CheckTgUserAuth(int(update.Message.From.ID))
		if err != nil {
			log.Error().Msg(err.Error())
		}
		if auth != true {
			//TODO
			// не авторизован
		}
		//TODO
		// авторизован
		next(ctx, b, update)
	}
}
