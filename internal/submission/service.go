package submission

type SubmissionService struct {
	repo              SubmissionRepository
	jobRmqProducer    JobProducer
	statusRmqProducer StatusProducer
}

func NewSubmissionService(repo SubmissionRepository, jobRmqProducer *JobProducer, statusConsumer *StatusProducer) *SubmissionService {
	return &SubmissionService{repo: repo, jobRmqProducer: *jobRmqProducer, statusRmqProducer: *statusConsumer}
}

func (s *SubmissionService) CreateSubmission(submissionRequest *SubmissionRequest) (*Submission, error) {

	submission := &Submission{
		UserId:          submissionRequest.UserId,
		ProblemId:       submissionRequest.ProblemId,
		Code:            submissionRequest.Code,
		Status:          "Pending",
		TotalTestCases:  0,
		TestCasesPassed: 0,
		Language:        submissionRequest.Language,
	}

	err := s.repo.CreateSubmission(submission)
	if err != nil {
		return nil, err
	}

	job := &SubmissionJob{
		SubmissionId: submission.ID,
		ProblemId:    submission.ProblemId,
		Code:         submission.Code,
		UserId:       submission.UserId,
		Language:     submission.Language,
	}

	s.jobRmqProducer.ProduceJob(job)
	s.statusRmqProducer.ProduceStatus(&SubmissionStatus{
		SubmissionId: submission.ID,
		Status:       "In QUEUE",
	})

	return submission, nil
}
