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

	jobRmqProducer := NewJobProducer(channel, "code_execution_job_queue")
	if jobRmqProducer == nil {
		panic("Failed to create job producer")
	}

	statusProducer := NewStatusProducer(channel, "submission_status_queue")
	if statusProducer == nil {
		panic("Failed to create status producer")
	}

	submissionService := NewSubmissionService(NewSubmissionRepository(db), jobRmqProducer, statusProducer)
	if submissionService == nil {
		panic("Failed to create submission service")
	}
	contestStatusProducer := NewContestStatusProducer(channel, "contest_user_submissions")
	if contestStatusProducer == nil {
		panic("Failed to create contest status producer")
	}

	submissionHandler := NewSubmissionHandler(submissionService, contestStatusProducer)
	if submissionHandler == nil {
		panic("Failed to create submission handler")
	}

	fmt.Println("Registering routes for submission", submissionHandler)
	r.POST("/api/submission/", submissionHandler.CreateSubmission)
	r.GET("/api/submission/status/:id", func(ctx *gin.Context) {
		submissionHandler.SubmissionStatus(ctx, hub)
	})
	r.GET("/api/submission/:id", submissionHandler.GetSubmissionById)
	r.GET("/api/submission/", submissionHandler.GetSubmissionByQueryParam)
	r.PUT("/api/submission/:id/testcases", submissionHandler.AddSubmissionTestCases)
	r.PUT("/api/submission/:id/status", submissionHandler.UpdateSubmissionStatus)
	// r.GET("/submissions/:id", GetSubmissionById)
	// r.GET("/submissions", submissionHandler.)
}
