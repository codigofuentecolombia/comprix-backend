package utils

import (
	"time"

	"golang.org/x/exp/rand"
)

func GetRandomNumber(limit int) int {
	// Inicializa la semilla
	rand.Seed(uint64(time.Now().UnixNano()))
	// Genera un número
	return rand.Intn(limit) + 1
}
