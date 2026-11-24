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

var baseUrl = "https://wger.de/api/v2/"

type Exercices struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Results `json:"results"`
}

type Results struct {
	Id                 int                  `json:"id"`
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
	req, err := http.NewRequest(http.MethodGet, baseUrl+chemin, nil)
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

func fetchExercices() (Exercices, error) {
	var listeExercices Exercices
	err := getJson("exerciseinfo", &listeExercices)
	return listeExercices, err
}

func afficherListeExercices(w http.ResponseWriter, r *http.Request) {
	listeExercices, err := fetchExercices()
	if err != nil {
		renderAPIError(w, err)
		return
	}
	t, err := template.ParseFiles("template/exercices.html")
	if err != nil {
		log.Print("Problème avec Parsefiles :", err)
		http.Error(w, "Un problème est survenu avec notre serveur", http.StatusInternalServerError)
		return
	}
	t.Execute(w, listeExercices)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/exercices", afficherListeExercices)
	http.ListenAndServe(":8080", mux)
}
