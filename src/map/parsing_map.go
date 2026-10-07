package _map

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ParsingMap(nameCard string) []string {
	path := filepath.Join("../assets", "cartes", nameCard)
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr,
			"lecture impossible :", err)
		os.Exit(1)
	}
	contenu := string(data) // []byte → string
	lignes := strings.Split(contenu, "\n")
	return lignes
}

/*func MapIsValid(){
	// TODO:
}*/
