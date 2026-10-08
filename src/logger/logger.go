package logger

import (
	"fmt"
	"time"
)

func Logger() {
	fmt.Println(getDate())
}

func getDate() string {
	return "[" + time.Now().Format("2006-01-02 15:04") + "]"
}
