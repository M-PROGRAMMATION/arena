package _map

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
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
	"#", // Mur : infranchissable, arrête les tirs
	".", // Sol
	"~", // Eau : immobilise au tour suivant
	"^", // Piège : -15 PV
	"+", // Soin (v2)
	"*", // Munitions (v2)
	"!", // Bouclier (v2)
	"$", // Trésor (v2)
	"1", // Départ joueur 1
	"2", // Départ joueur 2
}

var headerKeys = []string{"NOM", "DIFFICULTE", "TOURS", "TAILLE"}

func LoadMap(path string) (Map, error) {
	file := filepath.Base(path)

	data, err := os.ReadFile(path)
	if err != nil {
		return Map{}, err
	}

	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if allBlank(lines) {
		return Map{}, fileErr(file, "fichier vide")
	}

	var m Map
	seen := map[string]int{}
	i := 0

	for ; i < len(lines); i++ {
		lineNo := i + 1
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			break
		}

		key, value, _ := strings.Cut(trimmed, " ")
		value = strings.TrimSpace(value)

		if !slices.Contains(headerKeys, key) {
			return Map{}, lineErr(file, lineNo, fmt.Sprintf("clé inconnue %q", key))
		}
		if _, dup := seen[key]; dup {
			return Map{}, lineErr(file, lineNo, fmt.Sprintf("%s défini deux fois", key))
		}

		switch key {
		case "NOM":
			if value == "" {
				return Map{}, lineErr(file, lineNo, "NOM attend une valeur")
			}
			m.Name = value
		case "DIFFICULTE":
			d, err := strconv.Atoi(value)
			if err != nil {
				return Map{}, lineErr(file, lineNo, "DIFFICULTE attend un entier")
			}
			if d < 1 || d > 3 {
				return Map{}, lineErr(file, lineNo, "DIFFICULTE attend un entier entre 1 et 3")
			}
			m.Difficulty = d
		case "TOURS":
			r, err := strconv.Atoi(value)
			if err != nil {
				return Map{}, lineErr(file, lineNo, "TOURS attend un entier")
			}
			if r <= 0 {
				return Map{}, lineErr(file, lineNo, "TOURS attend un entier positif")
			}
			m.Rounds = r
		case "TAILLE":
			fields := strings.Fields(value)
			if len(fields) != 2 {
				return Map{}, lineErr(file, lineNo, "TAILLE attend deux entiers")
			}
			w, errW := strconv.Atoi(fields[0])
			h, errH := strconv.Atoi(fields[1])
			if errW != nil || errH != nil {
				return Map{}, lineErr(file, lineNo, "TAILLE attend deux entiers")
			}
			if w < 3 || h < 3 {
				return Map{}, lineErr(file, lineNo, "TAILLE attend deux entiers d'au moins 3")
			}
			m.Width, m.Height = w, h
		}
		seen[key] = lineNo
	}

	for _, key := range headerKeys {
		if _, ok := seen[key]; !ok {
			return Map{}, fileErr(file, fmt.Sprintf("en-tête incomplet : %s manquant", key))
		}
	}

	grid := lines[i:]
	if len(grid) != m.Height {
		return Map{}, fileErr(file, fmt.Sprintf("%d lignes de grille, %d attendues", len(grid), m.Height))
	}

	p1Found, p2Found := false, false

	for r, row := range grid {
		lineNo := i + r + 1
		runes := []rune(row)
		if len(runes) != m.Width {
			return Map{}, lineErr(file, lineNo, fmt.Sprintf("%d cases, %d attendues", len(runes), m.Width))
		}

		for c, ch := range runes {
			col := c + 1
			sym := string(ch)

			if !slices.Contains(validCharacters, sym) {
				return Map{}, colErr(file, lineNo, col, fmt.Sprintf("symbole inconnu %q", sym))
			}

			isBorder := r == 0 || r == m.Height-1 || c == 0 || c == m.Width-1
			if isBorder && ch != '#' {
				return Map{}, colErr(file, lineNo, col, "la bordure doit être un mur")
			}

			switch ch {
			case '1':
				if p1Found {
					return Map{}, colErr(file, lineNo, col, "départ du joueur 1 déjà défini")
				}
				p1Found = true
			case '2':
				if p2Found {
					return Map{}, colErr(file, lineNo, col, "départ du joueur 2 déjà défini")
				}
				p2Found = true
			}
		}

		m.Grid = append(m.Grid, string(runes))
	}

	if !p1Found {
		return Map{}, fileErr(file, "aucun départ pour le joueur 1")
	}
	if !p2Found {
		return Map{}, fileErr(file, "aucun départ pour le joueur 2")
	}

	return m, nil
}

func allBlank(lines []string) bool {
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			return false
		}
	}
	return true
}

func fileErr(file, msg string) error {
	return fmt.Errorf("%s : %s", file, msg)
}

func lineErr(file string, line int, msg string) error {
	return fmt.Errorf("%s ligne %d : %s", file, line, msg)
}

func colErr(file string, line, col int, msg string) error {
	return fmt.Errorf("%s ligne %d colonne %d : %s", file, line, col, msg)
}
