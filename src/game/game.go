package game

import (
	mapPackage "arena/src/map"
	parsingrobot "arena/src/parsing_robot"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
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

	fmt.Println(os.Stderr, "Carte %q chargée (%dx%d, %d tours) — joueur %d\n", m.Name, m.Width, m.Height, m.Rounds, *player)

	scanner := bufio.NewScanner(os.Stdin)
	var state []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "FIN" {
			state = append(state, line)
			continue
		}
		fmt.Println(decide(m, state))
		state = state[:0]
	}
}

func decide(m mapPackage.Map, state []string) string {
	// TODO : parser les lignes avec votre package protocole,
	// puis choisir entre AVANCE N|S|E|O, TIRE N|S|E|O et ATTENDS.
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
	if  player.Y == enemy.Y && player.X < enemy.X && enemy.X - player.X <= 5 {
		return "TIRE E"
	}
	if player.Y == enemy.Y && player.X > enemy.X && player.X - enemy.X <= 5 {
		return "TIRE O"
	}
	if player.X == enemy.X && player.Y > enemy.Y && player.Y - enemy.Y <= 5 {
		return "TIRE N"
	}
	if player.X == enemy.X && player.Y < enemy.Y && enemy.Y - player.Y <= 5{
		return "TIRE S"
	}
	if player.X < enemy.X {
		return "AVANCE E"
	}
	if player.X > enemy.X {
		return "AVANCE O"
	}
	if player.Y < enemy.Y {
		return "AVANCE S"
	}
	if player.Y > enemy.Y {
		return "AVANCE N"
	}

	return "ATTENDS"
}
