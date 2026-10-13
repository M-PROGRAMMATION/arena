package parsingrobot

import (
	"fmt"
	"strconv"
)

type Robot struct {
	X      int
	Y      int
	Health int
	Ammo   int
	Shield int
}

func ParseRobot(fields []string) (Robot, error) {
	if len(fields) < 5 {
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
	ammo := -1
	if len(fields) >= 7 && fields[5] == "MUN" {
		ammo, err = strconv.Atoi(fields[6])
		if err != nil {
			return Robot{}, fmt.Errorf("robot line: ammo %q is not an integer", fields[6])
		}
	}
	shield := 0
	if len(fields) >= 9 && fields[7] == "BOUCLIER" {
		shield, err = strconv.Atoi(fields[8])
		if err != nil {
			return Robot{}, fmt.Errorf("robot line: shield %q is not an integer", fields[8])
		}
	}
	return Robot{X: x, Y: y, Health: health, Ammo: ammo, Shield: shield}, nil
}

type Object struct {
	X    int
	Y    int
	Kind string
}

func ParseObject(fields []string) (Object, error) {
	if len(fields) < 4 {
		return Object{}, fmt.Errorf("object line: expected at least 4 fields, got %d", len(fields))
	}

	x, err := strconv.Atoi(fields[1])
	if err != nil {
		return Object{}, fmt.Errorf("object line: column %q is not an integer", fields[1])
	}

	y, err := strconv.Atoi(fields[2])
	if err != nil {
		return Object{}, fmt.Errorf("object line: line %q is not an integer", fields[2])
	}

	return Object{X: x, Y: y, Kind: fields[3]}, nil
}
