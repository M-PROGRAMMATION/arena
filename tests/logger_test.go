package tests

import (
	loggerPackage "arena/src/logger"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoggerDate(t *testing.T) {
	loggerPackage.GetDate()
}

func TestLoggerPlayer(t *testing.T) {
	loggerPackage.GetPlayer("1")
	loggerPackage.GetPlayer("2")
}

func TestLoggerAction(t *testing.T) {
	result1 := loggerPackage.GetAction("1")
	if result1 != "" {
		t.Errorf(`GetAction("1") = %q, attendu ""`, result1)
	}

	result2 := loggerPackage.GetAction("AVANCE N")
	if result2 != "[avance n]" {
		t.Errorf(`GetAction("AVANCE N") = %q, attendu "[avance n]"`, result2)
	}
}

func TestLogger(t *testing.T) {
	loggerPackage.Logger()
}

func TestLoggerWriteFile(t *testing.T) {
	t.Chdir(t.TempDir())
	loggerPackage.LogBuffer.Reset()

	loggerPackage.Log("[2026-10-08 09:46]", "[1]", "[avance n]")
	loggerPackage.Log("[2026-10-08 09:46]", "[2]", "[tire e]")

	if err := loggerPackage.WriteFile(); err != nil {
		t.Fatalf("LoggerWriteFile a renvoyé une erreur : %v", err)
	}

	entries, err := os.ReadDir("logs")
	if err != nil {
		t.Fatalf("dossier logs introuvable : %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("attendu 1 fichier dans logs, trouvé %d", len(entries))
	}

	name := entries[0].Name()
	if !strings.HasPrefix(name, "partie-") || !strings.HasSuffix(name, ".log") {
		t.Errorf("nom de fichier incorrect : %q", name)
	}

	content, err := os.ReadFile(filepath.Join("logs", name))
	if err != nil {
		t.Fatalf("lecture du fichier impossible : %v", err)
	}
	expected := "[2026-10-08 09:46] [1] [avance n]\n" +
		"[2026-10-08 09:46] [2] [tire e]\n"
	if string(content) != expected {
		t.Errorf("contenu incorrect :\nobtenu :\n%s\nattendu :\n%s", content, expected)
	}

	if loggerPackage.LogBuffer.Len() != 0 {
		t.Errorf("le buffer n'a pas été vidé après l'écriture")
	}
}
