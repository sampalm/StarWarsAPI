package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sampalm/StarWarsAPI/config"

	"github.com/stretchr/testify/assert"
)

func init() {
	config.Conn("starwars_test")
}

func request(h http.Handler, method string, path string, data []byte) *httptest.ResponseRecorder {
	r, _ := http.NewRequest(method, path, bytes.NewBuffer(data))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

// Cadastra o planeta "Tatooine" no banco de dados. Deve retornar o codigo 201.
func TestCreatePlaneta(t *testing.T) {
	data := []byte(`{"nome": "Tatooine","clima": "temperate","terreno": "grasslands, mountains"}`)
	router := NewRouter()
	w := request(router, "POST", "/planetas/create", data)
	assert.Equal(t, http.StatusCreated, w.Code)

	data = []byte(`{"nome": "Tatooine","clima": " ","terreno": ""}`)
	w = request(router, "POST", "/planetas/create", data)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	data = []byte(`{"nome": " ","clima": "temperate","terreno": "grasslands, mountains"}`)
	w = request(router, "POST", "/planetas/create", data)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Lista todos os planetas cadastrados no banco de dados. Deve retornar o codigo 200.
func TestListPlaneta(t *testing.T) {
	router := NewRouter()
	w := request(router, "GET", "/planetas/", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Seleciona o planeta "Eriadu" do banco de dados. Deve retornar o codigo 200.
func TestQueryPlanetaByName(t *testing.T) {
	router := NewRouter()
	w := request(router, "GET", "/planetas/Tatooine", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestQueryPlanetaById Seleciona o item 1 do banco de dados. Deve retornar o codigo 200.
func TestQueryPlanetaById(t *testing.T) {
	router := NewRouter()
	w := request(router, "GET", "/planetas/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Remove o item 8 do banco de dados. Deve retornar o codigo 200.
func TestRemovePlaneta(t *testing.T) {
	router := NewRouter()
	w := request(router, "DELETE", "/planetas/remove/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}
