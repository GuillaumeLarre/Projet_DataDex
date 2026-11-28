package api

////////////////////////////////////////////////
//////////////// LES STRUCTURES ////////////////
////////////////////////////////////////////////

// Structure des données brutes reçues de l'API //
type ApiData[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}

type Exercise struct {
	Id               int           `json:"id"`
	Name             string        `json:"-"`
	Category         Category      `json:"category"`
	Muscles          []Muscle      `json:"muscles"`
	MusclesSecondary []Muscle      `json:"muscles_secondary"`
	Equipments       []Equipment   `json:"equipment"`
	Translations     []Translation `json:"translations"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Muscle struct {
	Id          int    `json:"id"`
	DisplayName string `json:"-"`
	Name        string `json:"name"`
	NameEn      string `json:"name_en"`
}

type Equipment struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Translation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    int    `json:"language"`
}
