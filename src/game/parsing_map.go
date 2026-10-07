package game

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var validCharacters = []string{
	"#", // = Mur Infranchissable (Arrête les tirs)
	"~", // Eau Immobilisé (Immobilisé au Tour Suivant)
	"^", // Piège (−15 PV)
	"+", // Soin
	"*", // Munitions
	"!", // Bouclier
	"$", // Trésor (Ramassés en marchant)
	"1", // Joueur 1
	"2", // Joueur 2
	".", // Endroit ou le joueur peut bouger
}

func ParsingMap(nameCard string) []string {
	path := filepath.Join("../assets", "cartes", nameCard)
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Carte introuvable :", err)
		os.Exit(1)
	}
	contenu := string(data) // []byte → string
	lignes := strings.Split(contenu, "\n")
	return lignes[4:]
}

/*
// Une carte valide respecte toutes ces règles. Votre parser doit produire **exactement** ces messages, comme `./arbitre --verifier` :
//
| Fichier vide | `<fichier> : fichier vide` |
| Clé inconnue | `<fichier> ligne N : clé inconnue "X"` |
| Clé répétée | `<fichier> ligne N : X défini deux fois` |
| Clé absente | `<fichier> : en-tête incomplet : X manquant` |
| Valeur incorrecte | `<fichier> ligne N : TAILLE attend deux entiers` (et messages équivalents pour les autres clés) |
| Ligne de mauvaise largeur | `<fichier> ligne N : 12 cases, 11 attendues` |
| Mauvais nombre de lignes | `<fichier> : 6 lignes de grille, 7 attendues` |
| Symbole inconnu | `<fichier> ligne N colonne C : symbole inconnu "?"` |  ✅
| Bordure ouverte | `<fichier> ligne N colonne C : la bordure doit être un mur` |
| Départ répété | `<fichier> ligne N colonne C : départ du joueur 1 déjà défini` | ✅
| Départ absent | `<fichier> : aucun départ pour le joueur 2` | ✅
*/

func MapIsValid(maps []string) bool {
	numberOneIsUsed := false
	numberTwoIsUsed := false

	for _, item := range maps {
		for _, char := range item {
			charStr := string(char)

			if !slices.Contains(validCharacters, charStr) {
				println("La game n'est pas valide, symbole inconnu!")
				return false
			}

			if charStr == "1" {
				if numberOneIsUsed {
					println("Impossible! Départ double!")
					return false
				}
				numberOneIsUsed = true
			}

			if charStr == "2" {
				if numberTwoIsUsed {
					println("Impossible! Départ double!")
					return false
				}
				numberTwoIsUsed = true
			}
		}
	}

	if !numberOneIsUsed {
		println("La game n'est pas valide, car le joueur 1 n'est pas présent!")
		return false
	} else if !numberTwoIsUsed {
		println("La game n'est pas valide, car le joueur 2 n'est pas présent!")
	}

	return true
}
