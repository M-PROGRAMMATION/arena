package objects

import mapPackage "arena/src/map"

func IsWall(m mapPackage.Map, x int, y int) bool {
	return m.Grid[y][x] == '#'
}

func WallInRow(m mapPackage.Map, y int, from int, to int) bool {
	for x := from + 1; x < to; x++ {
		if IsWall(m, x, y) {
			return true
		}
	}
	return false
}

func WallInColumn(m mapPackage.Map, x int, from int, to int) bool {
	for y := from + 1; y < to; y++ {
		if IsWall(m, x, y) {
			return true
		}
	}
	return false
}
