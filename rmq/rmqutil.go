package rmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func SetupRabbitMQ() *amqp091.Connection {
	conn, err := amqp091.Dial("amqps://RMQ_HOST/zoiqirnx")
	if err != nil {
		log.Fatal(err)
	}
	// defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"code_execution_job_queue",
		true,
		false,
		false,
		false,
		amqp091.Table{},
	)
	if err != nil {
		log.Fatal(err)
	}
	_, err = ch.QueueDeclare(
		"submision_status_stream",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-queue-type": "stream",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return conn

}
