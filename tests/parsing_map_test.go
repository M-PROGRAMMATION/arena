package tests

import (
	gamePackage "arena/src/map"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var validMapDirs = []string{
	"niveau1",
	"niveau2",
	"niveau3",
}

func TestLoadMap_ValidMaps(t *testing.T) {
	for _, dirName := range validMapDirs {
		dirPath := filepath.Join("../assets/cartes", dirName)

		files, err := os.ReadDir(dirPath)
		if err != nil {
			t.Fatal(err)
		}

		for _, file := range files {
			if file.IsDir() || filepath.Ext(file.Name()) != ".map" {
				continue
			}

			t.Run(dirName+"/"+file.Name(), func(t *testing.T) {
				m, err := gamePackage.LoadMap(filepath.Join(dirPath, file.Name()))
				if err != nil {
					t.Fatalf("carte attendue valide mais rejetée : %v", err)
				}

				if m.Name == "" {
					t.Error("NOM vide")
				}
				if m.Difficulty < 1 || m.Difficulty > 3 {
					t.Errorf("DIFFICULTE hors bornes : %d", m.Difficulty)
				}
				if m.Rounds <= 0 {
					t.Errorf("TOURS invalide : %d", m.Rounds)
				}
				if len(m.Grid) != m.Height {
					t.Errorf("Grid a %d lignes, attendu %d", len(m.Grid), m.Height)
				}
				for y, row := range m.Grid {
					if n := len([]rune(row)); n != m.Width {
						t.Errorf("ligne %d : largeur %d, attendu %d", y, n, m.Width)
					}
				}
			})
		}
	}
}

func TestLoadMap_InvalidFixtures(t *testing.T) {
	dirPath := filepath.Join("../assets/cartes/invalides")

	data, err := os.ReadFile(filepath.Join(dirPath, "attendus.txt"))
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range strings.Split(strings.TrimRight(string(data), "\n"), "\n") {
		want := want
		fileName, _, ok := strings.Cut(want, " ")
		if !ok {
			t.Fatalf("ligne attendus.txt invalide : %q", want)
		}

		t.Run(fileName, func(t *testing.T) {
			_, err := gamePackage.LoadMap(filepath.Join(dirPath, fileName))
			if err == nil {
				t.Fatalf("carte attendue invalide, mais acceptée")
			}
			if err.Error() != want {
				t.Errorf("message = %q\nattendu  = %q", err.Error(), want)
			}
		})
	}
}

func TestLoadMap_FileNotFound(t *testing.T) {
	_, err := gamePackage.LoadMap(filepath.Join(t.TempDir(), "absent.map"))
	if err == nil {
		t.Fatal("attendu une erreur pour un fichier inexistant")
	}
}

func baseMapLines() []string {
	return []string{
		"NOM Arene",
		"DIFFICULTE 1",
		"TOURS 150",
		"TAILLE 11 7",
		"###########",
		"#1........#",
		"#...#.#...#",
		"#*...#...*#",
		"#...#.#...#",
		"#........2#",
		"###########",
	}
}

func writeMap(t *testing.T, lines []string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "carte.map")
	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadMap_HeaderCleInconnue(t *testing.T) {
	lines := baseMapLines()
	lines[0] = "NAME Arene"

	_, err := gamePackage.LoadMap(writeMap(t, lines))
	assertErr(t, err, `carte.map ligne 1 : clé inconnue "NAME"`)
}

func TestLoadMap_HeaderCleRepetee(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func([]string) []string
		wantMsg string
	}{
		{
			"NOM",
			func(l []string) []string { return append([]string{l[0], "NOM Autre"}, l[1:]...) },
			`carte.map ligne 2 : NOM défini deux fois`,
		},
		{
			"DIFFICULTE",
			func(l []string) []string {
				return append(append([]string{}, l[:2]...), append([]string{"DIFFICULTE 2"}, l[2:]...)...)
			},
			`carte.map ligne 3 : DIFFICULTE défini deux fois`,
		},
		{
			"TOURS",
			func(l []string) []string {
				return append(append([]string{}, l[:3]...), append([]string{"TOURS 200"}, l[3:]...)...)
			},
			`carte.map ligne 4 : TOURS défini deux fois`,
		},
		{
			"TAILLE",
			func(l []string) []string {
				return append(append([]string{}, l[:4]...), append([]string{"TAILLE 9 9"}, l[4:]...)...)
			},
			`carte.map ligne 5 : TAILLE défini deux fois`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := gamePackage.LoadMap(writeMap(t, c.mutate(baseMapLines())))
			assertErr(t, err, c.wantMsg)
		})
	}
}

func TestLoadMap_HeaderCleAbsente(t *testing.T) {
	cases := []struct {
		name    string
		remove  int
		wantMsg string
	}{
		{"NOM", 0, `carte.map : en-tête incomplet : NOM manquant`},
		{"DIFFICULTE", 1, `carte.map : en-tête incomplet : DIFFICULTE manquant`},
		{"TOURS", 2, `carte.map : en-tête incomplet : TOURS manquant`},
		{"TAILLE", 3, `carte.map : en-tête incomplet : TAILLE manquant`},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lines := baseMapLines()
			lines = append(lines[:c.remove], lines[c.remove+1:]...)
			_, err := gamePackage.LoadMap(writeMap(t, lines))
			assertErr(t, err, c.wantMsg)
		})
	}
}

func TestLoadMap_HeaderValeurIncorrecte(t *testing.T) {
	cases := []struct {
		name    string
		line    int
		value   string
		wantMsg string
	}{
		{"DIFFICULTE non numérique", 1, "DIFFICULTE abc", `carte.map ligne 2 : DIFFICULTE attend un entier`},
		{"DIFFICULTE trop basse", 1, "DIFFICULTE 0", `carte.map ligne 2 : DIFFICULTE attend un entier entre 1 et 3`},
		{"DIFFICULTE trop haute", 1, "DIFFICULTE 4", `carte.map ligne 2 : DIFFICULTE attend un entier entre 1 et 3`},
		{"DIFFICULTE négative", 1, "DIFFICULTE -1", `carte.map ligne 2 : DIFFICULTE attend un entier entre 1 et 3`},
		{"TOURS non numérique", 2, "TOURS abc", `carte.map ligne 3 : TOURS attend un entier`},
		{"TOURS nul", 2, "TOURS 0", `carte.map ligne 3 : TOURS attend un entier positif`},
		{"TOURS négatif", 2, "TOURS -10", `carte.map ligne 3 : TOURS attend un entier positif`},
		{"TAILLE non numérique", 3, "TAILLE abc 7", `carte.map ligne 4 : TAILLE attend deux entiers`},
		{"TAILLE un seul entier", 3, "TAILLE 11", `carte.map ligne 4 : TAILLE attend deux entiers`},
		{"TAILLE trois entiers", 3, "TAILLE 11 7 3", `carte.map ligne 4 : TAILLE attend deux entiers`},
		{"TAILLE largeur < 3", 3, "TAILLE 2 7", `carte.map ligne 4 : TAILLE attend deux entiers d'au moins 3`},
		{"TAILLE hauteur < 3", 3, "TAILLE 11 2", `carte.map ligne 4 : TAILLE attend deux entiers d'au moins 3`},
		{"NOM vide", 0, "NOM", `carte.map ligne 1 : NOM attend une valeur`},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lines := baseMapLines()
			lines[c.line] = c.value
			_, err := gamePackage.LoadMap(writeMap(t, lines))
			assertErr(t, err, c.wantMsg)
		})
	}
}

func TestLoadMap_LigneMauvaiseLargeur(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func([]string) []string
		wantMsg string
	}{
		{
			"trop étroite",
			func(l []string) []string { l[6] = l[6][:len(l[6])-1]; return l },
			`carte.map ligne 7 : 10 cases, 11 attendues`,
		},
		{
			"trop large",
			func(l []string) []string { l[6] = l[6] + "#"; return l },
			`carte.map ligne 7 : 12 cases, 11 attendues`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := gamePackage.LoadMap(writeMap(t, c.mutate(baseMapLines())))
			assertErr(t, err, c.wantMsg)
		})
	}
}

func TestLoadMap_MauvaisNombreDeLignes(t *testing.T) {
	t.Run("trop peu", func(t *testing.T) {
		lines := baseMapLines()
		lines = lines[:len(lines)-1]
		_, err := gamePackage.LoadMap(writeMap(t, lines))
		assertErr(t, err, `carte.map : 6 lignes de grille, 7 attendues`)
	})

	t.Run("trop", func(t *testing.T) {
		lines := append(baseMapLines(), "###########")
		_, err := gamePackage.LoadMap(writeMap(t, lines))
		assertErr(t, err, `carte.map : 8 lignes de grille, 7 attendues`)
	})

	t.Run("ligne vide surnuméraire en fin de fichier", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "carte.map")
		content := strings.Join(baseMapLines(), "\n") + "\n\n"
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		_, err := gamePackage.LoadMap(path)
		assertErr(t, err, `carte.map : 8 lignes de grille, 7 attendues`)
	})
}

func TestLoadMap_SymboleInconnu(t *testing.T) {
	lines := baseMapLines()
	lines[7] = "#*.?.#...*#"
	_, err := gamePackage.LoadMap(writeMap(t, lines))
	assertErr(t, err, `carte.map ligne 8 colonne 4 : symbole inconnu "?"`)
}

func TestLoadMap_BordureOuverte(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func([]string) []string
		wantMsg string
	}{
		{
			"première ligne",
			func(l []string) []string { l[4] = l[4][:5] + "." + l[4][6:]; return l },
			`carte.map ligne 5 colonne 6 : la bordure doit être un mur`,
		},
		{
			"dernière ligne",
			func(l []string) []string { l[10] = l[10][:10] + "."; return l },
			`carte.map ligne 11 colonne 11 : la bordure doit être un mur`,
		},
		{
			"colonne de gauche",
			func(l []string) []string { l[6] = "." + l[6][1:]; return l },
			`carte.map ligne 7 colonne 1 : la bordure doit être un mur`,
		},
		{
			"colonne de droite",
			func(l []string) []string { l[6] = l[6][:10] + "."; return l },
			`carte.map ligne 7 colonne 11 : la bordure doit être un mur`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := gamePackage.LoadMap(writeMap(t, c.mutate(baseMapLines())))
			assertErr(t, err, c.wantMsg)
		})
	}
}

func TestLoadMap_DepartRepete(t *testing.T) {
	t.Run("joueur 1", func(t *testing.T) {
		lines := baseMapLines()
		lines[7] = "#*1..#...*#"
		_, err := gamePackage.LoadMap(writeMap(t, lines))
		assertErr(t, err, `carte.map ligne 8 colonne 3 : départ du joueur 1 déjà défini`)
	})

	t.Run("joueur 2", func(t *testing.T) {
		lines := baseMapLines()
		lines[6] = "#.2.#.#...#"
		_, err := gamePackage.LoadMap(writeMap(t, lines))
		assertErr(t, err, `carte.map ligne 10 colonne 10 : départ du joueur 2 déjà défini`)
	})
}

func TestLoadMap_DepartAbsent(t *testing.T) {
	t.Run("joueur 1", func(t *testing.T) {
		lines := baseMapLines()
		lines[5] = "#.........#"
		_, err := gamePackage.LoadMap(writeMap(t, lines))
		assertErr(t, err, `carte.map : aucun départ pour le joueur 1`)
	})

	t.Run("joueur 2", func(t *testing.T) {
		lines := baseMapLines()
		lines[9] = "#.........#"
		_, err := gamePackage.LoadMap(writeMap(t, lines))
		assertErr(t, err, `carte.map : aucun départ pour le joueur 2`)
	})
}

func TestLoadMap_FichierVide(t *testing.T) {
	cases := []string{"", "\n", "   \n\t\n  \n"}
	for _, content := range cases {
		path := filepath.Join(t.TempDir(), "carte.map")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		_, err := gamePackage.LoadMap(path)
		assertErr(t, err, `carte.map : fichier vide`)
	}
}

func TestLoadMap_PremiereErreurGagne(t *testing.T) {
	t.Run("largeur avant symbole plus loin dans le fichier", func(t *testing.T) {
		lines := baseMapLines()
		lines[6] = lines[6][:10]
		lines[8] = "#.?.#.#...#"
		_, err := gamePackage.LoadMap(writeMap(t, lines))
		assertErr(t, err, `carte.map ligne 7 : 10 cases, 11 attendues`)
	})

	t.Run("symbole inconnu avant bordure sur la même case", func(t *testing.T) {
		lines := baseMapLines()
		lines[4] = lines[4][:5] + "?" + lines[4][6:]
		_, err := gamePackage.LoadMap(writeMap(t, lines))
		assertErr(t, err, `carte.map ligne 5 colonne 6 : symbole inconnu "?"`)
	})

	t.Run("le repère de la grille est le symbole # en tout début de ligne", func(t *testing.T) {
		lines := baseMapLines()
		lines[4] = "." + lines[4][1:]
		_, err := gamePackage.LoadMap(writeMap(t, lines))
		assertErr(t, err, `carte.map ligne 5 : clé inconnue ".##########"`)
	})
}

func TestLoadMap_ToleranceCRLF(t *testing.T) {
	path := filepath.Join(t.TempDir(), "carte.map")
	content := strings.Join(baseMapLines(), "\r\n") + "\r\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	m, err := gamePackage.LoadMap(path)
	if err != nil {
		t.Fatalf("carte CRLF attendue valide : %v", err)
	}
	if m.Width != 11 || m.Height != 7 {
		t.Errorf("dimensions incorrectes : %dx%d", m.Width, m.Height)
	}
}

func TestLoadMap_LignesVidesIgnoreesDansEntete(t *testing.T) {
	lines := baseMapLines()
	lines = append([]string{lines[0], ""}, lines[1:]...)
	m, err := gamePackage.LoadMap(writeMap(t, lines))
	if err != nil {
		t.Fatalf("ligne vide dans l'en-tête attendue tolérée : %v", err)
	}
	if m.Rounds != 150 {
		t.Errorf("TOURS = %d, attendu 150", m.Rounds)
	}
}

func TestLoadMap_OrdreDesClesLibre(t *testing.T) {
	lines := []string{
		"TAILLE 11 7",
		"TOURS 150",
		"DIFFICULTE 1",
		"NOM Arene",
	}
	lines = append(lines, baseMapLines()[4:]...)

	m, err := gamePackage.LoadMap(writeMap(t, lines))
	if err != nil {
		t.Fatalf("ordre libre des clés attendu toléré : %v", err)
	}
	if m.Name != "Arene" || m.Difficulty != 1 || m.Rounds != 150 || m.Width != 11 || m.Height != 7 {
		t.Errorf("champs incorrects : %+v", m)
	}
}

func TestLoadMap_TailleMinimale(t *testing.T) {
	lines := []string{
		"NOM Mini",
		"DIFFICULTE 1",
		"TOURS 10",
		"TAILLE 4 3",
		"####",
		"#12#",
		"####",
	}

	m, err := gamePackage.LoadMap(writeMap(t, lines))
	if err != nil {
		t.Fatalf("carte minimale attendue valide : %v", err)
	}
	if m.Width != 4 || m.Height != 3 || len(m.Grid) != 3 {
		t.Errorf("carte minimale mal chargée : %+v", m)
	}
}

func TestLoadMap_TousLesSymbolesValides(t *testing.T) {
	lines := []string{
		"NOM Kitchen",
		"DIFFICULTE 1",
		"TOURS 10",
		"TAILLE 11 5",
		"###########",
		"#1.~^+!$..#",
		"#.........#",
		"#........2#",
		"###########",
	}

	m, err := gamePackage.LoadMap(writeMap(t, lines))
	if err != nil {
		t.Fatalf("tous les symboles valides devraient être acceptés : %v", err)
	}
	if len(m.Grid) != 5 {
		t.Errorf("Grid a %d lignes, attendu 5", len(m.Grid))
	}
}

func assertErr(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatalf("attendu l'erreur %q, aucune erreur renvoyée", want)
	}
	if err.Error() != want {
		t.Errorf("message = %q\nattendu  = %q", err.Error(), want)
	}
}
