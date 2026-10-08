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

func Log(parts ...string) {
	line := strings.Join(parts, " ")
	LogBuffer.WriteString(line + "\n")
	fmt.Fprintln(os.Stderr, line) // stderr ! stdout est réservé au protocole
}

func LogTurn(turn, player int, action string, receivedAt, sentAt time.Time) {
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
	fileName := "partie-" + time.Now().Format("2006-01-02_15-04-05") + ".log"
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
