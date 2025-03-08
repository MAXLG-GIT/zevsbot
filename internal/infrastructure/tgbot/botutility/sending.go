package botutility

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/rs/zerolog/log"
)

func SendTextMessage(ctx context.Context, b *bot.Bot, chatID int64, text, replyMarkup string) {
	resp := &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: replyMarkup,
	}
	_, err := b.SendMessage(ctx, resp)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
}
