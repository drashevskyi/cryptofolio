package model

import "strings"

type Asset struct {
	ID       int     `json:"id"`
	User     string  `json:"-"`
	Label    string  `json:"label"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

var AllowedCurrencies = []string{"BTC", "ETH", "LTC"}

func IsCurrencyAllowed(curr string) bool {
	for _, c := range AllowedCurrencies {
		if c == curr {
			return true
		}
	}
	return false
}

func AllowedCurrencyList() string {
	return strings.Join(AllowedCurrencies, ", ")
}
