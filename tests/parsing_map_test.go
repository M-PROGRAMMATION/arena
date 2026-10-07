package tests

import (
	gamePackage "arena/src/game"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestParsingMap(t *testing.T) {
	map1 := gamePackage.ParsingMap("niveau1/arene.map")
	fmt.Println(map1)
}

func TestMapIsValid(t *testing.T) {
	dirsNames := []string{
		"niveau1",
		"niveau2",
		"niveau3",
		"invalides",
	}

	for _, dirName := range dirsNames {
		path := filepath.Join("../assets/cartes", dirName)

		files, err := os.ReadDir(path)
		if err != nil {
			t.Fatal(err)
		}

		for _, file := range files {
			if !file.IsDir() {
				testMap := gamePackage.ParsingMap(filepath.Join(dirName, file.Name()))

				if !gamePackage.MapIsValid(testMap) {
					t.Errorf("Carte invalide : %s/%s", dirName, file.Name())
				}
			}
		}
	}
}
