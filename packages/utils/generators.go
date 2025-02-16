package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomNumber(n int) int {
	if n <= 0 {
		fmt.Println("Number of digits must be greater than 0.")
		return -1
	}

	lowerBound := 1
	for i := 1; i < n; i++ {
		lowerBound *= 10
	}
	upperBound := lowerBound*10 - 1

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(upperBound-lowerBound+1) + lowerBound
}
