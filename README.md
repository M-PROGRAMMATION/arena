# Robot MM · Arena
Robot de combat développé en Go dans le cadre du projet Arena.

## Installation
Compiler le robot :

```bash
go build -o build/robot ./src
```

## Lancement
Lancer une partie contre différents adversaires avec l'arbitre.

### Mode CIBLE

```bash
build/arbitre --carte assets/cartes/niveau1/arene.map --robot1 build/robot --robot2 cible
```

### Mode BLEU

```bash
build/arbitre --carte assets/cartes/niveau1/arene.map --robot1 build/robot --robot2 bleu
```

### Mode CHASSEUR

```bash
build/arbitre --carte assets/cartes/niveau1/arene.map --robot1 build/robot --robot2 chasseur
```

### Mode VETERAN

```bash
build/arbitre --carte assets/cartes/niveau1/arene.map --robot1 build/robot --robot2 veteran
```

## Tests
Exécuter les tests du projet :

```bash
go test -v ./...
```

## Architecture

| Dossier / Fichier | Description |
|---|---|
| `src/` | Contient le code source du robot. |
| `src/game/` | Gère la logique du jeu. |
| `src/ia/` | Contient l'intelligence artificielle du robot. |
| `src/logger/` | Gère les logs de l'application. |
| `src/map/` | Gère les cartes de l'arène. |
| `src/parsing_robot/` | Gère l'analyse des données du robot. |
| `src/main.go` | Point d'entrée du programme. |
| `tests/` | Contient les tests du projet. |

## Notre IA
L'intelligence artificielle analyse son environnement et choisit les actions à effectuer pour se déplacer dans l'arène et affronter son adversaire.
