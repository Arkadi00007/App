package tests

import (
	"Test_App/internal/domain"
	"context"
	"database/sql"
)

type subjectRepository struct {
	db *sql.DB
}

func NewSubjectRepository(db *sql.DB) domain.SubjectRepository {
	return &subjectRepository{db: db}
}

// GetAll — возвращает все предметы (студент видит список на главном экране)
func (r *subjectRepository) GetAllSubjects(ctx context.Context) ([]*domain.Subject, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description
		FROM subjects
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []*domain.Subject
	for rows.Next() {
		s := &domain.Subject{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Description); err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}

	// rows.Err() проверяет ошибки которые могли возникнуть ВО ВРЕМЯ итерации
	// (не на старте запроса, а где-то посреди чтения строк)
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}

func (r *subjectRepository) GetSubjectByID(ctx context.Context, id int64) (*domain.Subject, error) {
	s := &domain.Subject{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, description
		FROM subjects
		WHERE id = $1
	`, id).Scan(&s.ID, &s.Name, &s.Description)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return s, nil
}
