package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"cryptofolio/internal/auth"
	"cryptofolio/internal/config"
	"cryptofolio/internal/model"
	"cryptofolio/internal/rate"

	"github.com/gorilla/mux"
)

func getUser(r *http.Request) (string, bool) {
	return auth.GetUserFromContext(r)
}

func CreateAsset(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var asset model.Asset
		if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !config.ValidCurrencies[asset.Currency] {
			http.Error(w, "Invalid currency. Only BTC, ETH, LTC are supported.", http.StatusBadRequest)
			return
		}

		user, ok := getUser(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		asset.User = user

		_, err := db.Exec(`INSERT INTO assets("user", label, currency, amount) VALUES($1, $2, $3, $4)`,
			asset.User, asset.Label, asset.Currency, asset.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func ListAssets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := getUser(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`SELECT id, label, currency, amount FROM assets WHERE "user" = $1`, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var assets []model.Asset
		for rows.Next() {
			var asset model.Asset
			if err := rows.Scan(&asset.ID, &asset.Label, &asset.Currency, &asset.Amount); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			assets = append(assets, asset)
		}

		json.NewEncoder(w).Encode(assets)
	}
}

func GetAsset(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := getUser(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		id := mux.Vars(r)["id"]

		var asset model.Asset
		err := db.QueryRow(`SELECT id, label, currency, amount FROM assets WHERE "user" = $1 AND id = $2`, user, id).
			Scan(&asset.ID, &asset.Label, &asset.Currency, &asset.Amount)

		if err == sql.ErrNoRows {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rates, err := rate.FetchRates()
		if err != nil {
			http.Error(w, "Rate fetch failed", http.StatusBadGateway)
			return
		}

		usdValue := asset.Amount * rates[asset.Currency]

		response := map[string]interface{}{
			"id":        asset.ID,
			"label":     asset.Label,
			"currency":  asset.Currency,
			"amount":    asset.Amount,
			"usd_value": usdValue,
		}

		json.NewEncoder(w).Encode(response)
	}
}

func UpdateAsset(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := getUser(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		id := mux.Vars(r)["id"]
		var asset model.Asset
		if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !config.ValidCurrencies[asset.Currency] {
			http.Error(w, "Invalid currency. Only BTC, ETH, LTC are supported.", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(`UPDATE assets SET label=$1, currency=$2, amount=$3 WHERE id=$4 AND "user"=$5`,
			asset.Label, asset.Currency, asset.Amount, id, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteAsset(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := getUser(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		id := mux.Vars(r)["id"]

		_, err := db.Exec(`DELETE FROM assets WHERE id = $1 AND "user" = $2`, id, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func TotalValueUSD(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := getUser(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		rows, err := db.Query(`SELECT currency, amount FROM assets WHERE "user" = $1`, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		rates, err := rate.FetchRates()
		if err != nil {
			http.Error(w, "Rate fetch failed", http.StatusBadGateway)
			return
		}

		total := 0.0
		for rows.Next() {
			var currency string
			var amount float64
			_ = rows.Scan(&currency, &amount)
			total += amount * rates[currency]
		}

		json.NewEncoder(w).Encode(map[string]float64{"total_usd": total})
	}
}
