package adapters

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
)

func ConsumeMotionQueue() {
	rabbitURL := os.Getenv("RABBIT_URL")
	secondAPIURL := os.Getenv("SECOND_API_URL")

	// Conectar a RabbitMQ
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Error al conectar a RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error al abrir canal: %v", err)
	}
	defer ch.Close()

	// Declarar la cola "motion"
	_, err = ch.QueueDeclare(
		"motion",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error al declarar la cola: %v", err)
	}

	// Consumir mensajes de la cola "motion"
	msgs, err := ch.Consume(
		"motion",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error al consumir mensajes: %v", err)
	}

	
	go func() {
		for d := range msgs {
			log.Printf("Mensaje recibido en motion: %s", d.Body)
			url := secondAPIURL + "/sensor/motion"

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(d.Body))
			if err != nil {
				log.Printf("Error al crear request HTTP: %v", err)
				continue
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error al enviar request HTTP: %v", err)
				continue
			}
			resp.Body.Close()

			log.Printf("Mensaje reenviado a %s", url)
		}
	}()

	log.Println("Escuchando en la cola motion...")
	select {} // Bloquea la ejecución para que siga corriendo
}
