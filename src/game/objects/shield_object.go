package objects

import mapPackage "arena/src/map"

func IsShield(m mapPackage.Map, x int, y int) bool {
	return m.Grid[y][x] == '!'
}

func ShieldInRow(m mapPackage.Map, y int, from int, to int) bool {
	for x := from + 1; x < to; x++ {
		if IsShield(m, x, y) {
			return true
		}
	}
	return false
}

func ShieldInColumn(m mapPackage.Map, x int, from int, to int) bool {
	for y := from + 1; y < to; y++ {
		if IsShield(m, x, y) {
			return true
		}
	}
	return false
}
