package submission

import (
	"fmt"
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
	msgs, _ := c.channel.Consume(
		"submision_status_stream",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	for msg := range msgs {
		fmt.Println("Received status update:", string(msg.Body))
		c.hub.Broadcast <- msg.Body
	}
}
