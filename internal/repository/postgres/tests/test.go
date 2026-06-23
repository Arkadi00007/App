package tests

import (
	"Test_App/internal/domain"
	"context"
	"database/sql"
)

type testRepository struct {
	db *sql.DB
}

func NewTestRepository(db *sql.DB) domain.TestRepository {
	return &testRepository{db: db}
}

// GetBySectionID — список тестов внутри раздела
func (r *testRepository) GetTestsBySectionID(ctx context.Context, sectionID int64) ([]*domain.Test, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, section_id, title, description, show_answer_mode, is_premium, created_at
		FROM tests
		WHERE section_id = $1
		ORDER BY id
	`, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []*domain.Test
	for rows.Next() {
		t := &domain.Test{}
		if err := rows.Scan(
			&t.ID, &t.SectionID, &t.Title, &t.Description,
			&t.ShowAnswerMode, &t.IsPremium, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tests = append(tests, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tests, nil
}

// GetByID — конкретный тест (нужен когда студент открывает тест)
func (r *testRepository) GetTestByID(ctx context.Context, id int64) (*domain.Test, error) {
	t := &domain.Test{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, section_id, title, description, show_answer_mode, is_premium, created_at
		FROM tests
		WHERE id = $1
	`, id).Scan(
		&t.ID, &t.SectionID, &t.Title, &t.Description,
		&t.ShowAnswerMode, &t.IsPremium, &t.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return t, nil
}
