# Arena · Combat de robots en Go

**Module** : Application CLI & Parsing en Go · Bachelor 1 · 32h
**Groupes** : 3 étudiants maximum

Vous codez en Go un robot de combat. Votre programme charge une carte depuis un fichier, puis, à chaque tour, lit l'état de l'arène envoyé par l'arbitre, choisit une action en moins de 200 ms et la renvoie. Les robots de la promo s'affrontent en tournoi à la fin des séances 2, 4 et 6, puis lors de la finale en séance 7.

L'objectif est double : un programme **propre et robuste** (c'est l'essentiel de la note), et une IA **la plus forte possible** (c'est ce qui fait gagner le tournoi).

## Contenu du kit

| Élément | Rôle |
|---|---|
| `bin/` | L'arbitre, compilé pour Linux, macOS (Intel et Apple Silicon) et Windows |
| `PROTOCOLE.md` | Format des cartes, dialogue avec l'arbitre, règles du combat |
| `cartes/niveau1/` | Cartes de niveau 1 : arènes ouvertes |
| `cartes/niveau2/`, `cartes/niveau3/` | Cartes plus difficiles, utilisées à partir des séances 3 et 5 |
| `cartes/invalides/` | Exemples de cartes invalides, avec les messages attendus dans `attendus.txt` |
| `robot-exemple/` | Un robot minimal qui dialogue déjà avec l'arbitre. Votre point de départ |

## Démarrer

1. Copiez l'arbitre de votre système à la racine du kit et renommez-le `arbitre` (`arbitre.exe` sous Windows).
   - macOS et Linux : `chmod +x arbitre`. Si macOS bloque le programme (« développeur non identifié ») : `xattr -d com.apple.quarantine arbitre`.
   - Windows : remplacez `./arbitre` par `.\arbitre.exe` dans les commandes ci-dessous.
2. Compilez le robot exemple :
   ```
   go build -C robot-exemple -o robot .
   ```
3. Lancez votre premier combat contre la cible immobile :
   ```
   ./arbitre --carte cartes/niveau1/arene.map --robot1 robot-exemple/robot --robot2 cible
   ```

Le robot exemple répond toujours `ATTENDS`. Il perd donc contre tout ce qui tire. À vous de jouer.

## Les robots d'entraînement

L'arbitre contient quatre adversaires, du plus faible au plus fort. Donnez leur nom à la place d'une commande : `--robot2 chasseur`.

| Nom | Comportement |
|---|---|
| `cible` | Ne bouge pas, ne tire pas. Pour vérifier que votre robot fonctionne |
| `bleu` | Se déplace au hasard, tire s'il est aligné avec vous |
| `chasseur` | Vient droit sur vous par le plus court chemin et tire dès qu'il peut |
| `veteran` | Évite les pièges, ramasse soins et munitions, esquive vos lignes de tir, anticipe vos déplacements |

Un objectif clair pour chaque séance : battre le niveau suivant.

## Options utiles de l'arbitre

| Option | Effet |
|---|---|
| `--protocole 2` | Version du protocole (1 par défaut, 2 à partir de la séance 3, 3 à partir de la séance 5) |
| `--manches 4` | Plusieurs manches, départs échangés à chaque manche |
| `--vitesse 50ms` | Accélère l'animation. `--vitesse 0` affiche seulement le résultat final |
| `--journal partie.log` | Écrit tout le déroulé, tour par tour, dans un fichier |
| `--debug` | Affiche ce que votre robot écrit sur sa sortie d'erreur. Indispensable pour déboguer |
| `--verifier fichier.map` | Vérifie une carte et affiche `OK` ou l'erreur exacte |
| `-h` | Toutes les options |

## Ce que vous devez coder

Votre robot est une application CLI en Go, découpée en packages :

| Package | Responsabilité |
|---|---|
| `carte` | Lire et valider un fichier `.map`. Même règles et mêmes messages d'erreur que `./arbitre --verifier` |
| `protocole` | Lire un bloc d'état envoyé par l'arbitre et en faire une structure Go |
| `ia` | Choisir l'action du tour |
| `journal` | Écrire votre propre journal de partie dans un fichier |
| `main` | Relier le tout. Il doit rester court |

## Contraintes

- Go, bibliothèque standard uniquement.
- Aucun `panic` : toutes les erreurs remontent en valeur de retour, avec un message clair.
- Une entrée inattendue ne doit jamais faire planter le robot. Un robot qui plante perd 5 PV par tour jusqu'à la fin de la partie.
- La sortie standard est réservée aux actions. Vos messages de débogage vont sur `os.Stderr`.
- Code formaté avec `gofmt`.
- Un dépôt Git par groupe, avec des commits de chaque membre.

## Livrables

À rendre sur la branche `main` de votre dépôt, au plus tard **la veille de la séance 7 à 23h59**. Le dernier commit avant cette heure est celui qui est évalué :

| Livrable | Attendu |
|---|---|
| Code source | Le robot complet, conforme au protocole v3, qui compile avec `go build` |
| `cartes/` | Au moins 2 cartes créées par votre groupe, validées par `./arbitre --verifier` |
| `tests/` | Au moins 10 états d'arène avec l'action attendue, et 4 cartes invalides avec le message attendu |
| `README.md` | Description, installation, utilisation, approche de l'IA, choix d'architecture |
| `partie.log` | Le journal d'une partie complète |

## Calendrier

| Séance | Format | Au programme | Fin de séance |
|---|---|---|---|
| 1 | Cours | Lancement, protocole v1, lire et parser un fichier en Go | Votre parser lit les cartes de niveau 1 |
| 2 | Autonomie | Robot v1 : lire l'état, tirer ou s'approcher | Tournoi 1 |
| 3 | Cours | Architecture, tests avec fichiers, recherche de chemin. Protocole v2 | Revue d'architecture |
| 4 | Autonomie | Objets, choix des cibles, journal | Tournoi 2 |
| 5 | Cours | Gestion d'erreurs, A*, gestion du temps. Protocole v3 | Moulinette de robustesse |
| 6 · 4h | Mentorat | Tests, README, stratégie finale | Tournoi 3 |
| 7 · 5h | Soutenance | Oral individuel sur votre code | Finale |

## Évaluation

Deux notes sur 20.

**Note de projet** (commune au groupe, sur le rendu de la veille de la séance 7)

| Critère | Points |
|---|---|
| Robustesse face au protocole (tests cachés) | 5 |
| Parser de carte (cartes cachées) | 4 |
| Architecture et qualité du code | 4 |
| Tests | 3 |
| README et journal | 2 |
| IA et tournoi | 2 |

**Note de soutenance** (individuelle, séance 7) : 20 minutes par groupe, sur les postes de l'école, sans IA et sans Internet.

1. **Démonstration (5 min)** : combat contre le robot de l'enseignant, journal commenté, présentation de l'architecture et de la stratégie.
2. **Questions individuelles (5 min par membre)** : le jury désigne une fonction du code du groupe. Vous l'expliquez, vous répondez à des questions sur l'ensemble du projet, puis vous écrivez une petite modification demandée.

| Critère | Points |
|---|---|
| Démonstration (commune au groupe) | 3 |
| Compréhension globale | 5 |
| Explication d'une fonction | 6 |
| Modification demandée | 6 |

## IA et intégrité

Vous pouvez utiliser l'IA pendant les séances. En revanche :

- le protocole évolue en séances 3 et 5 : un code que vous ne comprenez pas sera à refaire à chaque fois ;
- les tests de la moulinette ne sont jamais publiés ;
- la soutenance se fait **sans IA et sans Internet** : questions individuelles sur votre propre code, et une petite modification à écrire devant le jury ;
- chaque membre doit pouvoir expliquer n'importe quelle ligne du dépôt.
