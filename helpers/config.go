package helpers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env, se utilizarán variables de entorno del sistema")
	}
}
