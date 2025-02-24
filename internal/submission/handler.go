package submission

import (
	"fmt"
	"net/http"
	"submission-service/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type SubmissionHandler struct {
	service *SubmissionService
}

func NewSubmissionHandler(service *SubmissionService) *SubmissionHandler {
	return &SubmissionHandler{service: service}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {
	var submission SubmissionRequest
	if err := c.BindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdSubmission, err := h.service.CreateSubmission(&submission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Submission created successfully",
		"submission": createdSubmission,
		"ws_url":     fmt.Sprintf("ws://%s/submission/status/%s", c.Request.Host, createdSubmission.ID),
	})

}

func (h *SubmissionHandler) SubmissionStatus(c *gin.Context, hub *ws.Hub) {
	clientID := c.Param("id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	// defer conn.Close()

	_, err = h.service.repo.GetSubmissionById(clientID)

	if err != nil {
		conn.WriteJSON(gin.H{"error": "Submission not found"})
		return
	}

	client := &ws.Client{
		ID:     clientID,
		Conn:   conn,
		SendCh: make(chan []byte),
	}

	hub.Register <- client

	// go client.ReadPump(hub)
	go client.WritePump()

}
