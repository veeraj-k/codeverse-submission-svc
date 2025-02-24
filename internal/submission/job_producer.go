package submission

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type Job struct {
	SubmissionId uuid.UUID
	ProblemId    uint
	UserId       uint
	Code         string
	Language     string
}

type JobProducer struct {
	channel *amqp091.Channel
	queue   string
}

func NewJobProducer(channel *amqp091.Channel, queue string) *JobProducer {
	return &JobProducer{
		channel: channel,
		queue:   queue,
	}
}

func (j *JobProducer) ProduceJob(job *Job) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(job)
	if err != nil {
		panic(err)

	}
	err = j.channel.PublishWithContext(
		ctx,
		"",
		j.queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/json",
			Body:        body,
		})

	if err != nil {
		panic(err)
	}
}
