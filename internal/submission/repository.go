package submission

import "gorm.io/gorm"

type SubmissionRepository interface {
	CreateSubmission(submission *Submission) error
	GetSubmissionById(id string) (*Submission, error)
	GetSubmissionsByQueryParams(params map[string]interface{}) ([]Submission, error)
}

type submissionRepository struct {
	DB *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{DB: db}
}

func (r *submissionRepository) CreateSubmission(submission *Submission) error {

	if err := r.DB.Create(submission).Error; err != nil {
		return err
	}

	return nil

}

func (r *submissionRepository) GetSubmissionById(id string) (*Submission, error) {

	var submission Submission
	if err := r.DB.Where("id = ?", id).First(&submission).Error; err != nil {
		return nil, err
	}

	return &submission, nil
}

func (r *submissionRepository) GetSubmissionsByQueryParams(params map[string]interface{}) ([]Submission, error) {
	var submissions []Submission
	query := r.DB

	for key, value := range params {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}
