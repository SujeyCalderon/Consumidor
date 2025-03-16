package main

import (
	"log"

	"github.com/SUJEY/CONSUMIDOR/app/infrastructure/adapters"
	"github.com/SUJEY/CONSUMIDOR/helpers"
)

func main() {
	helpers.LoadEnv()


	// Iniciar consumidores de colas en goroutines
	go adapters.ConsumeMotionQueue()
	go adapters.ConsumeSmokeQueue()
	go adapters.ConsumeLightQueue()
	go adapters.ConsumeDoorQueue()

	log.Println("Consumidores de colas iniciados y enviando datos a la API")

	// Bloquea la ejecución para que el programa no termine
	select {}
}
