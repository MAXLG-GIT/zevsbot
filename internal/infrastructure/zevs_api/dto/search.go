package dto

import (
	"zevsbot/internal/domain/entities"
	"zevsbot/internal/infrastructure/utilities"
)

type WebSearchReq struct {
	Token string `json:"token"`
	Data  struct {
		Query string `json:"query"`
	} `json:"data"`
}

type WebSearchResp map[int]Item

type Warehouse struct {
	Warehouse string                   `json:"warehouse"`
	Quantity  utilities.FlexibleString `json:"quantity"`
}

type Item struct {
	Sku           string      `json:"sku"`
	Name          string      `json:"name"`
	Brand         string      `json:"brand"`
	Price         string      `json:"price"`
	PriceDiscount string      `json:"price_discount"`
	PriceVat      string      `json:"price_vat"`
	Weight        string      `json:"weight"`
	Multiplicity  string      `json:"multiplicity"`
	Url           string      `json:"url"`
	Warehouses    []Warehouse `json:"warehouses"`
	Images        []string    `json:"images"`
	Thumbnail     string      `json:"thumbnail"`
}

func (i *Item) ToEntityItem(id int) entities.Item {
	item := entities.Item{
		Id:            id,
		Sku:           i.Sku,
		Name:          i.Name,
		Price:         i.Price,
		PriceDiscount: i.PriceDiscount,
		PriceVat:      i.PriceVat,
		Weight:        i.Weight,
		Multiplicity:  i.Multiplicity,
		Url:           i.Url,
		Images:        i.Images,
	}
	for _, val := range i.Warehouses {
		item.Warehouses = append(item.Warehouses, entities.Warehouse{
			Warehouse: val.Warehouse,
			Quantity:  string(val.Quantity),
		})
	}

	return item
}
