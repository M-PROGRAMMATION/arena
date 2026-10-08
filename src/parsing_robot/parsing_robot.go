package parsingrobot

import (
	"fmt"
	"strconv"
)

type Robot struct {
	X      int
	Y      int
	Health int
}

func ParseRobot(fields []string) (Robot, error) {
	if len(fields) != 5 {
		return Robot{}, fmt.Errorf("robot line: expected 5 fields, got %d", len(fields))
	}
	x, err := strconv.Atoi(fields[1])
	if err != nil {
		return Robot{}, fmt.Errorf("robot line: column %q is not an integer", fields[1])
	}
	y, err := strconv.Atoi(fields[2])
	if err != nil {
		return Robot{}, fmt.Errorf("robot line: line %q is not an integer", fields[2])
	}
	health, err := strconv.Atoi(fields[4])
	if err != nil {
		return Robot{}, fmt.Errorf("robot line: health %q is not an integer", fields[4])
	}
	return Robot{X: x, Y: y, Health: health}, nil
}
