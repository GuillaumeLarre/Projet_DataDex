package main

import (
	"Projet_DataDex/internal/api"
	"Projet_DataDex/internal/web"
	"log"
	"net/http"
)

func main() {
	exercices, err := api.ChargerToutesLesDonnees[api.Exercise](api.UrlExercices)
	if err != nil {
		web.ErreurChargementExercice = err
		log.Println("Problème de chargement des exercices :", err)
	}
	categories, err := api.ChargerToutesLesDonnees[api.Category](api.UrlExerciseCategory)
	if err != nil {
		log.Println("Problème de chargement des catégories :", err)
	}
	web.EnsembleCategories = categories
	log.Println("Nombre de catégories chargés :", len(web.EnsembleCategories))
	muscles, err := api.ChargerToutesLesDonnees[api.Muscle](api.UrlMuscle)
	if err != nil {
		log.Println("Problème de chargement des muscles :", err)
	}
	web.EnsembleMuscles = muscles
	api.DefinirNomsAffichageMuscle(web.EnsembleMuscles)
	log.Println("Nombre de muscles chargés :", len(web.EnsembleMuscles))
	equipements, err := api.ChargerToutesLesDonnees[api.Equipment](api.UrlEquipment)
	if err != nil {
		log.Println("Problème de chargement des équipements :", err)
	}
	web.EnsembleEquipements = equipements
	log.Println("Nombre d'équipements chargés :", len(web.EnsembleEquipements))
	if len(exercices) != 0 {
		web.EnsembleExercices = api.GarderExercicesAvecNomEnAnglais(exercices)
		log.Println("Nombre d'exercices chargés :", len(web.EnsembleExercices))
		for i := 0; i < len(web.EnsembleExercices); i++ {
			api.DefinirNomsAffichageMuscle(web.EnsembleExercices[i].Muscles)
			api.DefinirNomsAffichageMuscle(web.EnsembleExercices[i].MusclesSecondary)
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/exercices", web.AfficherListeExercices)
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Problème de lancement du serveur :", err)
	}
}
