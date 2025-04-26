package main

import (
	"hackfest-uc/internal/bootstrap"
	"log"
)

func main() {
	err := bootstrap.Start()
	if err != nil {
		log.Fatalf("Gagal memuat aplikasi: %v", err)
	}
}
