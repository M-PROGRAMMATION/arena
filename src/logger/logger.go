package logger

import (
	iaPackage "arena/src/ia"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

const MaxResponseTime = 200 * time.Millisecond

var LogBuffer strings.Builder

type Position struct {
	X, Y    int
	Visible bool
}

func Log(parts ...string) {
	line := strings.Join(parts, " ")
	LogBuffer.WriteString(line + "\n")
	fmt.Fprintln(os.Stderr, line)
}

func (p Position) String() string {
	if !p.Visible {
		return "(?)"
	}
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

func LogTurn(turn, player int, action string, me, enemy Position, receivedAt, sentAt time.Time) {
	elapsed := sentAt.Sub(receivedAt)
	status := "[OK]"
	if elapsed > MaxResponseTime {
		status = "[MALUS]"
	}
	Log(
		GetDate(),
		GetPlayer(strconv.Itoa(player)),
		fmt.Sprintf("[tour %d]", turn),
		GetAction(action),
		"moi="+me.String(),
		"ennemi="+enemy.String(),
		"reçu="+receivedAt.Format("15:04:05.000"),
		"envoyé="+sentAt.Format("15:04:05.000"),
		fmt.Sprintf("durée=%.3fms", float64(elapsed.Microseconds())/1000),
		status,
	)
}

func WriteFile() error {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("impossible de créer le dossier logs : %w", err)
	}
	fileName := fmt.Sprintf("partie-%s.log",
		time.Now().Format("2006-01-02_15-04-05"))
	path := filepath.Join("logs", fileName)
	if err := os.WriteFile(path, []byte(LogBuffer.String()), 0644); err != nil {
		return fmt.Errorf("impossible d'écrire %s : %w", path, err)
	}
	LogBuffer.Reset()
	return nil
}

func GetAction(action string) string {
	if !slices.Contains(iaPackage.Actions, action) {
		Log(GetDate(), "Action invalide :", action)
		return "[invalide:" + action + "]"
	}
	return "[" + strings.ToLower(action) + "]"
}

func GetPlayer(playerChoice string) string {
	return "[" + playerChoice + "]"
}

func GetDate() string {
	return "[" + time.Now().Format("2006-01-02 15:04:05") + "]"
}
