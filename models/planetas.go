package models

type Planeta struct {
	ID          int32  `json:"id"`
	Nome        string `json:"nome" binding:"required"`
	Clima       string `json:"clima" binding:"required"`
	Terreno     string `json:"terreno" binding:"required"`
	QtdAparicao int    `json:"qtd_aparicao"`
}

type PlanetaResultAPI struct {
	Count    int64       `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		Climate        string        `json:"climate"`
		Created        string        `json:"created"`
		Diameter       string        `json:"diameter"`
		Edited         string        `json:"edited"`
		Films          []interface{} `json:"films"`
		Gravity        string        `json:"gravity"`
		Name           string        `json:"name"`
		OrbitalPeriod  string        `json:"orbital_period"`
		Population     string        `json:"population"`
		Residents      []string      `json:"residents"`
		RotationPeriod string        `json:"rotation_period"`
		SurfaceWater   string        `json:"surface_water"`
		Terrain        string        `json:"terrain"`
		URL            string        `json:"url"`
	} `json:"results"`
}
