package dto

import (
	"time"
	"zevsbot/internal/domain/entities"
)

type WebAuthResp struct {
	Data struct {
		Token   string    `json:"token"`
		Expires time.Time `json:"expires"`
		User    struct {
			Id          int       `json:"id"`
			Email       string    `json:"email"`
			Name        string    `json:"name"`
			LastName    string    `json:"last_name"`
			MiddleName  string    `json:"middle_name"`
			Phone       string    `json:"phone"`
			IsActivated bool      `json:"is_activated"`
			LastLogin   time.Time `json:"last_login"`
			IsSuperuser int       `json:"is_superuser"`
			Property    struct {
				Payer            string      `json:"'payer'"`
				PriceType        string      `json:"price-type"`
				CustomerComment  string      `json:"customer-comment"`
				CustomerWithNds  string      `json:"customer-with-nds"`
				CustomerDiscount int         `json:"customer-discount"`
				CustomerCart     interface{} `json:"customer-cart"`
				CustomerShipment interface{} `json:"customer-shipment"`
				CustomerPayment  interface{} `json:"customer-payment"`
				Keep             string      `json:"keep"`
			} `json:"property"`

			VdomahRoleId       interface{} `json:"vdomah_role_id"`
			Withnds            string      `json:"withnds"`
			Discount           int         `json:"discount"`
			Warehouse          string      `json:"warehouse"`
			Shipment           int         `json:"shipment"`
			Payment            int         `json:"payment"`
			Keep               string      `json:"keep"`
			Company            string      `json:"company"`
			Oneccomment        string      `json:"oneccomment"`
			ActiveCurrencyCode string      `json:"active_currency_code"`
			Hideprices         int         `json:"hideprices"`
			PhoneList          []string    `json:"phone_list"`
		} `json:"user"`
	} `json:"data"`
}

func (a *WebAuthResp) ToEntityUser() *entities.User {

	return &entities.User{
		Id:          a.Data.User.Id,
		Email:       a.Data.User.Email,
		IsActivated: a.Data.User.IsActivated,
		Token:       a.Data.Token,
		TokenExp:    a.Data.Expires,
	}
}
