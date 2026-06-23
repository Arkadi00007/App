// package postgres
//
// import (
//
//	"Test_App/internal/domain"
//	"database/sql"
//
// )
//
//	type userRepository struct {
//		db *sql.DB
//	}
//
//	func NewUserRepository(db *sql.DB) domain.UserRepository {
//		return &userRepository{db: db}
//	}
//
//	func (r *userRepository) CreateUser(user *domain.User) error {
//		return r.db.QueryRow(`
//	       INSERT INTO users (email, password_hash, name, role)
//	       VALUES ($1, $2, $3, 'student')
//	       RETURNING id, created_at
//	   `, user.Email, user.PasswordHash, user.Name).
//			Scan(&user.ID, &user.CreatedAt)
//	}
//
// // ищем по email — нужно при логине
//
//	func (r *userRepository) GetUserByEmail(email string) (*domain.User, error) {
//		user := &domain.User{}
//		err := r.db.QueryRow(`
//	       SELECT id, email, password_hash, name, role, created_at
//	       FROM users
//	       WHERE email = $1
//	   `, email).Scan(
//			&user.ID,
//			&user.Email,
//			&user.PasswordHash,
//			&user.Name,
//			&user.Role,
//			&user.CreatedAt,
//		)
//
//		if err == sql.ErrNoRows {
//			return nil, nil // не найден — возвращаем nil без ошибки
//		}
//		if err != nil {
//			return nil, err // реальная ошибка БД
//		}
//
//		return user, nil
//	}
//
// // ищем по id — нужно для middleware (проверка токена)
//
//	func (r *userRepository) GetUserByID(id int64) (*domain.User, error) {
//		user := &domain.User{}
//		err := r.db.QueryRow(`
//	       SELECT id, email, name, role, created_at
//	       FROM users
//	       WHERE id = $1
//	   `, id).Scan(
//			&user.ID,
//			&user.Email,
//			&user.Name,
//			&user.Role,
//			&user.CreatedAt,
//		)
//
//		if err == sql.ErrNoRows {
//			return nil, nil
//		}
//		if err != nil {
//			return nil, err
//		}
//
//		return user, nil
//	}
//
//	func (r *userRepository) Update(user *domain.User) error {
//		_, err := r.db.Exec(`
//			UPDATE users SET
//				password_hash = $1,
//				is_verified   = $2
//			WHERE id = $3
//		`, user.PasswordHash, user.IsVerified, user.ID)
//		return err
//	}
package postgres

import (
	"Test_App/internal/domain"
	"context"
	"database/sql"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, name, role, is_verified)
		VALUES ($1, $2, $3, 'student', false)
		RETURNING id, created_at
	`, user.Email, user.PasswordHash, user.Name).
		Scan(&user.ID, &user.CreatedAt)
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, name, role, is_verified, created_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Role,
		&user.IsVerified,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, name, role, is_verified, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.IsVerified,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET
			password_hash = $1,
			is_verified   = $2
		WHERE id = $3
	`, user.PasswordHash, user.IsVerified, user.ID)
	return err
}
