package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 3 * time.Second,
}

func isTimeout(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}
	return false
}

type DonneesApi struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Exercice `json:"results"`
}

type Exercice struct {
	Id                 int                  `json:"id"`
	Nom                string               `json:"-"`
	Category           Category             `json:"category"`
	Muscles            []Muscles            `json:"muscles"`
	MusclesSecondaires []MusclesSecondaires `json:"muscles_secondary"`
	Equipement         []Equipement         `json:"equipment"`
	Translations       []Translations       `json:"translations"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Muscles struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type MusclesSecondaires struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Equipement struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Translations struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    int    `json:"language"`
}

var (
	ErrNotFound     = errors.New("Cette ressource est introuvable")
	ErrUnauthorized = errors.New("L'API nous refuse l'accès")
)

var urlExercices = "https://wger.de/api/v2/exerciseinfo/"
var ensembleExercices []Exercice
var erreurApi error

const langueAnglaise = 2

func renderAPIError(w http.ResponseWriter, err error) {
	log.Println("erreur API :", err)
	switch {
	case errors.Is(err, ErrNotFound):
		http.Error(w, "Cette ressource n'existe pas.", 404)
	case errors.Is(err, ErrUnauthorized):
		http.Error(w, "L'API refuse notre accès.", 502)
	case isTimeout(err):
		http.Error(w, "L'API met trop de temps.", 504)
	default:
		http.Error(w, "L'API est indisponible.", 502)
	}
}

func getJson(chemin string, variable any) error {
	req, err := http.NewRequest(http.MethodGet, chemin, nil)
	if err != nil {
		return fmt.Errorf("Problème lors de la préparation de la requête avec le chemin %s : %w", chemin, err)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Problème lors de la requête avec le chemin %s : %w", chemin, err)
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == http.StatusNotFound:
		return ErrNotFound
	case resp.StatusCode == http.StatusUnauthorized, resp.StatusCode == http.StatusForbidden:
		return ErrUnauthorized
	case resp.StatusCode != http.StatusOK:
		return fmt.Errorf("Un problème est survenu avec l'API sur le chemin %s : %d", chemin, resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(variable)
	if err != nil {
		return fmt.Errorf("Problème lors du décodage Json avec le chemin %s : %w", chemin, err)
	}
	return nil
}

func fetchExercices(chemin string) (DonneesApi, error) {
	var donneesApi DonneesApi
	err := getJson(chemin, &donneesApi)
	return donneesApi, err
}

func chargerTousLesExercices() ([]Exercice, error) {
	var exercices []Exercice
	url := urlExercices
	for url != "" {
		donneesApi, err := fetchExercices(url)
		if err != nil {
			return nil, err
		}
		exercices = append(exercices, donneesApi.Results...)
		url = donneesApi.Next
	}
	return exercices, nil
}

func afficherListeExercices(w http.ResponseWriter, r *http.Request) {
	if erreurApi != nil {
		renderAPIError(w, erreurApi)
		return
	}
	t, err := template.ParseFiles("template/exercices.html")
	if err != nil {
		log.Print("Problème avec Parsefiles :", err)
		http.Error(w, "Un problème est survenu avec notre serveur", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, ensembleExercices)
	if err != nil {
		log.Print("Problème avec Execute :", err)
		http.Error(w, "Un problème est survenu avec notre serveur", http.StatusInternalServerError)
		return
	}
}

func garderExercicesAvecNomEnAnglais(listeExercices []Exercice) []Exercice {
	var nouvelleListe []Exercice
	for i := 0; i < len(listeExercices); i++ {
		for j := 0; j < len(listeExercices[i].Translations); j++ {
			if listeExercices[i].Translations[j].Language == langueAnglaise {
				listeExercices[i].Nom = listeExercices[i].Translations[j].Name
				nouvelleListe = append(nouvelleListe, listeExercices[i])
				break
			}
		}
	}
	return nouvelleListe
}

func main() {
	exercices, err := chargerTousLesExercices()
	if err != nil {
		erreurApi = err
		log.Println("Problème de chargement du cache :", err)
	}
	if len(exercices) != 0 {
		ensembleExercices = garderExercicesAvecNomEnAnglais(exercices)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/exercices", afficherListeExercices)
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Problème de lancement du serveur :", err)
	}
}
