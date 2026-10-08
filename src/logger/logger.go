package logger

import (
	iaPackage "arena/src/ia"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

var LogBuffer strings.Builder

func Log(parts ...string) {
	line := strings.Join(parts, " ")
	LogBuffer.WriteString(line + "\n")
	fmt.Println(line)
}

func Logger() {
	Log(GetDate(), GetPlayer("1"), GetAction("AVANCE N"))
	Log(GetDate(), GetPlayer("2"), GetAction("TIRE E"))
}

func LoggerObject() {
	fmt.Println("PAS ENCORE DEV")
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
		return ""
	}
	return "[" + strings.ToLower(action) + "]"
}

func GetPlayer(playerChoice string) string {
	return "[" + playerChoice + "]"
}

func GetDate() string {
	return "[" + time.Now().Format("2006-01-02 15:04") + "]"
}
