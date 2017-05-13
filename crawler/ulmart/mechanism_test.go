package mvideo

import (
	"hecatoncheir/crawler"
	"testing"
	"time"
	"net/http"
	"net/http/cookiejar"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

func TestCookies(test *testing.T) {
	iri := "https://www.ulmart.ru/catalog/communicators?sort=5&viewType=2&rec=true"

	cookie, _ := cookiejar.New(nil)
	// 18414 - Москва
	// 1688 - Алексин
	city := &http.Cookie{Name: "city", Value: "18414"}
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
		test.Log(err)
	}

	cityName := document.Find("#load-cities").Text()
	cityName = strings.TrimSpace(cityName)

	if cityName != "Москва" {
		test.Fail()
	}
}

func TestCrawlerCanGetDocymentByConfig(test *testing.T) {
	smartphonesPage := Page{
		Path:                     "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		ItemConfig: ItemConfig{
			ItemSelector:        ".grid-view .product-tile",
			NameOfItemSelector:  ".product-tile-title",
			PriceOfItemSelector: ".product-price-current",
		},
	}

	configuration := EntityConfig{
		Company: crawler.Company{
			Iri:        "http://www.mvideo.ru/",
			Name:       "M.Video",
			Categories: []string{"Телефоны"},
		},
		Pages: []Page{smartphonesPage},
	}

	mechanism := NewCrawler()

	go mechanism.RunWithConfiguration(configuration)

	isRightItems := false

	go func() {
		time.Sleep(time.Second * 3)
		close(mechanism.Items)
	}()

	for item := range mechanism.Items {
		if item.Name != "" && item.Price.Value != "" {
			isRightItems = true
			break
		}
	}

	if isRightItems == false {
		test.Fail()
	}
}
