package tests

import (
	loggerPackage "arena/src/logger"
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
