package mvideo

import "hecatoncheir/crawler"

type ItemConfig struct {
	ItemSelector        string
	NameOfItemSelector  string
	PriceOfItemSelector string
}

type Page struct {
	ItemConfig
	CityParamPath            string
	CityParam                string
	Path                     string
	PageInPaginationSelector string
	PageParamPath            string
}

type EntityConfig struct {
	crawler.Company
	Pages []Page
}
