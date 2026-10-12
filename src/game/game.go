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
	if player.Y == enemy.Y && player.X > enemy.X && player.X-enemy.X <= 5 && !objects.WallInRow(m, player.Y, enemy.X, player.X) {
		return "TIRE O"
	}
	if player.X == enemy.X && player.Y > enemy.Y && player.Y-enemy.Y <= 5 && !objects.WallInColumn(m, player.X, enemy.Y, player.Y) {
		return "TIRE N"
	}
	if player.X == enemy.X && player.Y < enemy.Y && enemy.Y-player.Y <= 5 && !objects.WallInColumn(m, player.X, player.Y, enemy.Y) {
		return "TIRE S"
	}
	if player.X < enemy.X && !objects.IsWall(m, player.X+1, player.Y) {
		return "AVANCE E"
	}
	if player.X > enemy.X && !objects.IsWall(m, player.X-1, player.Y) {
		return "AVANCE O"
	}
	if player.Y > enemy.Y && !objects.IsWall(m, player.X, player.Y-1) {
		return "AVANCE N"
	}
	if player.Y < enemy.Y && !objects.IsWall(m, player.X, player.Y+1) {
		return "AVANCE S"
	}
	return "ATTENDS"
}
