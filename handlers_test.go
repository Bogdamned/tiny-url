package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
)

func TestTinifyHandler(t *testing.T) {
	testUrl := "https://www.google.com/"
	type test struct {
		id       uint64
		expected string
	}

	payload := struct {
		URL string `json:"url"`
	}{testUrl}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// InitHandler
	h, err := initHandler()
	if err != nil {
		t.Errorf("Failed to initialize handler: %v", err.Error())
	}

	id, err := h.urlDB.GetID()
	if err != nil {
		t.Errorf("Failure: %v", err.Error())
	}

	resultMask := "{\"tinyURL\":\"%s\"}"
	expected := fmt.Sprintf(resultMask, Encode(uint64(id+1)))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Tinify)
	handler.ServeHTTP(rr, req)

	//Cleanup test data from DB
	h.urlDB.DeleteLast()

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	result := rr.Body.String()
	if result != expected {
		t.Errorf("tinify encoder returned wrong value: got %v want %v", result, expected)
	}

}

func TestTinyRedirectHandler(t *testing.T) {
	// InitHandler
	h, err := initHandler()
	if err != nil {
		t.Errorf("Failed to initialize handler: %v", err.Error())
	}
	// Get last ID
	id, _ := h.urlDB.GetID()
	// Get url
	url, _ := h.urlDB.Get(id)

	req, err := http.NewRequest("GET", "/"+Encode(uint64(id)), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.TinyRedirect)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTemporaryRedirect)
	}
	if rr.HeaderMap["Location"][0] != url {
		t.Errorf("handler redirected to wrong url: got %v want %v", url, rr.HeaderMap["Location"][0])
	}
}

func initHandler() (*Handler, error) {
	db, err := NewRedisDB(Config{
		Addr:     viper.GetString("db.addr"),
		Password: "",
		DB:       viper.GetInt("db.db"),
	})
	if err != nil {
		return nil, err
	}

	return NewHandler(NewLocalCache(), NewUrlRedis(db)), nil
}
