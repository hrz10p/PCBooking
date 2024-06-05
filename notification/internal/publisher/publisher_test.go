package publisher_test

import (
	"encoding/json"
	"notification/internal/publisher"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublishMessage(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	require.NoError(t, err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	require.NoError(t, err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	require.NoError(t, err, "Failed to declare a queue")

	msg := publisher.Message{
		Type:      "test",
		Message:   "This is a test message",
		Recipient: "test@example.com",
	}

	err = publisher.PublishMessage(msg)
	require.NoError(t, err, "Failed to publish message")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	require.NoError(t, err, "Failed to register a consumer")

	done := make(chan bool)

	go func() {
		for d := range msgs {
			var receivedMessage publisher.Message
			err := json.Unmarshal(d.Body, &receivedMessage)
			require.NoError(t, err, "Failed to unmarshal message")

			assert.Equal(t, msg, receivedMessage, "The received message does not match the sent message")
			done <- true
			break
		}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("No messages received within the timeout period")
	}
}
