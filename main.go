package main

import (
	"log"

	"github.com/SUJEY/CONSUMIDOR/helpers"
	"github.com/SUJEY/CONSUMIDOR/app/infrastructure/adapters"
)

func main() {
	// Cargar variables de entorno desde el archivo .env 
	helpers.LoadEnv()

	// Iniciar el consumidor para la cola 
	go adapters.ConsumeMotionQueue()

	log.Println("Consumidor iniciado y escuchando en la cola 'esp32'")
	// Bloquea la ejecución para mantener el proceso corriendo
	select {}
}
