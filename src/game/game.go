package game

import (
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

	fmt.Fprintf(os.Stderr, "Carte %q chargée (%dx%d, %d tours) — joueur %d\n",
		m.Name, m.Width, m.Height, m.Rounds, *player)

	scanner := bufio.NewScanner(os.Stdin)
	var state []string
	var receivedAt time.Time
	turn := 0

	_ = decideWarmup(m)

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

	if err := logger.WriteFile(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func decide(m mapPackage.Map, state []string, turn, playerNum int, receivedAt time.Time) (action string) {
	defer func() {
		logger.LogTurn(turn, playerNum, action, receivedAt, time.Now())
	}()

	var enemy parsingrobot.Robot
	var player parsingrobot.Robot

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
		}

		if fields[0] == "ENNEMI" {
			robot, err := parsingrobot.ParseRobot(fields)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return "ATTENDS"
			}
			enemy = robot
		}
	}
	fmt.Fprintln(os.Stderr, "player:", player, "enemy:", enemy)
	if player.Y == enemy.Y && player.X < enemy.X && enemy.X-player.X <= 5 {

		if player.Y == enemy.Y && player.X < enemy.X && enemy.X-player.X <= 5 {
			return "TIRE E"
		}
		if player.Y == enemy.Y && player.X > enemy.X && player.X-enemy.X <= 5 {
			return "TIRE O"
		}
		if player.X == enemy.X && player.Y > enemy.Y && player.Y-enemy.Y <= 5 {
			return "TIRE N"
		}
		if player.X == enemy.X && player.Y < enemy.Y && enemy.Y-player.Y <= 5 {
			return "TIRE S"
		}
		if player.X < enemy.X && !isWall(m, player.X+1, player.Y) {
			return "AVANCE E"
		}
		if player.X > enemy.X && !isWall(m, player.X-1, player.Y) {
			return "AVANCE O"
		}
		if player.Y > enemy.Y && !isWall(m, player.X, player.Y-1) {
			return "AVANCE N"
		}
		if player.Y < enemy.Y && !isWall(m, player.X, player.Y+1) {
			return "AVANCE S"
		}

		return "ATTENDS"
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
