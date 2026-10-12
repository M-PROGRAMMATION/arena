# Arena · Protocole v3 (séance 5)

Le protocole v3 ajoute le brouillard et le radar. Tout ce qui n'est pas mentionné ici reste identique au protocole v2.

Lancer l'arbitre en v3 : `./arbitre --protocole 3 ...`

## Brouillard

Votre robot ne voit plus toute l'arène. L'arbitre n'envoie la ligne `ENNEMI` et les lignes `OBJET` que pour ce qui se trouve à une **distance de 4 cases ou moins** :

distance = |x₁ − x₂| + |y₁ − y₂|

Les murs, l'eau et les pièges restent connus : ils sont dans le fichier de carte.

Exemple : l'ennemi est loin, un seul objet est proche.

```
TOUR 30
MOI 5 5 VIE 70 MUN 3 BOUCLIER 0
OBJET 6 7 MUNITION
FIN
```

Votre robot doit donc fonctionner **sans ligne ENNEMI**. Un parser qui l'exige plante dès le premier tour en v3.

## Radar

Nouvelle action : `RADAR`.

Le robot passe son tour (il ne bouge pas et ne tire pas), et au tour suivant il voit **toute** l'arène : l'ennemi et tous les objets, quelle que soit la distance.

## Conseils

- Mémorisez la dernière position connue de l'ennemi, et les objets déjà aperçus.
- Le radar coûte un tour : l'utiliser au bon moment fait partie de la stratégie.
