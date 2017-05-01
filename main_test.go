package main

import (
	"hecatonhair/crawler"
	socket "hecatonhair/socket_engine"
	"sync"
	"testing"

	"golang.org/x/net/websocket"
)

var (
	once       sync.Once
	goroutines sync.WaitGroup
)

func SetUpSocketServer() {
	server := socket.NewEngine("v1.0")
	goroutines.Done()
	server.Listen("localhost", 8181)
	defer server.Server.Close()
}

// {
// 	"Message": "Get items from categories of company",
// 	"Data": {
// 			"Iri": "http://www.mvideo.ru/",
//			"Name": "M.Video",
//			"Categories": ["Телефоны"],
// 			"Pages": [{
// 				"Path": "smartfony-i-svyaz/smartfony-205",
// 				"PageInPaginationSelector": ".pagination-list .pagination-item",
// 				"PageParamPath": "/f/page=",
// 				"ItemSelector": ".grid-view .product-tile",
// 				"NameOfItemSelector": ".product-tile-title",
// 				"PriceOfItemSelector": ".product-price-current"
// 			}]
// 	}
// }
func TestSocketCanParseDocumentOfEntity(test *testing.T) {
	goroutines.Add(1)
	go once.Do(SetUpSocketServer)
	goroutines.Wait()

	client, err := websocket.Dial("ws://localhost:8181", "", "http://localhost:8181")

	if err != nil {
		test.Fatal()
	}

	smartphonesPage := crawler.Page{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		ItemConfig: crawler.ItemConfig{
			ItemSelector:        ".grid-view .product-tile",
			NameOfItemSelector:  ".product-tile-title",
			PriceOfItemSelector: ".product-price-current",
		},
	}

	configuration := crawler.EntityConfig{
		Company: crawler.Company{
			Iri:        "http://www.mvideo.ru/",
			Name:       "M.Video",
			Categories: []string{"Телефоны"},
		},
		Pages: []crawler.Page{smartphonesPage},
	}

	websocket.JSON.Send(client, socket.MessageEvent{Message: "Get items from categories of company", Data: configuration})

	message := make(chan socket.MessageEvent)

	go func() {
		for {
			socketEvent := socket.MessageEvent{}
			err := websocket.JSON.Receive(client, &socketEvent)
			if err != nil {
				test.Error(err)
			}
			message <- socketEvent
			break
		}
	}()

	for event := range message {
		if event.Message != "Item from categories of company parsed" ||
			event.Data.(map[string]interface{})["Item"] == nil {
			test.Fail()
		}
		break
	}

}
