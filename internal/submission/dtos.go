package submission

type SubmissionRequest struct {
	UserId    uint   `json:"user_id"`
	ProblemId uint   `json:"problem_id"`
	Language  string `json:"language"`
	Code      string `json:"code"`
}
