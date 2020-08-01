package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/sampalm/StarWarsAPI/models"
)

func SearchPlanetaAPI(s string) (count int, err error) {
	var data models.PlanetaResultAPI
	res, err := http.Get("https://swapi.dev/api/planets/?search=" + url.QueryEscape(s))
	if err != nil {
		return count, err
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return count, err
	}

	if data.Count == 0 {
		return count, fmt.Errorf("SearchPlanetaAPI: Planeta nao encontrado")
	}
	if data.Count > 1 {
		return count, fmt.Errorf("SearchPlanetaAPI: Sua pesquisa precisa ser mais especifica")
	}

	return len(data.Results[0].Films), nil
}
