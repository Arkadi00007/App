package tests

import (
	"Test_App/internal/domain"
	"context"
	"database/sql"
)

type sectionRepository struct {
	db *sql.DB
}

func NewSectionRepository(db *sql.DB) domain.SectionRepository {
	return &sectionRepository{db: db}
}

func (r *sectionRepository) GetSectionBySubjectID(ctx context.Context, subjectID int64) ([]*domain.Section, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, subject_id, name, description
		FROM sections
		WHERE subject_id = $1
		ORDER BY id
	`, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []*domain.Section
	for rows.Next() {
		s := &domain.Section{}
		if err := rows.Scan(&s.ID, &s.SubjectID, &s.Name, &s.Description); err != nil {
			return nil, err
		}
		sections = append(sections, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sections, nil
}

func (r *sectionRepository) GetSectionByID(ctx context.Context, id int64) (*domain.Section, error) {
	s := &domain.Section{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, subject_id, name, description
		FROM sections
		WHERE id = $1
	`, id).Scan(&s.ID, &s.SubjectID, &s.Name, &s.Description)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return s, nil
}
