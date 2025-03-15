package adapters

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
)

// ConsumeMotionQueue se conecta a RabbitMQ, consume mensajes de la cola "esp32"
// y los reenvía al endpoint /sensor/motion de tu API.
func ConsumeMotionQueue() {
	rabbitURL := os.Getenv("RABBIT_URL")
	secondAPIURL := os.Getenv("SECOND_API_URL")

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Error al conectar a RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error al abrir canal: %v", err)
	}

	// Declarar la cola "esp32"
	q, err := ch.QueueDeclare(
		"esp32", // nombre de la cola
		true,    // durable
		false,   // auto-delete
		false,   // exclusivo
		false,   // no-wait
		nil,     // argumentos
	)
	if err != nil {
		log.Fatalf("Error al declarar la cola: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // nombre de la cola
		"",     // consumer tag
		true,   // auto-acknowledge
		false,  // exclusivo
		false,  // no-local
		false,  // no-wait
		nil,    // argumentos
	)
	if err != nil {
		log.Fatalf("Error al consumir mensajes: %v", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Mensaje recibido en cola 'esp32': %s", d.Body)

			// Construir la URL de reenvío: se asume que el endpoint es /sensor/motion
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
			defer resp.Body.Close()

			log.Printf("Mensaje reenviado a %s", url)
		}
	}()
}
