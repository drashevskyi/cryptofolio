package handler

import (
	"bytes"
	"cryptofolio/internal/auth"
	"cryptofolio/internal/config"
	"cryptofolio/internal/store"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret")
	config.StaticUsers = map[string]string{
		"test-user1": "password1",
		"test-user2": "password2",
	}
	os.Exit(m.Run())
}

func setupTestRouter(db *store.TestableDB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", Login).Methods("POST")

	secured := r.NewRoute().Subrouter()
	secured.Use(auth.JWTAuthMiddleware)
	secured.HandleFunc("/assets", CreateAsset(db.DB)).Methods("POST")
	return r
}

func TestLoginSuccess(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/login", Login).Methods("POST")

	body := `{"username":"test-user1","password":"password1"}`
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	var res map[string]string
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatal("Failed to decode JSON")
	}

	if res["token"] == "" {
		t.Error("expected token in response")
	}
}

func TestCreateAssetUnauthorized(t *testing.T) {
	db := store.NewTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	body := `{"label":"Test","currency":"BTC","amount":1.0}`
	req := httptest.NewRequest("POST", "/assets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", w.Code)
	}
}

func TestCreateAssetAuthorized(t *testing.T) {
	db := store.NewTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	loginPayload := `{"username":"test-user1","password":"password1"}`
	loginReq := httptest.NewRequest("POST", "/login", bytes.NewBufferString(loginPayload))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()

	router.ServeHTTP(loginRes, loginReq)

	if loginRes.Code != http.StatusOK {
		t.Fatalf("expected login 200 OK, got %d", loginRes.Code)
	}

	var loginBody map[string]string
	if err := json.NewDecoder(loginRes.Body).Decode(&loginBody); err != nil {
		t.Fatalf("failed to decode login response: %v", err)
	}
	token := loginBody["token"]
	if token == "" {
		t.Fatal("expected token in login response")
	}

	assetPayload := `{"label":"Test Wallet","currency":"ETH","amount":2.0}`
	req := httptest.NewRequest("POST", "/assets", bytes.NewBufferString(assetPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", w.Code)
	}
}
