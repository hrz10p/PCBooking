package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type Message struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	Recipient string `json:"recipient"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email_queue", // Имя очереди
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go processMessages(msgs)

	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	forever := make(chan struct{})
	<-forever
}

func processMessages(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		var message Message
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}
		go sendEmail(message)
	}
}
