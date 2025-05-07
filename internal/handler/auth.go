package handler

import (
	"cryptofolio/internal/auth"
	"encoding/json"
	"log"
	"net/http"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !auth.Authenticate(req.Username, req.Password) {
		log.Printf("Invalid login attempt for user: %s", req.Username)
		writeError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := auth.GenerateJWT(req.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
