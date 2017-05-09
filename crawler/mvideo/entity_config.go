package mvideo

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

type Company struct {
	ID         string
	Iri        string
	Name       string
	Categories []string
}

type EntityConfig struct {
	Company
	Pages []Page
}
