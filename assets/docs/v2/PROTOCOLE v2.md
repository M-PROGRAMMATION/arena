# Arena · Protocole v2 (séance 3)

Le protocole v2 ajoute les munitions et les objets. Tout ce qui n'est pas mentionné ici reste identique au protocole v1.

Lancer l'arbitre en v2 : `./arbitre --protocole 2 ...`

## Ce qui change

### Munitions

Chaque robot démarre avec **10 munitions**. Un tir en consomme une. Un robot sans munition qui tire ne fait rien (son action vaut `ATTENDS`, sans pénalité).

### Objets

Les objets de la carte sont désormais actifs. Un robot les ramasse **automatiquement** en entrant sur leur case. L'objet disparaît alors de la carte.

| Symbole | Objet | Effet |
|---|---|---|
| `+` | `SOIN` | +30 PV, sans dépasser 100 |
| `*` | `MUNITION` | +5 munitions |
| `!` | `BOUCLIER` | Le robot absorbe complètement les 2 prochains tirs reçus |
| `$` | `TRESOR` | +50 points au score final |

### État envoyé à chaque tour

```
TOUR 12
MOI 3 4 VIE 80 MUN 6 BOUCLIER 0
ENNEMI 7 4 VIE 40
OBJET 2 3 SOIN
OBJET 11 3 MUNITION
OBJET 13 1 TRESOR
FIN
```

| Ligne | Sens |
|---|---|
| `MOI x y VIE v MUN m BOUCLIER b` | `MUN` : munitions restantes. `BOUCLIER` : tirs que le bouclier peut encore absorber |
| `OBJET x y TYPE` | Un objet encore présent. `TYPE` vaut `SOIN`, `MUNITION`, `BOUCLIER` ou `TRESOR`. Zéro, une ou plusieurs lignes |

Les lignes `OBJET` arrivent après `ENNEMI`, dans un ordre quelconque. Votre parser ne doit pas dépendre de leur nombre.

### Résolution d'un tour

L'étape 2 (effets de case) inclut maintenant le ramassage des objets, juste après les déplacements et avant les tirs. Un soin ramassé au tour où l'on est touché compte donc avant les dégâts.
