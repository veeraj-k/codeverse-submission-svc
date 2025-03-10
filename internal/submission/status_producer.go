package submission

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type SubmissionStatus struct {
	SubmissionId uuid.UUID `json:"job_id"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
}

type StatusProducer struct {
	channel *amqp091.Channel
	queue   string
}

func NewStatusProducer(channel *amqp091.Channel, queue string) *StatusProducer {
	return &StatusProducer{
		channel: channel,
		queue:   queue,
	}
}

func (s *StatusProducer) ProduceStatus(status *SubmissionStatus) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(status)
	if err != nil {
		panic(err)

	}
	if s.channel.IsClosed() {
		panic("channel/connection is not open")
	}

	err = s.channel.PublishWithContext(
		ctx,
		"",
		s.queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/json",
			Body:        body,
		})

	if err != nil {
		panic(err)
	}
	fmt.Println("Published status update:", status)
}
