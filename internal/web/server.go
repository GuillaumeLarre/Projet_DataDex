package web

import (
	"Projet_DataDex/internal/api"
	"errors"
	"log"
	"net/http"
)

var EnsembleExercices []api.Exercise
var EnsembleMuscles []api.Muscle
var EnsembleCategories []api.Category
var EnsembleEquipements []api.Equipment
var ErreurChargementExercice error

func renderAPIError(w http.ResponseWriter, err error) {
	log.Println("erreur API :", err)
	switch {
	case errors.Is(err, api.ErrNotFound):
		http.Error(w, "Cette ressource n'existe pas.", 404)
	case errors.Is(err, api.ErrUnauthorized):
		http.Error(w, "L'API refuse notre accès.", 502)
	case api.IsTimeout(err):
		http.Error(w, "L'API met trop de temps.", 504)
	default:
		http.Error(w, "L'API est indisponible.", 502)
	}
}
