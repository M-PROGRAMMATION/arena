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

func TestParsingMap(t *testing.T) {
	for _, dirName := range dirsNames {
		path := filepath.Join("../assets/cartes", dirName)

		files, err := os.ReadDir(path)
		if err != nil {
			t.Fatal(err)
		}

		for _, file := range files {
			if !file.IsDir() {
				testMap := gamePackage.ParsingMap("test_void.map")
				fmt.Println(testMap)
			}
		}
	}
}

func TestMapIsValid(t *testing.T) {
	for _, dirName := range dirsNames {
		path := filepath.Join("../assets/cartes", dirName)

		files, err := os.ReadDir(path)
		if err != nil {
			t.Fatal(err)
		}

		for _, file := range files {
			if !file.IsDir() {
				testMap := gamePackage.ParsingMap(filepath.Join(dirName, file.Name()))

				//TODO: Change with function write by Mathieu for get width and height.
				if !gamePackage.MapIsValid(testMap, 1, 2) {
					t.Errorf("Carte invalide : %s/%s", dirName, file.Name())
				}
			}
		}
	}
}
