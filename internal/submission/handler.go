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

func (h *SubmissionHandler) GetSubmissionById(c *gin.Context) {

	sid := c.Param("id")

	submission, err := h.service.repo.GetSubmissionById(sid)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Submission with id not found",
		})
		return
	}
	c.JSON(http.StatusOK, submission)

}

func (h *SubmissionHandler) GetSubmissionByQueryParam(c *gin.Context) {
	params := make(map[string]interface{})
	for key, value := range c.Request.URL.Query() {
		params[key] = value[0]
	}

	submissions, err := h.service.repo.GetSubmissionsByQueryParams(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, submissions)
}

func (h *SubmissionHandler) AddSubmissionTestCases(c *gin.Context) {
	var testCases []SubmissionTestCases
	if err := c.BindJSON(&testCases); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sid := c.Param("id")

	submission, err := h.service.repo.GetSubmissionById(sid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	var test_cases_passed uint
	var total_test_cases uint

	for _, testCase := range testCases {
		if testCase.Status == "passed" {
			test_cases_passed++
		}
		total_test_cases++
	}

	submission.TestCasesPassed = test_cases_passed
	submission.TotalTestCases = total_test_cases

	if err := h.service.repo.AddSubmissionTestCases(submission, testCases); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Test cases added successfully",
	})
}

func (h *SubmissionHandler) UpdateSubmissionStatus(c *gin.Context) {

	sid := c.Param("id")

	submission, err := h.service.repo.GetSubmissionById(sid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	var statusBody SubmissionStatus
	if err := c.BindJSON(&statusBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// status := c.Request.URL.Query().Get("status")

	// if status == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
	// 	return
	// }

	if err := h.service.repo.UpdateStatus(submission, statusBody.Status, statusBody.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.service.statusRmqProducer.ProduceStatus(&SubmissionStatus{
		SubmissionId: submission.ID,
		Status:       statusBody.Status,
		Message:      statusBody.Message,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Submission status updated successfully",
	})
}
