package usecase

import (
	"Test_App/internal/domain"
	"context"
	"errors"
)

var ErrSectionNotFound = errors.New("section not found")

type sectionUseCase struct {
	sectionRepo domain.SectionRepository
	subjectRepo domain.SubjectRepository
}

func NewSectionUseCase(
	sectionRepo domain.SectionRepository,
	subjectRepo domain.SubjectRepository,
) domain.SectionUseCase {
	return &sectionUseCase{
		sectionRepo: sectionRepo,
		subjectRepo: subjectRepo,
	}
}

// GetSectionsBySubjectID — список разделов предмета
// Проверяем что предмет существует прежде чем искать его разделы
func (uc *sectionUseCase) GetSectionsBySubjectID(ctx context.Context, subjectID int64) ([]*domain.Section, error) {
	// 1. проверяем что предмет существует
	subject, err := uc.subjectRepo.GetSubjectByID(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	if subject == nil {
		return nil, ErrSubjectNotFound
	}

	// 2. загружаем разделы
	sections, err := uc.sectionRepo.GetSectionsBySubjectID(ctx, subjectID)
	if err != nil {
		return nil, err
	}

	return sections, nil
}
