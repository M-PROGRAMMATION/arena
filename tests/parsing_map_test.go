package tests

import (
	mapPackage "arena/src/map"
	"fmt"
	"testing"
)

func TestParsingMap(t *testing.T) {
	map1 := mapPackage.ParsingMap("niveau1/arene.map")
	fmt.Println(map1)
}
