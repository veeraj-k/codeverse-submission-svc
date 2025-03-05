package submission

import (
	"fmt"
	"submission-service/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, rconn *amqp091.Connection) {

	hub := ws.NewHub()

	go hub.Run()

	channel, err := rconn.Channel()
	if err != nil {
		panic(fmt.Sprintf("Failed to open a channel: %v", err))
	}
	// defer channel.Close()

	stsConsumer := NewStatusConsumer(channel, hub)
	if stsConsumer == nil {
		panic("Failed to create status consumer")
	}
	go stsConsumer.StartConsuming()

	rmqProducer := NewJobProducer(channel, "code_execution_job_queue")
	if rmqProducer == nil {
		panic("Failed to create job producer")
	}

	submissionService := NewSubmissionService(NewSubmissionRepository(db), rmqProducer)
	if submissionService == nil {
		panic("Failed to create submission service")
	}

	submissionHandler := NewSubmissionHandler(submissionService)
	if submissionHandler == nil {
		panic("Failed to create submission handler")
	}

	fmt.Println("Registering routes for submission", submissionHandler)
	r.POST("/submissions", submissionHandler.CreateSubmission)
	r.GET("/submission/status/:id", func(ctx *gin.Context) {
		submissionHandler.SubmissionStatus(ctx, hub)
	})
	r.GET("/submission/:id", submissionHandler.GetSubmissionById)
	r.GET("/submissions", submissionHandler.GetSubmissionByQueryParam)
	r.PUT("/submission/:id/testcases", submissionHandler.AddSubmissionTestCases)
	r.PUT("/submission/:id/status", submissionHandler.UpdateSubmissionStatus)
	// r.GET("/submissions/:id", GetSubmissionById)
	// r.GET("/submissions", submissionHandler.)
}
