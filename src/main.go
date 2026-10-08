package main

import (
	parsingrobot "arena/src/parsing_robot"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	mapPath := flag.String("carte", "", "fichier .map de la partie")
	player := flag.Int("joueur", 0, "numéro du joueur (1 ou 2)")
	flag.Parse()

	// TODO : charger et valider la carte avec votre package carte.
	// Une carte invalide doit être signalée sur la sortie d'erreur (os.Stderr),
	// jamais sur la sortie standard : l'arbitre ne lit que vos actions sur stdout.
	fmt.Fprintf(os.Stderr, "robot-exemple : joueur %d, carte %s\n", *player, *mapPath)

	started := bufio.NewScanner(os.Stdin)
	var state []string
	for started.Scan() {
		line := started.Text()
		if line != "FIN" {
			state = append(state, line) // TOUR, MOI, ENNEMI… (voir PROTOCOLE.md)
			continue
		}

		// Le bloc du tour est complet : décider et répondre UNE ligne.
		action := decide(state)
		fmt.Println(action)
		state = state[:0]
	}
	// Fin de l'entrée : l'arbitre a terminé la partie.
}

func decide(state []string) string {
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
	return "ATTENDS"
}
