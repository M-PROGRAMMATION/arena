package tests

import (
	gamePackage "arena/src/map"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var dirsNames = []string{
	"niveau1",
	"niveau2",
	"niveau3",
	"invalides",
}

var validMapDirs = []string{
	"niveau1",
	"niveau2",
	"niveau3",
}

func TestParsingMap(t *testing.T) {
	for _, dirName := range dirsNames {
		path := filepath.Join("../assets/cartes", dirName)

		files, err := os.ReadDir(path)
		if err != nil {
			t.Fatal(err)
		}

		for _, file := range files {
			if !file.IsDir() {
				testMap, err := gamePackage.ParsingMap(filepath.Join(path, file.Name()))
				if err != nil {
					t.Fatal(err)
				}
				fmt.Println(testMap)
			}
		}
	}
}

func TestMapIsValid(t *testing.T) {
	for _, dirName := range validMapDirs {
		path := filepath.Join("../assets/cartes", dirName)

		files, err := os.ReadDir(path)
		if err != nil {
			t.Fatal(err)
		}

		for _, file := range files {
			if file.IsDir() || filepath.Ext(file.Name()) != ".map" {
				continue
			}

			m, err := gamePackage.LoadMap(filepath.Join(path, file.Name()))
			if err != nil {
				t.Fatalf("carte %s/%s : %v", dirName, file.Name(), err)
			}

			if !gamePackage.MapIsValid(m.Grid, m.Width, m.Height) {
				t.Errorf("carte attendue valide mais jugée invalide : %s/%s", dirName, file.Name())
			}
		}
	}
}

func TestMapIsInvalid(t *testing.T) {
	path := filepath.Join("../assets/cartes", "invalides")

	files, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".map" {
			continue
		}

		m, err := gamePackage.LoadMap(filepath.Join(path, file.Name()))
		if err != nil {
			continue
		}

		if gamePackage.MapIsValid(m.Grid, m.Width, m.Height) {
			t.Errorf("carte attendue invalide mais jugée valide : invalides/%s", file.Name())
		}
	}
}
