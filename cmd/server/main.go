package main

import (
	"log"
	"net/http"

	"cryptofolio/internal/auth"
	"cryptofolio/internal/handler"
	"cryptofolio/internal/store"

	"github.com/gorilla/mux"
)

func main() {
	db := store.InitDB()
	store.Migrate(db)

	r := mux.NewRouter()

	r.HandleFunc("/login", handler.Login).Methods("POST")

	secured := r.NewRoute().Subrouter()
	secured.Use(auth.JWTAuthMiddleware)

	secured.HandleFunc("/assets", handler.CreateAsset(db)).Methods("POST")
	secured.HandleFunc("/assets", handler.ListAssets(db)).Methods("GET")
	secured.HandleFunc("/assets/{id}", handler.GetAsset(db)).Methods("GET")
	secured.HandleFunc("/assets/{id}", handler.UpdateAsset(db)).Methods("PUT")
	secured.HandleFunc("/assets/{id}", handler.DeleteAsset(db)).Methods("DELETE")
	secured.HandleFunc("/assets/value/total", handler.TotalValueUSD(db)).Methods("GET")

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
