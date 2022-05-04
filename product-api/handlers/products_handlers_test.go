package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barmansurajit/go-microservices/product-api/data"
	assert "github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler := http.Handler(new(Products))
	handler.ServeHTTP(rr, req)

	data := rr.Body.String()

	assert.Equal(rr.Code, http.StatusOK)
	assert.NotEmpty(data)
	assert.Contains(data, "Latte")
}

func TestAddProduct(t *testing.T) {
	assert := assert.New(t)

	p := &data.Product{Name: "Mocachino", Price: 1.25, Description: "Another type of coffee", SKU: "N078"}
	pjson, _ := json.Marshal(p)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(pjson))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.Handler(new(Products))
	handler.ServeHTTP(rr, req)

	assert.Equal(rr.Code, http.StatusCreated)
}

func TestUpdateProduct(t *testing.T) {
	assert := assert.New(t)

	p := &data.Product{Name: "Mocachino", Price: 1.25, Description: "Another type of coffee", SKU: "N078"}
	pjson, _ := json.Marshal(p)
	req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(pjson))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.Handler(new(Products))
	handler.ServeHTTP(rr, req)

	assert.Equal(rr.Code, http.StatusOK)
}
