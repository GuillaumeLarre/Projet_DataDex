package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 3 * time.Second,
}

const UrlExercices = "https://wger.de/api/v2/exerciseinfo/"
const UrlMuscle = "https://wger.de/api/v2/muscle/"
const UrlEquipment = "https://wger.de/api/v2/equipment/"
const UrlExerciseCategory = "https://wger.de/api/v2/exercisecategory/"

func IsTimeout(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}
	return false
}

var (
	ErrNotFound     = errors.New("Cette ressource est introuvable")
	ErrUnauthorized = errors.New("L'API nous refuse l'accès")
)

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

func recupererDonnees[T any](chemin string) (ApiData[T], error) {
	var donneesApi ApiData[T]
	err := getJson(chemin, &donneesApi)
	return donneesApi, err
}

func ChargerToutesLesDonnees[T any](urlInitiale string) ([]T, error) {
	var donnees []T
	url := urlInitiale
	for url != "" {
		donneesApi, err := recupererDonnees[T](url)
		if err != nil {
			return nil, err
		}
		donnees = append(donnees, donneesApi.Results...)
		url = donneesApi.Next
	}
	return donnees, nil
}
