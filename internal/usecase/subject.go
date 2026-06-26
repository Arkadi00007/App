package usecase

import (
	"Test_App/internal/domain"
	"context"
	"errors"
)

var ErrSubjectNotFound = errors.New("subject not found")

type subjectUseCase struct {
	subjectRepo domain.SubjectRepository
}

func NewSubjectUseCase(subjectRepo domain.SubjectRepository) domain.SubjectUseCase {
	return &subjectUseCase{subjectRepo: subjectRepo}
}

// GetAllSubjects — список всех предметов для главного экрана
func (uc *subjectUseCase) GetAllSubjects(ctx context.Context) ([]*domain.Subject, error) {
	subjects, err := uc.subjectRepo.GetAllSubjects(ctx)
	if err != nil {
		return nil, err
	}
	return subjects, nil
}

// GetSubjectByID — конкретный предмет (может понадобиться для заголовка страницы)
func (uc *subjectUseCase) GetSubjectByID(ctx context.Context, id int64) (*domain.Subject, error) {
	subject, err := uc.subjectRepo.GetSubjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if subject == nil {
		return nil, ErrSubjectNotFound
	}
	return subject, nil
}
