package domain

import "context"

type Subject struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SubjectRepository interface {
	GetAllSubjects(ctx context.Context) ([]*Subject, error)
	GetSubjectByID(ctx context.Context, id int64) (*Subject, error)
}

type SubjectService interface {
	GetAll(ctx context.Context) ([]*Subject, error)
	GetByID(ctx context.Context, id int64) (*Subject, error)
}
