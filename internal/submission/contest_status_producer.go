package submission

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type ContestStatus struct {
	SubmissionId uuid.UUID `json:"job_id"`
	ProblemId    uint      `json:"problem_id"`
	UserId       uint      `json:"user_id"`
	TotalTests   uint      `json:"total_tests"`
	PassedTests  uint      `json:"test_cases_passed"`
	Timestamp    int32     `json:"timestamp"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
}

type ContestStatusProducer struct {
	channel *amqp091.Channel
	queue   string
}

func NewContestStatusProducer(channel *amqp091.Channel, queue string) *ContestStatusProducer {
	return &ContestStatusProducer{
		channel: channel,
		queue:   queue,
	}
}

func (s *ContestStatusProducer) ProduceStatus(status *ContestStatus) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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
