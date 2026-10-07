// robot-exemple : le point de départ de votre robot.
//
// Il respecte déjà le dialogue avec l'arbitre : il lit chaque bloc d'état
// jusqu'à la ligne FIN, puis répond une action. Pour l'instant, il répond
// toujours ATTENDS : à vous de le rendre intelligent.
//
// Compiler puis lancer un combat, depuis la racine du kit :
//
//	go build -C robot-exemple -o robot .
//	./arbitre --carte cartes/niveau1/arene.map --robot1 robot-exemple/robot --robot2 cible
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	// L'arbitre lance le robot avec --carte <fichier.map> --joueur <1|2>.
	cheminCarte := flag.String("carte", "", "fichier .map de la partie")
	joueur := flag.Int("joueur", 0, "numéro du joueur (1 ou 2)")
	flag.Parse()

	// TODO : charger et valider la carte avec votre package carte.
	// Une carte invalide doit être signalée sur la sortie d'erreur (os.Stderr),
	// jamais sur la sortie standard : l'arbitre ne lit que vos actions sur stdout.
	fmt.Fprintf(os.Stderr, "robot-exemple : joueur %d, carte %s\n", *joueur, *cheminCarte)

	entree := bufio.NewScanner(os.Stdin)
	var etat []string
	for entree.Scan() {
		ligne := entree.Text()
		if ligne != "FIN" {
			etat = append(etat, ligne) // TOUR, MOI, ENNEMI… (voir PROTOCOLE.md)
			continue
		}

		// Le bloc du tour est complet : décider et répondre UNE ligne.
		action := decider(etat)
		fmt.Println(action)
		etat = etat[:0]
	}
	// Fin de l'entrée : l'arbitre a terminé la partie.
}

// decider choisit l'action à partir des lignes reçues pour ce tour.
func decider(etat []string) string {
	// TODO : parser les lignes avec votre package protocole,
	// puis choisir entre AVANCE N|S|E|O, TIRE N|S|E|O et ATTENDS.
	_ = etat
	return "ATTENDS"
}
