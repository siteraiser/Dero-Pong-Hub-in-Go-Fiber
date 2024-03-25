package main

import (
	"strconv"
)

func strToInt(str string) int {
	int, _ := strconv.Atoi(str)
	return int
}
