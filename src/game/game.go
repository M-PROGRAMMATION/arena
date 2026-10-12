package game

import (
	"arena/src/game/objects"
	"arena/src/logger"
	mapPackage "arena/src/map"
	parsingrobot "arena/src/parsing_robot"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Point struct {
	X int
	Y int
}

func StartGame() {
	mapPath := flag.String("carte", "", "fichier .map de la partie")
	player := flag.Int("joueur", 1, "numéro du joueur (1 ou 2)")
	flag.Parse()

	if *mapPath == "" {
		fmt.Fprintln(os.Stderr, "Erreur : option -carte obligatoire")
		os.Exit(1)
	}

	if *player != 1 && *player != 2 {
		fmt.Fprintln(os.Stderr, "Erreur : -joueur doit valoir 1 ou 2")
		os.Exit(1)
	}

	m, err := mapPackage.LoadMap(*mapPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur de carte :", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	var state []string
	var receivedAt time.Time
	turn := 0

	for scanner.Scan() {
		if receivedAt.IsZero() {
			receivedAt = time.Now() // première ligne du tour = début du chrono
		}

		line := scanner.Text()
		if line != "FIN" {
			state = append(state, line)
			continue
		}

		turn++
		fmt.Println(decide(m, state, turn, *player, receivedAt))
		state = state[:0]
		receivedAt = time.Time{}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Erreur de lecture :", err)
	}

	if err := logger.WriteFile(*player); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func decide(m mapPackage.Map, state []string, turn, playerNum int, receivedAt time.Time) (action string) {
	var enemy parsingrobot.Robot
	var player parsingrobot.Robot
	var seenMe, seenEnemy bool

	defer func() {
		logger.LogTurn(turn, playerNum, action,
			logger.Position{X: player.X, Y: player.Y, Visible: seenMe},
			logger.Position{X: enemy.X, Y: enemy.Y, Visible: seenEnemy},
			receivedAt, time.Now(),
		)
	}()

	for _, line := range state {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		if fields[0] == "MOI" {
			robot, err := parsingrobot.ParseRobot(fields)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return "ATTENDS"
			}
			player = robot
			seenMe = true
		}

		if fields[0] == "ENNEMI" {
			robot, err := parsingrobot.ParseRobot(fields)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return "ATTENDS"
			}
			enemy = robot
			seenEnemy = true
		}
	}

	if player.Y == enemy.Y && player.X < enemy.X && enemy.X-player.X <= 5 && !objects.WallInRow(m, player.Y, player.X, enemy.X) {
		return "TIRE E"
	}
	if player.Y == enemy.Y && player.X > enemy.X && player.X-enemy.X <= 5 && !wallInRow(m, player.Y, enemy.X, player.X) {
		return "TIRE O"
	}
	if player.X == enemy.X && player.Y > enemy.Y && player.Y-enemy.Y <= 5 && !objects.WallInColumn(m, player.X, enemy.Y, player.Y) {
		return "TIRE N"
	}
	if player.X == enemy.X && player.Y < enemy.Y && enemy.Y-player.Y <= 5 && !objects.WallInColumn(m, player.X, player.Y, enemy.Y) {
		return "TIRE S"
	}

	dist := distancesFrom(m, Point{X: enemy.X, Y: enemy.Y})
	here := dist[player.Y][player.X]

	if dist[player.Y][player.X+1] == here-1 {
		return "AVANCE E"
	}
	if dist[player.Y][player.X-1] == here-1 {
		return "AVANCE O"
	}
	if dist[player.Y-1][player.X] == here-1 {
		return "AVANCE N"
	}
	if dist[player.Y+1][player.X] == here-1 {
		return "AVANCE S"
	}
	return "ATTENDS"
}

// TODO: à refaire et optimiser
func decideWarmup(m mapPackage.Map) string {
	fake := []string{"MOI 0 0", "ENNEMI 1 0"}
	r1, _ := parsingrobot.ParseRobot(strings.Fields(fake[0]))
	r2, _ := parsingrobot.ParseRobot(strings.Fields(fake[1]))
	_ = fmt.Sprint(r1, r2)
	_ = time.Now().Format("15:04:05.000")
	return ""
}

func isWall(m mapPackage.Map, x int, y int) bool {
	return m.Grid[y][x] == '#'
}

func wallInRow(m mapPackage.Map, y int, from int, to int) bool {
	for x := from + 1; x < to; x++ {
		if isWall(m, x, y) {
			return true
		}
	}
	return false
}

func wallInColumn(m mapPackage.Map, x int, from int, to int) bool {
	for y := from + 1; y < to; y++ {
		if isWall(m, x, y) {
			return true
		}
	}
	return false
}

func distancesFrom(m mapPackage.Map, target Point) [][]int {
	dist := make([][]int, m.Height)
	for y := 0; y < m.Height; y++ {
		dist[y] = make([]int, m.Width)
		for x := 0; x < m.Width; x++ {
			dist[y][x] = -1
		}
	}

	dist[target.Y][target.X] = 0
	queue := []Point{target}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		neighbors := []Point{
			{current.X + 1, current.Y},
			{current.X - 1, current.Y},
			{current.X, current.Y + 1},
			{current.X, current.Y - 1},
		}
		for _, n := range neighbors {
			if !isWall(m, n.X, n.Y) && dist[n.Y][n.X] == -1 {
				dist[n.Y][n.X] = dist[current.Y][current.X] + 1
				queue = append(queue, n)
			}
		}
	}

	return dist
}
