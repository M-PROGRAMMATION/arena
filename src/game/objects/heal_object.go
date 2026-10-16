package objects

import mapPackage "arena/src/map"

func IsHeal(m mapPackage.Map, x int, y int) bool {
	return m.Grid[y][x] == '+'
}

func HealInRow(m mapPackage.Map, y int, from int, to int) bool {
	for x := from + 1; x < to; x++ {
		if IsHeal(m, x, y) {
			return true
		}
	}
	return false
}

func HealInColumn(m mapPackage.Map, x int, from int, to int) bool {
	for y := from + 1; y < to; y++ {
		if IsHeal(m, x, y) {
			return true
		}
	}
	return false
}
