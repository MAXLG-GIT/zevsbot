package botutility

import (
	"fmt"
	"github.com/bojanz/currency"
	"github.com/go-telegram/bot/models"
	"os"
	"strings"
	"zevsbot/internal/domain/entities"
	"zevsbot/internal/infrastructure/messages"
)

func RenderMediaGroupPhotos(item *entities.Item) ([]models.InputMedia, error) {
	var resMedia []models.InputMedia

	warehouses := renderWarehouses(item.Warehouses)
	caption := fmt.Sprintf("<b>%s</b>\n<i>%s</i>\n%s: %s, %s: %s\n%s: %s\n%s\n<a href=\"%s\">%s</a>",
		item.Sku, item.Name, messages.GetMessage("price"), renderPrice(item.PriceDiscount),
		messages.GetMessage("with_nds"), renderPrice(item.PriceVat),
		messages.GetMessage("weight"), item.Weight, strings.Join(warehouses, ", "), item.Url,
		messages.GetMessage("website_link"),
	)

	if len(item.Images) < 1 {
		resMedia = append(resMedia, &models.InputMediaPhoto{
			Media:     os.Getenv("ZEVS_BLANK_PHOTO_IMG"),
			Caption:   caption,
			ParseMode: "HTML",
		})
		return resMedia, nil
	}

	resMedia = append(resMedia, &models.InputMediaPhoto{
		Caption:   caption,
		ParseMode: "HTML",
		Media:     item.Images[0],
	})
	for i := 1; i < len(item.Images); i++ {
		resMedia = append(resMedia, &models.InputMediaPhoto{
			Media: item.Images[i],
		})
	}

	return resMedia, nil
}

func renderWarehouses(warehouses []entities.Warehouse) []string {
	//
	var resStr []string
	for _, val := range warehouses {
		var whName string
		if whName = messages.GetMessage(val.Warehouse); whName == "" {
			return resStr
		}
		resStr = append(resStr, fmt.Sprintf("%s: %sшт", whName, val.Quantity))
	}
	return resStr
}

func renderPrice(price string) string {

	amount, _ := currency.NewAmount(price, os.Getenv("ZEVS_TG_CURRENCY"))
	locale := currency.NewLocale(os.Getenv("ZEVS_TG_CURRENCY_CODE"))
	formatter := currency.NewFormatter(locale)
	return formatter.Format(amount)
}
