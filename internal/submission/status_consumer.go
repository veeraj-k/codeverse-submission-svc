package submission

import (
	"submission-service/internal/ws"

	"github.com/rabbitmq/amqp091-go"
)

type StatusConsumer struct {
	channel *amqp091.Channel
	hub     *ws.Hub
}

func NewStatusConsumer(channel *amqp091.Channel, hub *ws.Hub) *StatusConsumer {
	return &StatusConsumer{
		channel: channel,
		hub:     hub,
	}
}

func (c *StatusConsumer) StartConsuming() {
	msgs, err := c.channel.Consume(
		"submission_status_queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		// fmt.Println("Received status update:", string(msg.Body))
		c.hub.Broadcast <- msg.Body
	}
}
