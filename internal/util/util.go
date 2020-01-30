package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetNum(min, max int) int {
	return (rand.Intn(max-min+1) + min)
}
