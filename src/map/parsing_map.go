package _map

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Map struct {
	Name       string
	Difficulty int
	Rounds     int
	Width      int
	Height     int
	Grid       []string
}

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

func LoadMap(nameCard string) (Map, error) {
	lignes, err := ParsingMap(nameCard)
	if err != nil {
		return Map{}, err
	}
	var m Map
	i := 0

	for ; i < len(lignes); i++ {
		ligne := strings.TrimRight(lignes[i], "\r")
		if strings.TrimSpace(ligne) == "" {
			continue
		}
		if strings.HasPrefix(ligne, "#") {
			break
		}

		cle, valeur, ok := strings.Cut(ligne, " ")
		if !ok {
			return m, fmt.Errorf("ligne %d invalide : %q", i+1, ligne)
		}
		valeur = strings.TrimSpace(valeur)

		var err error
		switch cle {
		case "NOM":
			m.Name = valeur
		case "DIFFICULTE":
			m.Difficulty, err = strconv.Atoi(valeur)
		case "TOURS":
			m.Rounds, err = strconv.Atoi(valeur)
		case "TAILLE":
			_, err = fmt.Sscanf(valeur, "%d %d", &m.Width, &m.Height)
		default:
			err = fmt.Errorf("clé inconnue %q", cle)
		}
		if err != nil {
			return m, fmt.Errorf("ligne %d (%s) : %w", i+1, cle, err)
		}
	}

	for ; i < len(lignes); i++ {
		ligne := strings.TrimRight(lignes[i], "\r")
		if ligne == "" {
			continue
		}
		m.Grid = append(m.Grid, ligne)
	}

	if len(m.Grid) != m.Height {
		return m, fmt.Errorf("hauteur attendue %d, trouvée %d", m.Height, len(m.Grid))
	}
	for y, row := range m.Grid {
		if len([]rune(row)) != m.Width {
			return m, fmt.Errorf("ligne %d de la grille : largeur %d au lieu de %d",
				y+1, len([]rune(row)), m.Width)
		}
	}

	return m, nil
}

func ParsingMap(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("carte introuvable : %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, fmt.Errorf("carte vide : %s", path)
	}
	return strings.Split(string(data), "\n"), nil
}

func MapIsValid(grid []string, width, height int) bool {
	numberOneIsUsed := false
	numberTwoIsUsed := false

	if width < 3 || height < 3 {
		println("TAILLE invalide : minimum 3 !")
		return false
	}

	if len(grid) != height {
		fmt.Printf("La map fait %d lignes au lieu de %d !\n", len(grid), height)
		return false
	}

	for y, line := range grid {
		if n := utf8.RuneCountInString(line); n != width {
			fmt.Printf("La ligne %d fait %d caractères au lieu de %d !\n", y+1, n, width)
			return false
		}
	}

	for _, item := range grid {
		for _, char := range item {
			charStr := string(char)

			if !slices.Contains(validCharacters, charStr) {
				println("La map n'est pas valide, symbole inconnu!")
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
		println("La map n'est pas valide, car le joueur 1 n'est pas présent!")
		return false
	} else if !numberTwoIsUsed {
		println("La map n'est pas valide, car le joueur 2 n'est pas présent!")
	}

	return true
}
