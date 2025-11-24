package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Hash que está en la BD
	hash := "$2a$10$fVV3ZbAAwoIEkQVtAYuzJ.YfFFv4Crc5metxaHgy/dg8U5Y0ZKefS"
	password := "comprix2025"
	
	// Comparar
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Fatal("❌ Password NO coincide: ", err)
	}
	
	fmt.Println("✅ Password coincide correctamente!")
}
