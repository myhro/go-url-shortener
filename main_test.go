package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	router := gin.New()
	setupRouter(router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Go URL Shortener", w.Body.String())
}

func TestMalformedRequest(t *testing.T) {
	setupDB(":memory:")
	router := gin.New()
	setupRouter(router)

	w := httptest.NewRecorder()
	body := []byte("not-a-json")
	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNewURL(t *testing.T) {
	setupDB(":memory:")
	router := gin.New()
	setupRouter(router)

	google := "https://google.com/"
	u := URL{Full: google}
	body, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(w, req)

	res := URL{}
	err = json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, google, res.Full)
	assert.Equal(t, "1", res.Hash)
}

func TestURLDetails(t *testing.T) {
	setupDB(":memory:")
	router := gin.New()
	setupRouter(router)

	google := "https://google.com/"
	u := URL{Full: google}
	body, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/1/details", nil)
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(w, req)

	res := URL{}
	err = json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, google, res.Full)
	assert.Equal(t, "1", res.Hash)
}

func TestURLDetailsNotFound(t *testing.T) {
	setupDB(":memory:")
	router := gin.New()
	setupRouter(router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/1/details", nil)
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Not Found", w.Body.String())
}

func TestURLNotFound(t *testing.T) {
	setupDB(":memory:")
	router := gin.New()
	setupRouter(router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/1", nil)
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Not Found", w.Body.String())
}

func TestURLRedirect(t *testing.T) {
	setupDB(":memory:")
	router := gin.New()
	setupRouter(router)

	google := "https://google.com/"
	u := URL{Full: google}
	body, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/1", nil)
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(w, req)

	location := ""
	if len(w.HeaderMap["Location"]) > 0 {
		location = w.HeaderMap["Location"][0]
	}

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
	assert.Equal(t, google, location)
}
