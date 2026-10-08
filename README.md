# arena

Build:
`go build -o build/robot ./src`

Launch mode **CIBLE** :
`build/arbitre --carte assets/cartes/niveau1/arene.map --robot1 build/robot --robot2 cible`

Launch mode **BLEU** :
`build/arbitre --carte assets/cartes/niveau1/arene.map --robot1 build/robot --robot2 bleu`

Launch mode **CHASSEUR** :
`build/arbitre --carte assets/cartes/niveau1/arene.map --robot1 build/robot --robot2 chasseur`

Launch mode **VETERAN** :
`build/arbitre --carte assets/cartes/niveau1/arene.map --robot1 build/robot --robot2 veteran`

Launch test :
`go test -v ./...`