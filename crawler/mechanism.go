package crawler

import (
	"regexp"
	//"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Item is a structure of one product from one page
type Item struct {
	Name     string
	Price    string
	Company  Company
	DateTime time.Time
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

		pageItem := Item{
			Name:     name,
			Price:    price,
			DateTime: time.Now().UTC(),
			Company:  company,
		}

		crawler.Items <- pageItem
	})

	return nil
}

// RunWithConfiguration can parse web documents and make Item structure for each product on page filtered by selectors
func (crawler *Crawler) RunWithConfiguration(config EntityConfig) error {
	patternForCutPrice, _ := regexp.Compile("р[уб]*?.")

	for _, pageConfig := range config.Pages {

		document, err := goquery.NewDocument(config.Company.Iri + pageConfig.Path + pageConfig.PageParamPath + "1")
		if err != nil {
			return err
		}



		go crawler.GetItemsFromPage(document, pageConfig, config.Company, patternForCutPrice)

		//pagesCount := document.Find(pageConfig.PageInPaginationSelector).Last().Find("a").Text()

		//countOfPages, err := strconv.Atoi(pagesCount)
		//if err != nil {
		//	return err
		//}

		// var iterator int
		// for iterator = 2; iterator <= countOfPages; iterator++ {

		// 	document, err := goquery.NewDocument(config.Company.Iri + pageConfig.Path + pageConfig.PageParamPath + strconv.Itoa(iterator))
		// 	if err != nil {
		// 		return err
		// 	}

		// 	go crawler.GetItemsFromPage(document, pageConfig, config.Company, patternForCutPrice)
		// }
	}

	return nil
}
