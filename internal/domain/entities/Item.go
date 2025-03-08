package entities

type Item struct {
	Id            int
	Sku           string `json:"sku"`
	Name          string `json:"name"`
	Price         string `json:"price"`
	PriceDiscount string `json:"price_discount"`
	PriceVat      string `json:"price_vat"`
	Weight        string `json:"weight"`
	Multiplicity  string `json:"multiplicity"`
	Url           string `json:"url"`
	Warehouses    []Warehouse
	Images        []string
}
