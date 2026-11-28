package api

var correspondanceIdMuscleEtNom = map[int]string{
	14: "External abdominal oblique",
}

const langueAnglaise = 2

func GarderExercicesAvecNomEnAnglais(listeExercices []Exercise) []Exercise {
	var nouvelleListe []Exercise
	for i := 0; i < len(listeExercices); i++ {
		for j := 0; j < len(listeExercices[i].Translations); j++ {
			if listeExercices[i].Translations[j].Language == langueAnglaise {
				listeExercices[i].Name = listeExercices[i].Translations[j].Name
				nouvelleListe = append(nouvelleListe, listeExercices[i])
				break
			}
		}
	}
	return nouvelleListe
}

func definirNomAffichageMuscle(muscle *Muscle) {
	if muscle.NameEn == "" {
		nom, ok := correspondanceIdMuscleEtNom[muscle.Id]
		if !ok {
			muscle.DisplayName = muscle.Name
		} else {
			muscle.DisplayName = nom
		}
	} else {
		muscle.DisplayName = muscle.NameEn
	}
}

func DefinirNomsAffichageMuscle(listeMuscle []Muscle) {
	for i := 0; i < len(listeMuscle); i++ {
		definirNomAffichageMuscle(&listeMuscle[i])
	}
}
