package bothandlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/rs/zerolog/log"
	"image"
	"os"
	"zevsbot/internal/application/app_interfaces"
	"zevsbot/internal/infrastructure/config"
	"zevsbot/internal/infrastructure/messages"
	"zevsbot/internal/infrastructure/tgbot/botutility"
	"zevsbot/internal/infrastructure/utilities"
	"zevsbot/internal/zevs_errors"
)

type Handlers interface {
	ImageHandler(ctx context.Context, b *bot.Bot, update *models.Update) error
	LoginHandler(ctx context.Context, b *bot.Bot, update *models.Update) error
	LogoutHandler(ctx context.Context, b *bot.Bot, update *models.Update) error
	TextHandler(ctx context.Context, b *bot.Bot, update *models.Update) error
}

type TgHandlers struct {
	useCases app_interfaces.UseCases
}

func InitTgHandlers(useCases app_interfaces.UseCases) Handlers {
	return &TgHandlers{useCases: useCases}
}

func (tgh TgHandlers) LoginHandler(ctx context.Context, b *bot.Bot, update *models.Update) error {

	replyToMessage := update.Message.ReplyToMessage
	if replyToMessage != nil {
		if replyToMessage.Entities != nil {
			if len(replyToMessage.Entities) > 1 && replyToMessage.Entities[1].Type == "email" {

				userEmail, err := utilities.GetSubstring(
					replyToMessage.Text,
					replyToMessage.Entities[1].Offset,
					replyToMessage.Entities[1].Length)
				if err != nil {
					botutility.SendTextMessage(ctx, b, update.Message.Chat.ID,
						messages.GetMessage("error_occurred"), "")
					return err
				}
				userPassword := update.Message.Text

				err = tgh.useCases.Um.Auth(int(update.Message.Chat.ID), userEmail, []byte(userPassword))
				if err != nil {
					var zevsAuthError *zevs_errors.PublicError
					if errors.As(err, &zevsAuthError) {
						botutility.SendTextMessage(ctx, b, update.Message.Chat.ID,
							err.Error(), "")
						return err
					}
					botutility.SendTextMessage(ctx, b, update.Message.Chat.ID,
						messages.GetMessage("authorization_failed"), "")
					return err
				}

				botutility.SendTextMessage(ctx, b, update.Message.Chat.ID, messages.GetMessage("authorized"), "")
				return nil
			}

			if update.Message.Entities != nil && len(update.Message.Entities) > 0 {
				if update.Message.Entities[0].Type == "email" {
					userEmail, err := utilities.GetSubstring(
						update.Message.Text,
						update.Message.Entities[0].Offset,
						update.Message.Entities[0].Length)
					if err != nil {
						botutility.SendTextMessage(ctx, b, update.Message.Chat.ID, messages.GetMessage("error_occurred"), "")
						return err
					}

					botutility.SendTextMessage(ctx, b, update.Message.Chat.ID,
						fmt.Sprintf(messages.GetMessage("command_auth"), userEmail),
						fmt.Sprintf("{\"force_reply\":true, \"input_field_placeholder\": \"%s\"}",
							messages.GetMessage("input_password")))
					return nil
				}
			}

		}
	}

	botutility.SendTextMessage(ctx, b, update.Message.Chat.ID,
		messages.GetMessage("command_input_email"),
		fmt.Sprintf("{\"force_reply\":true, \"input_field_placeholder\": \"%s\"}",
			messages.GetMessage("input_email")))
	return nil

}
func (tgh TgHandlers) LogoutHandler(ctx context.Context, b *bot.Bot, update *models.Update) error {
	err := tgh.useCases.Um.Logout(int(update.Message.Chat.ID))
	if err != nil {
		botutility.SendTextMessage(ctx, b, update.Message.Chat.ID,
			messages.GetMessage("error_occurred"), "")
		return err
	}
	botutility.SendTextMessage(ctx, b, update.Message.Chat.ID, messages.GetMessage("logged_out"), "")
	return nil
}
func (tgh TgHandlers) TextHandler(ctx context.Context, b *bot.Bot, update *models.Update) error {

	//TODO
	// implement this part properly
	return tgh.searchTextHandler(ctx, b, update.Message.Chat.ID, update.Message.Text)
}

func (tgh TgHandlers) ImageHandler(ctx context.Context, b *bot.Bot, update *models.Update) error {

	photo := update.Message.Photo[len(update.Message.Photo)-1]
	fileData, err := b.GetFile(ctx, &bot.GetFileParams{FileID: photo.FileID})
	if err != nil {
		//log.Error().Msg(err.Error())
		return err
	}
	file, err := utilities.DownloadFile(os.Getenv("ZEVS_TG_FILEPATH") + fileData.FilePath)

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
			log.Error().Msg(err.Error())
		}
		err = os.Remove(file.Name())
		if err != nil {
			fmt.Println(err)
			log.Error().Msg(err.Error())
		}
	}(file)

	img, _, err := image.Decode(file)

	if err != nil {

		botutility.SendTextMessage(ctx, b, update.Message.Chat.ID,
			messages.GetMessage("error_occurred"), "")
		return nil
	}

	searchString, err := tgh.useCases.Sm.SearchImageText(&img)
	if err != nil {
		return err
	}

	if searchString == "" {
		botutility.SendTextMessage(ctx, b, update.Message.Chat.ID,
			messages.GetMessage("bar_code_not_read"), "")
		return nil
	}
	botutility.SendTextMessage(ctx, b, update.Message.Chat.ID,
		fmt.Sprintf("%s %s", messages.GetMessage("currently_searching"), searchString), "")
	return tgh.searchTextHandler(ctx, b, update.Message.Chat.ID, searchString)

}

func (tgh TgHandlers) searchTextHandler(ctx context.Context, b *bot.Bot, chatId int64, searchString string) error {

	items, err := tgh.useCases.Sm.Search(int(chatId), searchString)
	if err != nil {
		var zevsAuthError *zevs_errors.PublicError
		if errors.As(err, &zevsAuthError) {
			botutility.SendTextMessage(ctx, b, chatId,
				err.Error(), "")
			return err
		}
		botutility.SendTextMessage(ctx, b, chatId,
			messages.GetMessage("error_occurred"), "")
		return err
	}

	if items == nil {
		botutility.SendTextMessage(ctx, b, chatId,
			messages.GetMessage("nothing_found"), "")
		return nil
	}

	var listCount int
	if listCount = config.ReadIntVal("ZEVS_TG_LIST_COUNT", 3); listCount > len(items) {
		listCount = len(items)
	}

	for i := 0; i < listCount; i++ {
		printableMedia, renderErr := botutility.RenderMediaGroupPhotos(&items[i])
		if renderErr != nil {
			botutility.SendTextMessage(ctx, b, chatId,
				messages.GetMessage("error_occurred"), "")
			return renderErr
		}
		_, sendErr := b.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
			ChatID: chatId,
			Media:  printableMedia,
		})
		if sendErr != nil {
			botutility.SendTextMessage(ctx, b, chatId,
				messages.GetMessage("error_occurred"), "")
			return sendErr
		}
	}

	return nil
}
