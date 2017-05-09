package mvideo

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"hecatoncheir/crawler"

	"github.com/PuerkitoBio/goquery"
)

var cities = crawler.Cities{
	"Москва":      "CityCZ_975",
	"Новосибирск": "CityCZ_2246",
}

// Crawler for parse documents
type Crawler struct {
	Items chan crawler.Item // For subscribe to events
}

// NewCrawler create a new Crawler object
func NewCrawler() *Crawler {
	crawler := Crawler{Items: make(chan crawler.Item)}
	return &crawler
}

// GetItemsFromPage can get product from html document by selectors in the configuration
func (cw *Crawler) GetItemsFromPage(document *goquery.Document, pageConfig Page, company crawler.Company, patternForCutPrice *regexp.Regexp) error {
	document.Find(pageConfig.ItemSelector).Each(func(iterator int, item *goquery.Selection) {
		var name, price string

		name = item.Find(pageConfig.NameOfItemSelector).Text()
		price = item.Find(pageConfig.PriceOfItemSelector).Text()

		name = strings.TrimSpace(name)
		price = strings.TrimSpace(price)

		// price = strings.Replace(price, "р.", "", -1)
		price = patternForCutPrice.ReplaceAllString(price, "")

		//fmt.Printf("Review %s: %s \n", name, price)

		cityName, err := cities.SearchCityByCode(pageConfig.CityParam)
		if err != nil {
			log.Println(err)
		}

		priceData := crawler.Price{
			Value:    price,
			City:     cityName,
			DateTime: time.Now().UTC(),
		}

		pageItem := crawler.Item{
			Name:    name,
			Price:   priceData,
			Company: company,
		}

		cw.Items <- pageItem
	})

	return nil
}

// RunWithConfiguration can parse web documents and make Item structure for each product on page filtered by selectors
func (cw *Crawler) RunWithConfiguration(config EntityConfig) error {
	patternForCutPrice, _ := regexp.Compile("р[уб]*?.")

	for _, pageConfig := range config.Pages {

		document, err := goquery.NewDocument(config.Company.Iri + pageConfig.Path + pageConfig.PageParamPath + "1" + pageConfig.CityParamPath + pageConfig.CityParam)
		if err != nil {
			return err
		}

		go cw.GetItemsFromPage(document, pageConfig, config.Company, patternForCutPrice)

		pagesCount := document.Find(pageConfig.PageInPaginationSelector).Last().Find("a").Text()

		countOfPages, err := strconv.Atoi(pagesCount)
		if err != nil {
			return err
		}

		pagesCrawling := make(chan func(), 6)

		go func() {
			for crawler := range pagesCrawling {
				go crawler()
			}
		}()

		var iterator int
		for iterator = 2; iterator <= countOfPages; iterator++ {
			document, err := goquery.NewDocument(config.Company.Iri + pageConfig.Path + pageConfig.PageParamPath + strconv.Itoa(iterator) + pageConfig.CityParamPath + pageConfig.CityParam)
			if err != nil {
				return err
			}

			pagesCrawling <- func() {
				cw.GetItemsFromPage(document, pageConfig, config.Company, patternForCutPrice)
			}
		}

		close(pagesCrawling)
	}

	return nil
}
