package main

import (
	"submission-service/config"
	"submission-service/db"
	"submission-service/internal/submission"
	"submission-service/rmq"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

var DB *gorm.DB
var conn *amqp091.Connection

func init() {
	config.Loadenv()
	DB = db.InitDB()
	db.MigrateDB()
	db.SeedDB()
	conn = rmq.SetupRabbitMQ()

}
func main() {

	r := gin.Default()

	submission.RegisterRoutes(r, DB, conn)

	r.Run("127.0.0.1:8082")

}
