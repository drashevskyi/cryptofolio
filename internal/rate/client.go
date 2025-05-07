package rate

import (
	"encoding/json"
	"net/http"
	"time"

	"cryptofolio/internal/config"
)

func FetchRates() (map[string]float64, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(config.CoinGeckoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]map[string]float64
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return map[string]float64{
		"BTC": data["bitcoin"]["usd"],
		"ETH": data["ethereum"]["usd"],
		"LTC": data["litecoin"]["usd"],
	}, nil
}
