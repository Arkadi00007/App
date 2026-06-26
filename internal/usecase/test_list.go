// usecase/test_list.go
package usecase

import (
	"Test_App/internal/domain"
	"context"
)

type testListUseCase struct {
	testRepo    domain.TestRepository
	sectionRepo domain.SectionRepository
}

func NewTestListUseCase(
	testRepo domain.TestRepository,
	sectionRepo domain.SectionRepository,
) domain.TestListUseCase {
	return &testListUseCase{
		testRepo:    testRepo,
		sectionRepo: sectionRepo,
	}
}

func (uc *testListUseCase) GetTestsBySectionID(ctx context.Context, sectionID int64) ([]*domain.Test, error) {
	// проверяем что раздел существует
	section, err := uc.sectionRepo.GetSectionByID(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	if section == nil {
		return nil, ErrSectionNotFound
	}

	return uc.testRepo.GetTestsBySectionID(ctx, sectionID)
}
