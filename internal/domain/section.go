package domain

import "context"

type Section struct {
	ID          int64  `json:"id"`
	SubjectID   int64  `json:"subject_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type SectionRepository interface {
	GetSectionsBySubjectID(ctx context.Context, subjectID int64) ([]*Section, error)
	GetSectionByID(ctx context.Context, id int64) (*Section, error)
}

type SectionUseCase interface {
	GetSectionsBySubjectID(ctx context.Context, subjectID int64) ([]*Section, error)
}
