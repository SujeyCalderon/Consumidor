package adapters

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
)

type LightSensorMessage struct {
	ID         string  `json:"id"`
	Timestamp  string  `json:"timestamp"`
	Luminosity float64 `json:"luminosity"`
}

func ConsumeLightQueue() {
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
	

	q, err := ch.QueueDeclare(
		"light",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error al declarar la cola: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error al consumir mensajes: %v", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Mensaje recibido en light: %s", d.Body)
			url := secondAPIURL + "/sensor/light"
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
}
