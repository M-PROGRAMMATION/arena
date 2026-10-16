package objects

import mapPackage "arena/src/map"

func IsTreasure(m mapPackage.Map, x int, y int) bool {
	return m.Grid[y][x] == '$'
}

func TreasureInRow(m mapPackage.Map, y int, from int, to int) bool {
	for x := from + 1; x < to; x++ {
		if IsTreasure(m, x, y) {
			return true
		}
	}
	return false
}

func TreasureInColumn(m mapPackage.Map, x int, from int, to int) bool {
	for y := from + 1; y < to; y++ {
		if IsTreasure(m, x, y) {
			return true
		}
	}
	return false
}
