package main

import (
	"os"
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
	if os.Getenv("LOAD_ENV") == "true" {
		config.Loadenv()
	}
	DB = db.InitDB()
	db.MigrateDB()
	db.SeedDB()
	conn = rmq.SetupRabbitMQ()

}
func main() {

	r := gin.Default()

	// r.Use(cors.New(cors.Config{
	// 	AllowAllOrigins:  true,
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	r.Handle("GET", "/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	submission.RegisterRoutes(r, DB, conn)

	r.Run("127.0.0.1:8082")

}
