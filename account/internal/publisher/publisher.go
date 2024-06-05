package publisher

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type Message struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	Recipient string `json:"recipient"`
}

type RabbitPublisher struct {
	rabbitUrl string
}

func NewRabbitPublisher(rabbitUrl string) *RabbitPublisher {
	return &RabbitPublisher{rabbitUrl: rabbitUrl}
}

func (r *RabbitPublisher) PublishMessage(msg Message) error {
	conn, err := amqp.Dial(r.rabbitUrl)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("RABBIT [x] Message sent to %s", msg.Recipient)
	return nil
}
