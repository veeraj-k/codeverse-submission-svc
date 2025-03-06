package rmq

import (
	"fmt"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func SetupRabbitMQ() *amqp091.Connection {

	conn, err := amqp091.Dial(fmt.Sprintf("amqps://%s/%s", os.Getenv("RMQ_HOST"), os.Getenv("RMQ_VHOST")))
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
		"submission_status_queue",
		true,
		false,
		false,
		false,
		amqp091.Table{},
	)
	if err != nil {
		log.Fatal(err)
	}

	return conn

}
