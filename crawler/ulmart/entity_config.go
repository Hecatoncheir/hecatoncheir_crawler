package ulmart

import "hecatoncheir/crawler"

type ItemConfig struct {
	ItemSelector        string
	NameOfItemSelector  string
	PriceOfItemSelector string
}

type Page struct {
	ItemConfig
	CityInCookieKey               string
	CityID                        string
	Path                          string
	TotalCountItemsOnPageSelector string
	MaxItemsOnPageSelector        string
	PagePath                      string
	PageParamPath                 string
}

type EntityConfig struct {
	crawler.Company
	Pages []Page
}
