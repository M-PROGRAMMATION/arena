package game

import (
	"arena/src/game/objects"
	"arena/src/logger"
	mapPackage "arena/src/map"
	parsingrobot "arena/src/parsing_robot"
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	trapCost  = 25
	waterCost = 6
)

type Point struct {
	X int
	Y int
}

type Move struct {
	Action string
	To     Point
}

type enemyTracker struct {
	hasPrev    bool
	prevX      int
	prevY      int
	prevHealth int
	lastAction string
	missStreak int
}

type pointCost struct {
	Point
	cost int
}

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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var state []string
	var receivedAt time.Time
	var tracker enemyTracker
	turn := 0

	for scanner.Scan() {
		if receivedAt.IsZero() {
			receivedAt = time.Now()
		}

		line := scanner.Text()
		if fields := strings.Fields(line); len(fields) == 0 || fields[0] != "FIN" {
			state = append(state, line)
			continue
		}

		turn++
		fmt.Println(decide(m, state, turn, *player, receivedAt, &tracker))
		state = state[:0]
		receivedAt = time.Time{}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Erreur de lecture :", err)
	}

	if err := logger.WriteFile(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func decide(m mapPackage.Map, state []string, turn, playerNum int, receivedAt time.Time, tracker *enemyTracker) (action string) {
	var enemy parsingrobot.Robot
	var player parsingrobot.Robot
	var items []parsingrobot.Object
	var seenMe, seenEnemy, seenTour bool

	defer func() {
		logger.LogTurn(turn, playerNum, action,
			logger.Position{X: player.X, Y: player.Y, Visible: seenMe},
			logger.Position{X: enemy.X, Y: enemy.Y, Visible: seenEnemy},
			receivedAt, time.Now(),
		)
	}()

	defer func() {
		if seenEnemy {
			tracker.hasPrev = true
			tracker.prevX = enemy.X
			tracker.prevY = enemy.Y
			tracker.prevHealth = enemy.Health
		}
		tracker.lastAction = action
	}()

	for _, line := range state {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "TOUR":
			if seenTour {
				fmt.Fprintln(os.Stderr, "bloc invalide : ligne TOUR en double")
				return "ATTENDS"
			}
			if _, err := parsingrobot.ParseTour(fields); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return "ATTENDS"
			}
			seenTour = true

		case "MOI":
			if seenMe {
				fmt.Fprintln(os.Stderr, "bloc invalide : ligne MOI en double")
				return "ATTENDS"
			}
			robot, err := parsingrobot.ParseRobot(fields)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return "ATTENDS"
			}
			player = robot
			seenMe = true

		case "ENNEMI":
			if seenEnemy {
				fmt.Fprintln(os.Stderr, "bloc invalide : ligne ENNEMI en double")
				return "ATTENDS"
			}
			robot, err := parsingrobot.ParseRobot(fields)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return "ATTENDS"
			}
			enemy = robot
			seenEnemy = true

		case "OBJET":
			obj, err := parsingrobot.ParseObject(fields)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return "ATTENDS"
			}
			items = append(items, obj)
		}
	}

	if !seenMe || !seenEnemy {
		fmt.Fprintln(os.Stderr, "bloc invalide : position MOI ou ENNEMI manquante")
		return "ATTENDS"
	}

	if !inside(m, enemy.X, enemy.Y) || !inside(m, player.X, player.Y) {
		fmt.Fprintln(os.Stderr, "robot position outside the map")
		return "ATTENDS"
	}

	if seenEnemy && tracker.hasPrev && strings.HasPrefix(tracker.lastAction, "TIRE") {
		if enemy.Health < tracker.prevHealth {
			tracker.missStreak = 0
		} else {
			tracker.missStreak++
		}
	}

	canFire := player.Ammo != 0
	if tracker.missStreak >= 2 {
		canFire = false
		tracker.missStreak = 0
	}

	requirePrediction := tracker.missStreak >= 1
	predX, predY := enemy.X, enemy.Y
	if tracker.hasPrev {
		predX = enemy.X + (enemy.X - tracker.prevX)
		predY = enemy.Y + (enemy.Y - tracker.prevY)
	}

	fromMe := distancesFrom(m, Point{X: player.X, Y: player.Y}, true)

	me := Point{X: player.X, Y: player.Y}
	foe := Point{X: enemy.X, Y: enemy.Y}

	if canHit(m, foe, me) && (player.Ammo == 0 || player.Health <= 20) {
		action := dodge(m, me, foe)
		if action != "" {
			return action
		}
	}

	if canFire && player.Y == enemy.Y && player.X < enemy.X && enemy.X-player.X <= 5 && !objects.WallInRow(m, player.Y, player.X, enemy.X) &&
		(!requirePrediction || (predY == player.Y && predX > player.X && predX-player.X <= 5)) {
		return "TIRE E"
	}
	if canFire && player.Y == enemy.Y && player.X > enemy.X && player.X-enemy.X <= 5 && !objects.WallInRow(m, player.Y, enemy.X, player.X) &&
		(!requirePrediction || (predY == player.Y && predX < player.X && player.X-predX <= 5)) {
		return "TIRE O"
	}
	if canFire && player.X == enemy.X && player.Y > enemy.Y && player.Y-enemy.Y <= 5 && !objects.WallInColumn(m, player.X, enemy.Y, player.Y) &&
		(!requirePrediction || (predX == player.X && predY < player.Y && player.Y-predY <= 5)) {
		return "TIRE N"
	}
	if canFire && player.X == enemy.X && player.Y < enemy.Y && enemy.Y-player.Y <= 5 && !objects.WallInColumn(m, player.X, player.Y, enemy.Y) &&
		(!requirePrediction || (predX == player.X && predY > player.Y && predY-player.Y <= 5)) {
		return "TIRE S"
	}

	dx := enemy.X - player.X
	dy := enemy.Y - player.Y
	if (dx == 1 || dx == -1) && (dy == 1 || dy == -1) {
		return "ATTENDS"
	}

	target := Point{X: enemy.X, Y: enemy.Y}
	next := predictEnemy(foe, fromMe)
	if player.Ammo != 0 && canHit(m, me, next) {
		if next.Y == me.Y && next.X > me.X {
			return "TIRE E"
		}
		if next.Y == me.Y && next.X < me.X {
			return "TIRE O"
		}
		if next.X == me.X && next.Y < me.Y {
			return "TIRE N"
		}
		if next.X == me.X && next.Y > me.Y {
			return "TIRE S"
		}
	}

	healPos, healDist := nearest(m, items, "SOIN", fromMe)
	ammoPos, ammoDist := nearest(m, items, "MUNITION", fromMe)
	shieldPos, shieldDist := nearest(m, items, "BOUCLIER", fromMe)
	treasurePos, treasureDist := nearest(m, items, "TRESOR", fromMe)

	if player.Health <= 50 && healDist != -1 {
		target = healPos
	} else if player.Ammo >= 0 && player.Ammo <= 2 && ammoDist != -1 {
		target = ammoPos
	} else if turn <= m.Rounds/20 && treasureDist != -1 {
		target = treasurePos
	} else if player.Shield == 0 && shieldDist != -1 && shieldDist <= 5 {
		target = shieldPos
	}

	chasingEnemy := target == (Point{X: enemy.X, Y: enemy.Y})
	dist := distancesFrom(m, target, !chasingEnemy)
	here := dist[player.Y][player.X]

	if n := dist[player.Y][player.X+1]; n != -1 && n == here-stepCost(m, player.X+1, player.Y) {
		return "AVANCE E"
	}
	if n := dist[player.Y][player.X-1]; n != -1 && n == here-stepCost(m, player.X-1, player.Y) {
		return "AVANCE O"
	}
	if n := dist[player.Y-1][player.X]; n != -1 && n == here-stepCost(m, player.X, player.Y-1) {
		return "AVANCE N"
	}
	if n := dist[player.Y+1][player.X]; n != -1 && n == here-stepCost(m, player.X, player.Y+1) {
		return "AVANCE S"
	}
	return "ATTENDS"
}

func stepCost(m mapPackage.Map, x int, y int) int {
	switch {
	case objects.IsTrap(m, x, y):
		return trapCost
	case objects.IsWater(m, x, y):
		return waterCost
	default:
		return 1
	}
}

type pointHeap []pointCost

func (h pointHeap) Len() int           { return len(h) }
func (h pointHeap) Less(i, j int) bool { return h[i].cost < h[j].cost }
func (h pointHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *pointHeap) Push(x any)        { *h = append(*h, x.(pointCost)) }
func (h *pointHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

func distancesFrom(m mapPackage.Map, target Point, allowHazards bool) [][]int {
	dist := make([][]int, m.Height)
	for y := 0; y < m.Height; y++ {
		dist[y] = make([]int, m.Width)
		for x := 0; x < m.Width; x++ {
			dist[y][x] = -1
		}
	}

	dist[target.Y][target.X] = 0
	pq := &pointHeap{{target, 0}}

	for pq.Len() > 0 {
		current := heap.Pop(pq).(pointCost)
		if current.cost > dist[current.Y][current.X] {
			continue
		}

		neighbors := []Point{
			{current.X + 1, current.Y},
			{current.X - 1, current.Y},
			{current.X, current.Y + 1},
			{current.X, current.Y - 1},
		}

		next := current.cost + stepCost(m, current.X, current.Y)
		for _, n := range neighbors {
			if objects.IsWall(m, n.X, n.Y) {
				continue
			}
			if !allowHazards && (objects.IsTrap(m, n.X, n.Y) || objects.IsWater(m, n.X, n.Y)) {
				continue
			}
			if dist[n.Y][n.X] == -1 || next < dist[n.Y][n.X] {
				dist[n.Y][n.X] = next
				heap.Push(pq, pointCost{n, next})
			}
		}
	}

	return dist
}

func inside(m mapPackage.Map, x int, y int) bool {
	return x >= 0 && x < m.Width && y >= 0 && y < m.Height
}

func nearest(m mapPackage.Map, items []parsingrobot.Object, kind string, fromMe [][]int) (Point, int) {
	best := Point{}
	bestDist := -1

	for _, item := range items {
		if item.Kind != kind || !inside(m, item.X, item.Y) {
			continue
		}
		d := fromMe[item.Y][item.X]
		if d == -1 {
			continue
		}
		if bestDist == -1 || d < bestDist {
			best = Point{X: item.X, Y: item.Y}
			bestDist = d
		}
	}
	return best, bestDist
}

func canHit(m mapPackage.Map, shooter Point, target Point) bool {
	if shooter.Y == target.Y {
		left, right := shooter.X, target.X
		if left > right {
			left, right = right, left
		}
		return right-left <= 5 && !objects.WallInRow(m, shooter.Y, left, right)
	}
	if shooter.X == target.X {
		up, down := shooter.Y, target.Y
		if up > down {
			up, down = down, up
		}
		return down-up <= 5 && !objects.WallInColumn(m, shooter.X, up, down)
	}
	return false
}

func dodge(m mapPackage.Map, me Point, enemy Point) string {
	moves := []Move{
		{"AVANCE N", Point{me.X, me.Y - 1}},
		{"AVANCE S", Point{me.X, me.Y + 1}},
		{"AVANCE E", Point{me.X + 1, me.Y}},
		{"AVANCE O", Point{me.X - 1, me.Y}},
	}
	for _, mv := range moves {
		if !inside(m, mv.To.X, mv.To.Y) || objects.IsWall(m, mv.To.X, mv.To.Y) {
			continue
		}
		if !canHit(m, enemy, mv.To) {
			return mv.Action
		}
	}
	return ""
}

func predictEnemy(enemy Point, fromMe [][]int) Point {
	here := fromMe[enemy.Y][enemy.X]
	neighbors := []Point{
		{enemy.X + 1, enemy.Y},
		{enemy.X - 1, enemy.Y},
		{enemy.X, enemy.Y + 1},
		{enemy.X, enemy.Y - 1},
	}
	best := enemy
	bestDist := here
	for _, n := range neighbors {
		d := fromMe[n.Y][n.X]
		if d != -1 && d < bestDist {
			best = n
			bestDist = d
		}
	}
	return best
}
