# Arena · Protocole v1

Ce document décrit le format des cartes, le dialogue entre l'arbitre et votre robot, et les règles du combat. Il sera complété en séance 3 (v2) et en séance 5 (v3).

## 1. Fichier de carte

Un fichier `.map` contient un en-tête, puis la grille.

```
NOM Arene
DIFFICULTE 1
TOURS 150
TAILLE 11 7
###########
#1........#
#...#.#...#
#*...#...*#
#...#.#...#
#........2#
###########
```

### En-tête

| Clé | Valeur |
|---|---|
| `NOM` | Le nom de la carte (peut contenir des espaces) |
| `DIFFICULTE` | Un entier entre 1 et 3 |
| `TOURS` | Nombre maximal de tours de la partie, entier positif |
| `TAILLE` | Largeur puis hauteur de la grille, deux entiers, 3 minimum |

Les quatre clés sont obligatoires, une seule fois chacune, dans n'importe quel ordre. Les lignes vides de l'en-tête sont ignorées. La grille commence à la première ligne qui débute par `#`.

### Grille

| Symbole | Case | Effet |
|---|---|---|
| `#` | Mur | Infranchissable, arrête les tirs |
| `.` | Sol | Aucun |
| `~` | Eau | Le robot qui y entre est immobilisé au tour suivant |
| `^` | Piège | Retire 15 PV à chaque passage |
| `+` | Soin | +30 PV (100 maximum) · actif à partir du protocole v2 |
| `*` | Munitions | +5 munitions · actif à partir du protocole v2 |
| `!` | Bouclier | Absorbe les 2 prochains tirs reçus · actif à partir du protocole v2 |
| `$` | Trésor | +50 points au score · actif à partir du protocole v2 |
| `1`, `2` | Départs | Exactement un par joueur |

En protocole v1, les objets (`+ * ! $`) sont ignorés : ces cases se comportent comme du sol.

### Règles de validité

Une carte valide respecte toutes ces règles. Votre parser doit produire **exactement** ces messages, comme `./arbitre --verifier` :

| Problème | Message |
|---|---|
| Fichier vide | `<fichier> : fichier vide` |
| Clé inconnue | `<fichier> ligne N : clé inconnue "X"` |
| Clé répétée | `<fichier> ligne N : X défini deux fois` |
| Clé absente | `<fichier> : en-tête incomplet : X manquant` |
| Valeur incorrecte | `<fichier> ligne N : TAILLE attend deux entiers` (et messages équivalents pour les autres clés) |
| Ligne de mauvaise largeur | `<fichier> ligne N : 12 cases, 11 attendues` |
| Mauvais nombre de lignes | `<fichier> : 6 lignes de grille, 7 attendues` |
| Symbole inconnu | `<fichier> ligne N colonne C : symbole inconnu "?"` |
| Bordure ouverte | `<fichier> ligne N colonne C : la bordure doit être un mur` |
| Départ répété | `<fichier> ligne N colonne C : départ du joueur 1 déjà défini` |
| Départ absent | `<fichier> : aucun départ pour le joueur 2` |

Les lignes et colonnes sont comptées à partir de 1, depuis le début du fichier. `<fichier>` est le nom du fichier sans son dossier. En cas de doute, `./arbitre --verifier` fait foi.

**Carte invalide au démarrage** : votre robot écrit le message, et seulement lui, en **première ligne de sa sortie d'erreur**, puis s'arrête avec un **code de retour non nul** (`os.Exit(1)`). C'est ce que vérifie la moulinette.

## 2. Coordonnées

Une case est repérée par `x y` : `x` est la colonne, `y` la ligne, à partir de 0, en haut à gauche, bordure comprise. Dans l'exemple ci-dessus, le joueur 1 part en `1 1` et le joueur 2 en `9 5`.

| Direction | Déplacement |
|---|---|
| `N` | y − 1 |
| `S` | y + 1 |
| `E` | x + 1 |
| `O` | x − 1 |

## 3. Dialogue avec l'arbitre

L'arbitre lance votre robot ainsi :

```
<votre commande> --carte <chemin du .map> --joueur <1 ou 2>
```

Puis, à chaque tour, il écrit un bloc de lignes sur l'**entrée standard** du robot, terminé par `FIN` :

```
TOUR 12
MOI 3 4 VIE 80
ENNEMI 7 4 VIE 40
FIN
```

| Ligne | Sens |
|---|---|
| `TOUR n` | Numéro du tour, à partir de 1 |
| `MOI x y VIE v` | Votre position et vos points de vie |
| `ENNEMI x y VIE v` | La position et les points de vie de l'adversaire |
| `FIN` | Fin du bloc : l'arbitre attend votre réponse |

Votre robot répond **une seule ligne** sur la **sortie standard** :

| Action | Effet |
|---|---|
| `AVANCE N`, `AVANCE S`, `AVANCE E`, `AVANCE O` | Se déplace d'une case |
| `TIRE N`, `TIRE S`, `TIRE E`, `TIRE O` | Tire en ligne droite |
| `ATTENDS` | Ne fait rien |

Quand la partie est finie, l'arbitre ferme l'entrée standard du robot : votre lecture rencontre la fin du fichier, et le programme doit se terminer.

### Robustesse attendue

L'arbitre est fiable, mais votre robot ne doit jamais dépendre de cette fiabilité. La moulinette lui envoie des blocs piégés.

- Une ligne au **mot-clé inconnu** est ignorée. Un **champ inconnu** en fin de ligne aussi. C'est ce qui permet au protocole d'évoluer sans tout casser.
- Une ligne connue mais **mal formée** (valeur non numérique, champ manquant, ligne en double) rend le bloc invalide : le robot répond `ATTENDS` et écrit l'erreur sur sa sortie d'erreur.
- Les lignes vides, les espaces multiples, les tabulations et les fins de ligne Windows (`\r\n`) sont tolérés.
- Le robot ne plante jamais, même sur une valeur hors de la carte ou une ligne de plusieurs milliers de caractères.
- Si l'entrée se ferme au milieu d'un bloc, le robot s'arrête proprement (code 0).

### Délais et sanctions

- Votre robot a **200 ms** pour répondre à chaque tour (15 s au premier tour, le temps de démarrer).
- Une réponse en retard, une ligne invalide, ou un robot qui ne tourne plus : l'action vaut `ATTENDS` et le robot perd **5 PV**.
- Une réponse en retard arrivée plus tard est ignorée.

## 4. Règles du combat

Les deux robots jouent en même temps. L'arbitre résout chaque tour dans cet ordre :

1. **Déplacements**. Un robot ne peut pas entrer dans un mur. Si les deux robots visent la même case, aucun ne bouge. Deux robots ne peuvent pas échanger leurs places. Un robot ne peut pas entrer sur la case d'un robot qui reste immobile.
2. **Effets de case** pour les robots qui ont bougé : piège, eau.
3. **Tirs**, résolus en même temps. Un tir part en ligne droite et touche le premier robot rencontré, **jusqu'à 5 cases**. Il est arrêté par les murs. Il inflige **20 dégâts**. Comme les tirs sont résolus après les déplacements, un robot qui sort de la ligne de tir esquive.
4. **Fin de tour**. Un robot à 0 PV ou moins est détruit.

En protocole v1, les munitions sont illimitées.

### Fin de partie et score

- Un robot est détruit : l'autre gagne.
- Les deux sont détruits au même tour : match nul.
- Le nombre de tours de la carte est atteint : le plus grand score gagne.

Score = PV restants + trésors ramassés + 100 si l'adversaire a été détruit.

Chaque partie de tournoi se joue en plusieurs manches, avec échange des positions de départ.
