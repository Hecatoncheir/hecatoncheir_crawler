# Hecatoncheir
Crawler with websocket and rest api

By default tcp server run on `8181` port.

<br>

> #### TODO:
> [0] Log <br>
> [0] REST POST method

## REST API:

[✓] GET `api/version` <br>
Response json: `{"apiVersion":"v1.0"}` 

## Socket
<br>
Send message:

```
{"Message":"Need api version"}
```
Response:

```
{"Message": "Version of API", "Data": {"API version": "v1.0"}
```
 
---
Send message:

```json
 {
 	"Message": "Get items from categories of company",
 	"Data": {
			"Iri": "http://link of company",
			"Name": "Name of company",
			"Categories": ["Some categories id or name"],
 			"Pages": [{
 				"Path": "path to search page",
 				"PageInPaginationSelector": ".pagination-list .pagination-item",
 				"PageParamPath": "/f/page=",
 				"CityParamPath": "?cityId=",
 				"CityParam": "CityCZ_975",
 				"ItemSelector": ".grid-view .product-tile",
 				"NameOfItemSelector": ".product-tile-title",
 				"PriceOfItemSelector": ".product-price-current"
 			}]
 	}
 }
```

Response for all connected clients:
```json
{
	"Data": {
		"Item": {
			"Name": "Смартфон Samsung Galaxy J5 Prime Black",
			"Price": {
 				"Value": "12990",
 				"DateTime": "2017-05-01T16:27:18.543653798Z",
 				"City": "Москва"
			},
			"Company": {
				"ID": "",
				"Iri": "link",
				"Name": "Company name",
				"Categories": ["Some categories id or name"]
			},
		}
	},
	"Message": "Item from categories of company parsed"
}
```
