package objects

import mapPackage "arena/src/map"

func IsAmmunition(m mapPackage.Map, x int, y int) bool {
	return m.Grid[y][x] == '*'
}

func AmmunitionInRow(m mapPackage.Map, y int, from int, to int) bool {
	for x := from + 1; x < to; x++ {
		if IsAmmunition(m, x, y) {
			return true
		}
	}
	return false
}

func AmmunitionInColumn(m mapPackage.Map, x int, from int, to int) bool {
	for y := from + 1; y < to; y++ {
		if IsAmmunition(m, x, y) {
			return true
		}
	}
	return false
}
