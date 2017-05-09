package mvideo

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
)

type Cities map[string]string

func (cities *Cities) searchCodeByCity(cityName string) (string, error) {
	for city, code := range *cities {
		if city == cityName {
			return code, nil
		}
	}
	return "", errors.New("City not exist")
}

func (cities *Cities) searchCityByCode(codeName string) (string, error) {
	for city, code := range *cities {
		if code == codeName {
			return city, nil
		}
	}
	return "", errors.New("Code not exist")
}

var cities = Cities{
	"Москва":      "CityCZ_975",
	"Новосибирск": "CityCZ_2246",
}

type Price struct {
	Value    string
	City     string
	DateTime time.Time
}

// Item is a structure of one product from one page
type Item struct {
	Name    string
	Price   Price
	Company Company
}

// Crawler for parse documents
type Crawler struct {
	Items chan Item // For subscribe to events
}

// NewCrawler create a new Crawler object
func NewCrawler() *Crawler {
	crawler := Crawler{Items: make(chan Item)}
	return &crawler
}

// GetItemsFromPage can get product from html document by selectors in the configuration
func (crawler *Crawler) GetItemsFromPage(document *goquery.Document, pageConfig Page, company Company, patternForCutPrice *regexp.Regexp) error {
	document.Find(pageConfig.ItemSelector).Each(func(iterator int, item *goquery.Selection) {
		var name, price string

		name = item.Find(pageConfig.NameOfItemSelector).Text()
		price = item.Find(pageConfig.PriceOfItemSelector).Text()

		name = strings.TrimSpace(name)
		price = strings.TrimSpace(price)

		// price = strings.Replace(price, "р.", "", -1)
		price = patternForCutPrice.ReplaceAllString(price, "")

		//fmt.Printf("Review %s: %s \n", name, price)

		cityName, err := cities.searchCityByCode(pageConfig.CityParam)
		if err != nil {
			log.Println(err)
		}

		priceData := Price{
			Value:    price,
			City:     cityName,
			DateTime: time.Now().UTC(),
		}

		pageItem := Item{
			Name:    name,
			Price:   priceData,
			Company: company,
		}

		crawler.Items <- pageItem
	})

	return nil
}

// RunWithConfiguration can parse web documents and make Item structure for each product on page filtered by selectors
func (crawler *Crawler) RunWithConfiguration(config EntityConfig) error {
	patternForCutPrice, _ := regexp.Compile("р[уб]*?.")

	for _, pageConfig := range config.Pages {

		document, err := goquery.NewDocument(config.Company.Iri + pageConfig.Path + pageConfig.PageParamPath + "1" + pageConfig.CityParamPath + pageConfig.CityParam)
		if err != nil {
			return err
		}

		go crawler.GetItemsFromPage(document, pageConfig, config.Company, patternForCutPrice)

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
				crawler.GetItemsFromPage(document, pageConfig, config.Company, patternForCutPrice)
			}
		}

		close(pagesCrawling)
	}

	return nil
}
