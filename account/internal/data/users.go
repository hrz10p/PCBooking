package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type UserModel struct {
	DB *sql.DB
}

var AnonymousUser = &User{}

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FName     string    `json:"f_name"`
	SName     string    `json:"s_name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	UserRole  string    `json:"user_role"`
	Activated bool      `json:"activated"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (id, created_at, updated_at, fname, sname, email, password_hash, user_role, activated)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`
	args := []any{time.Now(), time.Now(), user.FName, user.SName, user.Email, user.Password.hash, user.UserRole, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, created_at, updated_at, fname, sname, email, password_hash, user_role, activated
		FROM users
		WHERE email = $1`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FName,
		&user.SName,
		&user.Email,
		&user.Password.hash,
		&user.UserRole,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Get(id string) (*User, error) {
	query := `
		SELECT id, created_at, updated_at, fname, sname, email, password_hash, user_role, activated
		FROM users
		WHERE id = $1`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FName,
		&user.SName,
		&user.Email,
		&user.Password.hash,
		&user.UserRole,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET updated_at = $1, fname = $2, sname = $3, email = $4, password_hash = $5, user_role = $6, activated = $7
		WHERE id = $8
		RETURNING updated_at`
	args := []any{time.Now(), user.FName, user.SName, user.Email, user.Password.hash, user.UserRole, user.Activated, user.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt)
	return err
}

func (m UserModel) Delete(id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m UserModel) GetAll() ([]*User, error) {
	query := `
		SELECT id, created_at, updated_at, fname, sname, email, password_hash, user_role, activated
		FROM users`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.FName, &user.SName, &user.Email, &user.Password.hash, &user.UserRole, &user.Activated)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrRecordNotFound
	}
	return users, nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	query := `
		SELECT users.id, users.created_at, users.updated_at, users.fname, users.sname, users.email, users.password_hash, users.user_role, users.activated
		FROM users
		INNER JOIN tokens
		ON users.id = tokens.user_id
		WHERE tokens.hash = $1
		AND tokens.scope = $2
		AND tokens.expiry > $3`
	args := []any{tokenHash[:], tokenScope, time.Now()}
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FName,
		&user.SName,
		&user.Email,
		&user.Password.hash,
		&user.UserRole,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
