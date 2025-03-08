package tgbot

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"zevsbot/internal/application/app_interfaces"
	"zevsbot/internal/infrastructure/middleware"
	"zevsbot/internal/infrastructure/tgbot/bothandlers"

	"zevsbot/internal/infrastructure/utilities"
)

const (
	loginCommand  = "/auth"
	logoutCommand = "/logout"
)

type BotStruct struct {
	useCases app_interfaces.UseCases
	handlers bothandlers.Handlers
	mWare    *middleware.MWare
	//bot      *bot.Bot
}

func Init(useCases app_interfaces.UseCases, mWare *middleware.MWare) *BotStruct {
	handlers := bothandlers.InitTgHandlers(useCases)
	return &BotStruct{useCases: useCases, handlers: handlers, mWare: mWare}
}

func (bs BotStruct) Run(ctx context.Context) error {

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		//bot.WithMiddlewares(showMessageWithUserID, bs.mWare.CheckUserAuth),
		bot.WithDefaultHandler(bs.routerHandler),
	}
	b, err := bot.New(os.Getenv("ZEVS_TG_TOKEN"), opts...)
	if err != nil {
		return err
	}
	//bs.bot = b

	b.Start(ctx)
	return nil
}

//func showMessageWithUserID(next bot.HandlerFunc) bot.HandlerFunc {
//	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
//		if update.Message != nil {
//			log.Printf("%d say: %s", update.Message.From.ID, update.Message.Text)
//		}
//		next(ctx, b, update)
//	}
//}

func (bs BotStruct) routerHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	if update.Message == nil {
		return
	}
	if update.Message.Photo != nil && len(update.Message.Photo) > 0 {

		go func() {
			err := bs.handlers.ImageHandler(ctx, b, update)
			if err != nil {
				log.Error().Msg(err.Error())
				return
			}
		}()
		return
	}

	replyToMessage := update.Message.ReplyToMessage
	if replyToMessage != nil {
		if len(replyToMessage.Entities) > 0 && replyToMessage.Entities[0].Type == "bot_command" {
			bs.routeCommand(ctx, b, update, replyToMessage)
			return
		}
	}

	if len(update.Message.Entities) > 0 && update.Message.Entities[0].Type == "bot_command" {

		bs.routeCommand(ctx, b, update, update.Message)
		return
	}

	go func() {
		err := bs.handlers.TextHandler(ctx, b, update)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
	}()

}

func (bs BotStruct) routeCommand(ctx context.Context, b *bot.Bot, update *models.Update, message *models.Message) {
	replyToEntity := message.Entities[0]

	command, err := utilities.GetSubstring(message.Text, replyToEntity.Offset, replyToEntity.Length)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	switch command {
	case loginCommand:
		go func() {
			loginErr := bs.handlers.LoginHandler(ctx, b, update)
			if loginErr != nil {
				log.Error().Msg(loginErr.Error())
				return
			}
		}()
	case logoutCommand:
		go func() {
			logoutErr := bs.handlers.LogoutHandler(ctx, b, update)
			if logoutErr != nil {
				log.Error().Msg(logoutErr.Error())
				return
			}
		}()
	default:

	}

	log.Error().Msg("unknown command")
	return
}
