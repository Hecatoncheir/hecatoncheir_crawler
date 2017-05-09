package crawler

import (
	"errors"
	"time"
)

// Company type for parse
type Company struct {
	ID         string
	Iri        string
	Name       string
	Categories []string
}

// Price structure
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

// Cities codes for company
type Cities map[string]string

// SearchCodeByCity method of Cities type
func (cities *Cities) SearchCodeByCity(cityName string) (string, error) {
	for city, code := range *cities {
		if city == cityName {
			return code, nil
		}
	}
	return "", errors.New("City not exist")
}

// SearchCityByCode method of Cities type
func (cities *Cities) SearchCityByCode(codeName string) (string, error) {
	for city, code := range *cities {
		if code == codeName {
			return city, nil
		}
	}
	return "", errors.New("Code not exist")
}
