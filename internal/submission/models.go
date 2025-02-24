package submission

import (
	"github.com/google/uuid"
)

type Submission struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId          uint      `json:"user_id" gorm:"not null"`
	ProblemId       uint      `json:"problem_id" gorm:"not null"`
	Language        string    `json:"language" gorm:"not null"`
	Code            string    `json:"code" gorm:"not null"`
	Status          string    `json:"status" gorm:"not null"`
	TestCasesPassed uint      `json:"test_cases_passed" gorm:"not null"`
	TotalTestCases  uint      `json:"total_test_cases" gorm:"not null"`
}
