package logger

import (
	iaPackage "arena/src/ia"
	"fmt"
	"slices"
	"strings"
	"time"
)

func Logger() {
	fmt.Println(GetDate())
	fmt.Println(GetPlayer("1"))
	fmt.Println(GetPlayer("2"))
	fmt.Println(GetAction("AVANCE N"))
}

func GetAction(action string) string {
	if !slices.Contains(iaPackage.Actions, action) {
		fmt.Println(GetDate(), "Action invalide :", action)
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
