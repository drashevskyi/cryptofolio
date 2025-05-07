package config

var StaticUsers = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

const CoinGeckoURL = "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum,litecoin&vs_currencies=usd"
