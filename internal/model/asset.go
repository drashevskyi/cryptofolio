package model

type Asset struct {
	ID       int     `json:"id"`
	User     string  `json:"-"`
	Label    string  `json:"label"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}
