package ulmart

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"hecatoncheir/crawler"
)

var cities = crawler.Cities{
	"Москва":  "18414",
	"Алексин": "1688",
}

// Crawler for parse documents
type Crawler struct {
	Items chan crawler.Item // For subscribe to events
}

// NewCrawler create a new Crawler object
func NewCrawler() *Crawler {
	newCrawler := Crawler{Items: make(chan crawler.Item)}
	return &newCrawler
}

// GetItemsFromPage can get product from html document by selectors in the configuration
func (cw *Crawler) GetItemsFromPage(document *goquery.Document, pageConfig Page, company crawler.Company) error {
	document.Find(pageConfig.ItemSelector).Each(func(iterator int, item *goquery.Selection) {
		var name, price string

		name = item.Find(pageConfig.NameOfItemSelector).Text()
		price = item.Find(pageConfig.PriceOfItemSelector).Text()

		name = strings.TrimSpace(name)

		price = strings.Replace(price, " ", "", -1)
		price = strings.TrimSpace(price)

		cityName, err := cities.SearchCityByCode(pageConfig.CityID)
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

func (cw *Crawler) GetDocumentForUrl(iri string, pageConfig Page) (*goquery.Document, error) {
	cookie, _ := cookiejar.New(nil)
	city := &http.Cookie{Name: pageConfig.CityInCookieKey, Value: pageConfig.CityID}
	allCookies := []*http.Cookie{}
	allCookies = append(allCookies, city)

	pageUrl, _ := url.Parse(iri)
	cookie.SetCookies(pageUrl, allCookies)
	client := &http.Client{
		Jar: cookie,
	}

	request, _ := http.NewRequest("GET", iri, nil)
	response, _ := client.Do(request)

	document, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}

	return document, nil
}

// RunWithConfiguration can parse web documents and make Item structure for each product on page filtered by selectors
func (cw *Crawler) RunWithConfiguration(config EntityConfig) error {

	for _, pageConfig := range config.Pages {

		iri := config.Company.Iri + pageConfig.Path

		document, err := cw.GetDocumentForUrl(iri, pageConfig)
		if err != nil {
			return err
		}

		totalPerPageItems, err := strconv.Atoi(document.Find(pageConfig.TotalCountItemsOnPageSelector).Text())
		if err != nil {
			return err
		}

		maxItems, err := strconv.Atoi(document.Find(pageConfig.MaxItemsOnPageSelector).Text())
		if err != nil {
			return err
		}

		countOfPages := maxItems / totalPerPageItems

		if maxItems%totalPerPageItems != 0 {
			countOfPages += 1
		}

		pagesCrawling := make(chan func(), 6)

		go func() {
			for someCrawler := range pagesCrawling {
				go someCrawler()
			}
		}()

		var iterator int
		for iterator = 1; iterator <= countOfPages; iterator++ {
			iri := config.Company.Iri + pageConfig.PagePath + pageConfig.PageParamPath + strconv.Itoa(iterator)
			document, err := cw.GetDocumentForUrl(iri, pageConfig)
			if err != nil {
				return err
			}

			pagesCrawling <- func() {
				cw.GetItemsFromPage(document, pageConfig, config.Company)
			}
		}

		close(pagesCrawling)
	}

	return nil
}
