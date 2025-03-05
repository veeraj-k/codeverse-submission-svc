package submission

import (
	"github.com/google/uuid"
)

type Submission struct {
	ID              uuid.UUID             `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId          uint                  `json:"user_id" gorm:"not null"`
	ProblemId       uint                  `json:"problem_id" gorm:"not null"`
	Language        string                `json:"language" gorm:"not null"`
	Code            string                `json:"code" gorm:"not null"`
	Status          string                `json:"status" gorm:"not null"`
	TestCasesPassed uint                  `json:"test_cases_passed" gorm:"not null"`
	TotalTestCases  uint                  `json:"total_test_cases" gorm:"not null"`
	SubmissionTests []SubmissionTestCases `json:"submission_tests" gorm:"foreignKey:submission_id"`
}

type SubmissionTestCases struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	SubmissionId uuid.UUID `json:"submission_id" gorm:"not null"`
	Input        string    `json:"input" gorm:"not null"`
	Output       string    `json:"output" gorm:"not null"`
	Actual       string    `json:"actual" gorm:"not null"`
	Status       string    `json:"status" gorm:"not null"`
	Error        string    `json:"error" gorm:"not null"`
	Stdout       string    `json:"stdout" gorm:"not null"`
}
