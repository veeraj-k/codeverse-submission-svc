package submission

type SubmissionService struct {
	repo        SubmissionRepository
	rmqProducer JobProducer
}

func NewSubmissionService(repo SubmissionRepository, rmqProducer *JobProducer) *SubmissionService {
	return &SubmissionService{repo: repo, rmqProducer: *rmqProducer}
}

func (s *SubmissionService) CreateSubmission(submissionRequest *SubmissionRequest) (*Submission, error) {

	submission := &Submission{
		UserId:          submissionRequest.UserId,
		ProblemId:       submissionRequest.ProblemId,
		Code:            submissionRequest.Code,
		Status:          "Pending",
		TotalTestCases:  0,
		TestCasesPassed: 0,
	}

	err := s.repo.CreateSubmission(submission)
	if err != nil {
		return nil, err
	}

	job := &Job{
		SubmissionId: submission.ID,
		ProblemId:    submission.ProblemId,
		Code:         submission.Code,
		UserId:       submission.UserId,
		Language:     "fix this",
	}

	s.rmqProducer.ProduceJob(job)

	return submission, nil
}
