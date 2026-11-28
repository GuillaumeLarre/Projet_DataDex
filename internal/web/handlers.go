package web

import (
	"html/template"
	"log"
	"net/http"
)

func AfficherListeExercices(w http.ResponseWriter, r *http.Request) {
	if ErreurChargementExercice != nil {
		renderAPIError(w, ErreurChargementExercice)
		return
	}
	t, err := template.ParseFiles("template/exercices.html")
	if err != nil {
		log.Print("Problème avec Parsefiles :", err)
		http.Error(w, "Un problème est survenu avec notre serveur", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, EnsembleExercices)
	if err != nil {
		log.Print("Problème avec Execute :", err)
		http.Error(w, "Un problème est survenu avec notre serveur", http.StatusInternalServerError)
		return
	}
}
